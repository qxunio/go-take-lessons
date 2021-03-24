package service

import (
	"encoding/json"
	"errors"
	"github.com/go-xorm/xorm"
	"go-take-lessons/db"
	"go-take-lessons/domain"
	"go-take-lessons/domain/comm"
	"go-take-lessons/domain/comm/ere"
	"go-take-lessons/model"
	"go-take-lessons/tools"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"sync"
	"time"
)

type eventService struct {
	mysql *xorm.Engine
}

// 学生查询历史
func (e eventService) PostStuHistoryEvent(user *comm.SessionUSER) ([]domain.EventSimpleVo, error) {
	var eventIds []string
	sql := "select DISTINCT(event_id) FROM stu_subject WHERE user_id = ? AND `enable` = ?"
	if err := e.mysql.SQL(sql, user.Id, comm.Enable).Find(&eventIds); err != nil {
		zap.S().Error("学生查询历史，查询学生选课异常 user_id: ", user.Id, err.Error())
		return nil, errors.New("查询异常")
	}

	var eventLis []model.Event
	var vo []domain.EventSimpleVo

	if eventIds == nil {
		return vo, nil
	}

	if err := e.mysql.Where("enable = ?", comm.Enable).In("id", eventIds).Find(&eventLis); err != nil {
		zap.S().Error("学生查询历史，user_id: ", user.Id, err.Error())
		return nil, errors.New("查询异常")
	}

	for _, event := range eventLis {
		vo = append(vo, domain.CvEventSimpleToVo(event))
	}

	return vo, nil
}

// 修改学生是否可以修改选课
func (e eventService) PostUpdateCanUpdate(arg domain.EventUpdateCanUpdateArg, user *comm.SessionUSER) error {
	var event model.Event
	if _, err := e.mysql.Where("enable = ?", comm.Enable).ID(arg.Id).Get(&event); err != nil {
		return ere.ErrorCommFindError
	}

	event.CanUpdate = tools.StringToInt(arg.Status)
	event.Utime = time.Now()
	if _, err := e.mysql.Cols("can_update", "utime").ID(event.Id).Update(event); err != nil {
		zap.S().Error(err)
		return ere.ErrorCommUpdateError
	}
	e.updateCache(&event)
	return nil
}

// 重新激活选课
func (e eventService) PostReactivation(id string, etime time.Time) error {
	var event model.Event
	if _, err := e.mysql.Where("enable = ?", comm.Enable).ID(id).Get(&event); err != nil {
		return ere.ErrorCommFindError
	}
	event.Status = comm.EventStatusEnable
	event.Etime = etime
	event.Utime = time.Now()
	if _, err := e.mysql.Id(event.Id).Update(event); err != nil {
		zap.S().Error(err)
		return ere.ErrorCommUpdateError
	}
	e.updateCache(&event)
	return nil
}

// 检查登录用户是否有选课
func (e eventService) isEffectiveEvent(event *model.Event, sessionUser *comm.SessionUSER) bool {
	if time.Now().After(event.Etime) {
		event.Utime = time.Now()
		event.Status = comm.EventStatusHistory
		if _, err := e.mysql.Id(event.Id).Cols("status", "utime").Update(event); err != nil {
			zap.S().Error(err)
		}
		e.removeCache()
		return false
	}

	if sessionUser.UserType == comm.Admin {
		return true
	}

	if event.SchoolYear == "全部" {
		return true
	}

	for _, sy := range strings.Split(event.SchoolYear, ",") {
		if strconv.FormatInt(int64(sessionUser.SchoolYear), 10) == sy {
			return true
		}
	}

	split := strings.Split(event.TagIds, ",")
	var tagStu model.TagStu
	count, err := e.mysql.Where("uid = ?", sessionUser.Id).In("tag_id", split).Count(&tagStu)
	if err != nil {
		zap.S().Error(err)
		return false
	}
	return count >= 1
}

func (e eventService) GetRedisSet(key string, value string) interface{} {
	result, err := db.RedisClient.Set(key, value, 0).Result()
	if err != nil {
		zap.S().Error(err)
	}
	return result
}

func (e eventService) GetRedisDel(key string) interface{} {
	result, err := db.RedisClient.Del(key).Result()
	if err != nil {
		zap.S().Error(err)
	}
	return result
}

func (e eventService) GetRedisGet(id string) interface{} {

	get, err := db.RedisClient.Get(id).Result()
	if err != nil {
		zap.S().Error(err)
	}
	return get
}

var updateExistEventCacheLock sync.Mutex

