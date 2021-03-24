package service

import (
	"github.com/go-xorm/xorm"
	"go-take-lessons/db"
	"go-take-lessons/domain"
	"go-take-lessons/domain/comm"
	"go-take-lessons/domain/comm/ere"
	"go-take-lessons/model"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
)

type stuFocusService struct {
	mysql         *xorm.Engine
	configService *ConfigurationService
}

type StuFocusService interface {
	// 创建
	PostCreate(arg *domain.StuFocusCreateArg, sessionUser *comm.SessionUSER) error

	// 删除
	PostRemove(arg *domain.StuFocusCreateArg, sessionUser *comm.SessionUSER) error

	// 查询
	PostList(eventId string, sessionUser *comm.SessionUSER) *[]domain.ConfigurationSimpleVo
}

// 创建
func (s stuFocusService) PostCreate(arg *domain.StuFocusCreateArg, sessionUser *comm.SessionUSER) error {
	exist, err := s.exist(arg.EventId, arg.ConfigSubjectId, sessionUser.Id)
	if err != nil {
		zap.S().Error(err)
		return ere.ErrorCommFindError
	}

	if exist {
		return nil
	}

	stuFocus := domain.CvStuFocusCreateArgToModel(arg, sessionUser)
	if _, err := s.mysql.Insert(stuFocus); err != nil {
		zap.S().Error(err)
		return ere.ErrorCommSaveError
	}
	s.updateCacheStuFocusConfigSubjectIdsWhenExistKey(strconv.FormatInt(sessionUser.Id, 10),
		strconv.FormatInt(stuFocus.CsId, 10), true)
	return nil
}

// 删除
func (s stuFocusService) PostRemove(arg *domain.StuFocusCreateArg, sessionUser *comm.SessionUSER) error {
	focus := new(model.StuFocus)
	has, err := s.mysql.Where("user_id = ?", sessionUser.Id).And("event_id = ?", arg.EventId).And("cs_id = ?",
		arg.ConfigSubjectId).And("enable = ?",
		comm.Enable).Get(focus)

	if err != nil {
		zap.S().Error(err)
		return ere.ErrorCommDeleteError
	}
	if !has {
		return ere.ErrorCommNotFond
	}

	focus.Utime = time.Now()
	focus.Enable = comm.Disable
	if _, err := s.mysql.Id(focus.Id).Cols("enable", "utime").Update(focus); err != nil {
		zap.S().Error(err)
		return ere.ErrorCommDeleteError
	}
	s.updateCacheStuFocusConfigSubjectIdsWhenExistKey(strconv.FormatInt(sessionUser.Id, 10),
		strconv.FormatInt(focus.CsId, 10), false)
	return nil
}

// 查询
func (s stuFocusService) PostList(eventId string, sessionUser *comm.SessionUSER) *[]domain.ConfigurationSimpleVo {
	var configSubjectSimpleVos *[]domain.ConfigurationSimpleVo
	update := false
	stuFocusConfigSubjectIds := s.getCacheStuFocusConfigSubjectIds(strconv.FormatInt(sessionUser.Id, 10))
	if stuFocusConfigSubjectIds == nil || len(stuFocusConfigSubjectIds) == 0 {
		sql := "select s.cs_id from stu_focus s where s.event_id = ? and s.user_id = ? and s.`enable` = ?"
		if err := s.mysql.SQL(sql, eventId, sessionUser.Id, comm.Enable).Find(&stuFocusConfigSubjectIds); err != nil {
			zap.S().Error(err)
			return configSubjectSimpleVos
		}
		update = true
	}
	if stuFocusConfigSubjectIds == nil || len(stuFocusConfigSubjectIds) == 0 {
		return configSubjectSimpleVos
	}

	configService := *s.configService
	configSubjectSimpleVos = configService.PostListStu(eventId)
	if configSubjectSimpleVos == nil || len(stuFocusConfigSubjectIds) == 0 {
		return configSubjectSimpleVos
	}
	if update {
		s.loadCacheStuFocusConfigSubjectIds(strconv.FormatInt(sessionUser.Id, 10), stuFocusConfigSubjectIds)
	}
	newConfigSubjectSimpleVos := make([]domain.ConfigurationSimpleVo, 0)
	for _, configSubject := range *configSubjectSimpleVos {
		for _, stuFocusSubjectId := range stuFocusConfigSubjectIds {
			if strings.Compare(configSubject.SubjectId, stuFocusSubjectId) == 0 {
				newConfigSubjectSimpleVos = append(newConfigSubjectSimpleVos, configSubject)
			}
		}
	}
	return &newConfigSubjectSimpleVos
}

// 更新缓存
func (s stuFocusService) updateCacheStuFocusConfigSubjectIdsWhenExistKey(userId, stuFocusSubjectId string,
	append bool) {
	if result, err := db.RedisClient.Exists(comm.RedisStuFocusSelfConfigSubjectIdsKey + userId).Result(); err != nil || result <= 0 {
		return
	}
	if append {
		if _, err := db.RedisClient.LPush(comm.RedisStuFocusSelfConfigSubjectIdsKey+userId,
			stuFocusSubjectId).Result(); err != nil {
			zap.S().Error(err)
		}
	} else {
		if _, err := db.RedisClient.LRem(comm.RedisStuFocusSelfConfigSubjectIdsKey+userId, 1, stuFocusSubjectId).Result(); err != nil {
			zap.S().Error(err)
		}
	}
}

// 加载缓存
func (s stuFocusService) loadCacheStuFocusConfigSubjectIds(userId string, stuFocusSubjectIds []string) {
	if _, err := db.RedisClient.LPush(comm.RedisStuFocusSelfConfigSubjectIdsKey+userId,
		stuFocusSubjectIds).Result(); err != nil {
		zap.S().Error(err)
	}
}

//获取缓存
func (s stuFocusService) getCacheStuFocusConfigSubjectIds(userId string) []string {
	result, err := db.RedisClient.LRange(comm.RedisStuFocusSelfConfigSubjectIdsKey+userId, 0, -1).Result()
	if err != nil {
		zap.S().Error(err)
		return nil
	}
	return result
}

// 是否存在
func (s stuFocusService) exist(eventId, configSubjectId string, userId int64) (bool, error) {
	focus := new(model.StuFocus)
	i, e := s.mysql.Where("user_id = ?", userId).And("event_id = ?", eventId).And("cs_id = ?",
		configSubjectId).And("enable = ?",
		comm.Enable).Count(focus)
	if e != nil {
		zap.S().Error(e)
		return false, ere.ErrorCommFindError
	}
	return i > 0, nil
}

func NewStuFocusService(mysql *xorm.Engine, configService *ConfigurationService) StuFocusService {
	return &stuFocusService{
		mysql:         mysql,
		configService: configService,
	}
}