// 查询存在的有效选课
func (e eventService) PostExistEvent(sessionUser *comm.SessionUSER) (*domain.EventVo, error) {
	cache := e.getEventCache()
	if cache != nil {
		if cache.Status == comm.EventStatusEnable && e.isEffectiveEvent(cache, sessionUser) {
			return domain.CvEventToStuVo(cache), nil
		}
	} else {
		updateExistEventCacheLock.Lock()
		defer updateExistEventCacheLock.Unlock()
		// 这里从缓存拿是，防止下一个线程去查询数据表并更新缓存
		cache := e.getEventCache()
		if cache != nil && e.isEffectiveEvent(cache, sessionUser) {
			return domain.CvEventToStuVo(cache), nil
		}
		var event model.Event
		has, err := e.mysql.Where("enable = ?", comm.Enable).And("status = ?",
			comm.Enable).Get(&event)
		if err != nil {
			zap.S().Error(err)
			return nil, ere.ErrorCommFindError
		}
		if has && e.isEffectiveEvent(&event, sessionUser) {
			e.updateCache(&event)
			return domain.CvEventToStuVo(&event), nil
		}
	}
	return nil, errors.New("当前没有选课活动")
}

// 修改选课状态
func (e eventService) PostChange(eventId string, sessionUser *comm.SessionUSER) error {
	if sessionUser.UserType != comm.Admin {
		return errors.New("非法操作")
	}
	var event model.Event
	has, err := e.mysql.Id(eventId).Get(&event)
	if err != nil {
		zap.S().Error(err)
		return ere.ErrorCommFindError
	}
	if !has {
		return errors.New("选课不存在")
	}

	if event.Status == comm.EventStatusEnable {
		event.Status = comm.EventStatusDisable
	} else {
		if e.existDoEvent() {
			return errors.New("当前存在正在进行的选课，每次只能存在一条有效的选课")
		}
		event.Status = comm.EventStatusEnable
	}
	event.Utime = time.Now()
	if _, err := e.mysql.Id(eventId).Cols("status", "utime").Update(&event); err != nil {
		zap.S().Error(err)
		return ere.ErrorCommUpdateError
	}
	if event.Status == comm.EventStatusDisable {
		e.removeCache()
	} else {
		e.updateCache(&event)
	}
	return nil
}

// 删除
func (e eventService) PostRemove(eventId string, sessionUser *comm.SessionUSER) error {
	if sessionUser.UserType != comm.Admin {
		return errors.New("非法操作")
	}
	var event model.Event
	has, err := e.mysql.Id(eventId).Get(&event)
	if err != nil {
		zap.S().Error(err)
		return ere.ErrorCommFindError
	}
	if !has {
		return errors.New("选课不存在")
	}
	event.Enable = comm.Disable
	event.Utime = time.Now()
	if _, err := e.mysql.Id(eventId).Cols("enable", "utime").Update(&event); err != nil {
		zap.S().Error(err)
		return ere.ErrorCommDeleteError
	}
	return nil
}

// 修改
func (e eventService) PostUpdate(arg *domain.EventUpdateArg, sessionUser *comm.SessionUSER) error {
	var event model.Event
	has, err := e.mysql.Id(arg.Id).Where("enable = ?", comm.Enable).Get(&event)
	if err != nil {
		zap.S().Error(err)
		return ere.ErrorCommFindError
	}
	if !has {
		return ere.ErrorCommNotFond
	}
	nowTime := time.Now()
	if nowTime.Before(event.Etime) && nowTime.After(event.Stime) {
		return errors.New("选课期间不允许修改")
	}
	event.Name = arg.Name
	event.Num = tools.StringToUint8(arg.Num)
	stime, err := time.ParseInLocation(comm.TimeFormatTime, arg.Stime, time.Local)
	if err != nil {
		zap.S().Error(err)
		return errors.New("格式化开始时间失败！")
	}
	event.Stime = stime
	etime, err := time.ParseInLocation(comm.TimeFormatTime, arg.Etime, time.Local)
	if err != nil {
		zap.S().Error(err)
		return errors.New("格式化结束时间失败！")
	}

	if stime.After(etime) {
		return errors.New("开始时间不能小于结束时间")
	}

	includeAll := false
	for _, s := range arg.SchoolYear {
		if s == "全部" {
			includeAll = true
		}
	}
	if includeAll {
		event.SchoolYear = "全部"
	} else {
		event.SchoolYear = strings.Join(arg.SchoolYear, ",")
	}
	event.TagIds = strings.Join(arg.TagIds, ",")
	event.Etime = etime
	event.CanUpdate = tools.StringToInt(arg.CanUpdate)
	event.Utime = time.Now()
	if _, err := e.mysql.Id(event.Id).Update(event); err != nil {
		zap.S().Error(err)
		return ere.ErrorCommUpdateError
	}
	e.updateCache(&event)
	return nil
}

// 分页
func (e eventService) PostPage(arg *comm.PageParam, sessionUser *comm.SessionUSER) (interface{}, error) {
	eventSql := e.mysql.Where("enable = ?", comm.Enable)
	countSql := e.mysql.Where("enable = ?", comm.Enable)

	var events []model.Event
	if err := eventSql.Limit(arg.GetLimit(), arg.GetOffset()).Desc("ctime", "id").Find(&events); err != nil {
		zap.S().Error(err)
		return nil, ere.ErrorCommFindError
	}

	var selectEvent model.Event
	total, err := countSql.Count(selectEvent)
	if err != nil {
		zap.S().Error(err)
		return nil, ere.ErrorCommFindError
	}

	var pageVo comm.PageVo
	pageVo.TotalCount = total

	eventVo := make([]domain.EventVo, 0)
	for _, event := range events {
		eventVo = append(eventVo, *domain.CvEventToVo(&event))
	}
	pageVo.Data = eventVo
	return &pageVo, nil
}

// 创建
func (e eventService) PostCreate(arg *domain.EventCreateArg, sessionUser *comm.SessionUSER) (*domain.EventVo, error) {
	var event model.Event
	event.Id = tools.SnowFlake.Generate().Int64()
	event.Name = arg.Name
	event.Num = tools.StringToUint8(arg.Num)
	stime, err := time.ParseInLocation(comm.TimeFormatTime, arg.Stime, time.Local)
	if err != nil {
		zap.S().Error(err)
		return nil, errors.New("格式化开始时间失败！")
	}
	event.Stime = stime
	etime, err := time.ParseInLocation(comm.TimeFormatTime, arg.Etime, time.Local)
	if err != nil {
		zap.S().Error(err)
		return nil, errors.New("格式化结束时间失败！")
	}

	if stime.After(etime) {
		return nil, errors.New("开始时间不能小于结束时间")
	}

	includeAll := false
	for _, s := range arg.SchoolYear {
		if s == "全部" {
			includeAll = true
		}
	}
	if includeAll {
		event.SchoolYear = "全部"
	} else {
		event.SchoolYear = strings.Join(arg.SchoolYear, ",")
	}
	event.TagIds = strings.Join(arg.TagIds, ",")
	event.Etime = etime
	event.Status = comm.EventStatusDisable
	event.Creator = sessionUser.Id
	event.Enable = comm.Enable
	event.CanUpdate = tools.StringToInt(arg.CanUpdate)
	event.Ctime = time.Now()

	if _, err := e.mysql.Insert(event); err != nil {
		zap.S().Error(err)
		return nil, ere.ErrorCommSaveError
	}
	return domain.CvEventToVo(&event), nil
}

// 检查是否存在正在进行中的选课
func (e eventService) existDoEvent() bool {
	var selectEvent model.Event
	count, err := e.mysql.Where("enable = ? and status = ? ", comm.Enable, comm.EventStatusEnable).Count(selectEvent)
	if err != nil {
		zap.S().Error(err)
		return true
	}
	return count > 0
}

// 删除缓存
func (e eventService) removeCache() {
	if _, err := db.RedisClient.Del(comm.EffectiveTakeLessonsActivityKey).Result(); err != nil {
		zap.S().Error(err)
	}
}

// 更新缓存
func (e eventService) updateCache(event *model.Event) {
	eventJsonByte, err := json.Marshal(*event)
	if err != nil {
		zap.S().Error(err)
		return
	}
	if _, err := db.RedisClient.Set(comm.EffectiveTakeLessonsActivityKey, eventJsonByte,
		time.Hour*24).Result(); err != nil {
		zap.S().Error(err)
	}
}

// 获取缓存
func (e eventService) getEventCache() *model.Event {
	if result, err := db.RedisClient.Get(comm.EffectiveTakeLessonsActivityKey).Result(); err == nil {
		event := new(model.Event)
		if err = json.Unmarshal([]byte(result), event); err != nil {
			zap.S().Error(err)
			return nil
		}
		return event
	}
	return nil
}

type EventService interface {
	// 创建
	PostCreate(arg *domain.EventCreateArg, sessionUser *comm.SessionUSER) (*domain.EventVo, error)

	// 分页
	PostPage(arg *comm.PageParam, sessionUser *comm.SessionUSER) (interface{}, error)

	// 修改
	PostUpdate(arg *domain.EventUpdateArg, sessionUser *comm.SessionUSER) error

	// 删除
	PostRemove(eventId string, sessionUser *comm.SessionUSER) error

	// 修改选课状态
	PostChange(eventId string, sessionUser *comm.SessionUSER) error

	// 查询存在的有效选课
	PostExistEvent(sessionUser *comm.SessionUSER) (*domain.EventVo, error)

	// 检查登录用户是否有选课
	isEffectiveEvent(event *model.Event, sessionUser *comm.SessionUSER) bool

	GetRedisGet(id string) interface{}
	GetRedisSet(key string, value string) interface{}
	GetRedisDel(key string) interface{}

	// 重新激活选课
	PostReactivation(id string, etime time.Time) error

	// 修改学生是否可以修改选课
	PostUpdateCanUpdate(arg domain.EventUpdateCanUpdateArg, user *comm.SessionUSER) error

	// // 学生查询历史
	PostStuHistoryEvent(user *comm.SessionUSER) ([]domain.EventSimpleVo, error)
}

func NewEventService(mysql *xorm.Engine) EventService {
	return &eventService{mysql: mysql}
}
