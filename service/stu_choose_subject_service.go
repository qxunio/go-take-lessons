package service

import (
	"encoding/json"
	"errors"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/go-xorm/xorm"
	"go-take-lessons/db"
	"go-take-lessons/domain"
	"go-take-lessons/domain/comm"
	"go-take-lessons/model"
	"go-take-lessons/tools"
	"go-take-lessons/tools/excel"
	"go.uber.org/zap"
	"strconv"
	"sync"
	"time"
)

type stuChooseSubjectService struct {
	mysql        *xorm.Engine
	eventService *EventService
	csService    *ConfigurationService
}

// 管理员导出报表
func (s stuChooseSubjectService) PostExport(arg domain.IdArg, user *comm.SessionUSER) (*excelize.File, error) {
	head := []string{"主题", "学生课程限制", "创建时间", "开始时间", "结束时间"}
	subjectHead := []string{"姓名", "账号", "班级", "入学年", "选择时间"}
	file := excel.CreateExcelModel(head, "选课详情", "Sheet1")
	headStyle, _ := file.NewStyle(excel.GetHeadCellStyle())
	titleStyle, _ := file.NewStyle(excel.GetTitleCellStyle())
	contentStyle, _ := file.NewStyle(excel.GetContentCellStyle())

	var event model.Event
	if _, err := s.mysql.Where("enable = ? AND id = ?", comm.Enable, arg.Id).Get(&event); err != nil {
		zap.S().Error("管理员导出报表，查询选课设置错误,eventId: %v", arg.Id, err)
		return nil, errors.New("导出异常")
	}

	file.SetCellValue("Sheet1", "A3", event.Name)
	file.SetCellStyle("Sheet1", "A3", "A3", contentStyle)
	file.SetCellValue("Sheet1", "B3", event.Num)
	file.SetCellStyle("Sheet1", "B3", "B3", contentStyle)
	file.SetCellValue("Sheet1", "C3", event.Ctime.Format(comm.TimeFormatTime))
	file.SetCellStyle("Sheet1", "C3", "C3", contentStyle)
	file.SetCellValue("Sheet1", "D3", event.Stime.Format(comm.TimeFormatTime))
	file.SetCellStyle("Sheet1", "D3", "D3", contentStyle)
	file.SetCellValue("Sheet1", "E3", event.Etime.Format(comm.TimeFormatTime))
	file.SetCellStyle("Sheet1", "E3", "E3", contentStyle)

	cService := *s.csService
	details, err := cService.PostCsDetails(arg.Id, user)
	if err != nil {
		return nil, err
	}

	file.SetCellValue("Sheet1", "A5", "课堂详情")
	file.MergeCell("Sheet1", "A5", "G5")
	file.SetCellStyle("Sheet1", "A5", "G5", titleStyle)

	file.SetCellValue("Sheet1", "A6", "课堂名称")
	file.SetCellStyle("Sheet1", "A6", "A6", headStyle)
	file.SetCellValue("Sheet1", "B6", "学科")
	file.SetCellStyle("Sheet1", "B6", "B6", headStyle)
	file.SetCellValue("Sheet1", "C6", "教学地点")
	file.SetCellStyle("Sheet1", "C6", "C6", headStyle)
	file.SetCellValue("Sheet1", "D6", "教学时间")
	file.SetCellStyle("Sheet1", "D6", "D6", headStyle)
	file.SetCellValue("Sheet1", "E6", "教师")
	file.SetCellStyle("Sheet1", "E6", "E6", headStyle)
	file.SetCellValue("Sheet1", "F6", "人数限制")
	file.SetCellStyle("Sheet1", "F6", "F6", headStyle)
	file.SetCellValue("Sheet1", "G6", "已选人数")
	file.SetCellStyle("Sheet1", "G6", "G6", headStyle)

	for index, detail := range details {
		file.SetCellValue("Sheet1", "A"+strconv.FormatInt(int64(index+7), 10), detail.ClassName)
		file.SetCellStyle("Sheet1", "A"+strconv.FormatInt(int64(index+7), 10), "A"+strconv.FormatInt(int64(index+7), 10), contentStyle)
		file.SetCellValue("Sheet1", "B"+strconv.FormatInt(int64(index+7), 10), detail.SubjectName)
		file.SetCellStyle("Sheet1", "B"+strconv.FormatInt(int64(index+7), 10), "B"+strconv.FormatInt(int64(index+7), 10), contentStyle)
		file.SetCellValue("Sheet1", "C"+strconv.FormatInt(int64(index+7), 10), detail.TeachAddress)
		file.SetCellStyle("Sheet1", "C"+strconv.FormatInt(int64(index+7), 10), "C"+strconv.FormatInt(int64(index+7), 10), contentStyle)
		file.SetCellValue("Sheet1", "D"+strconv.FormatInt(int64(index+7), 10), detail.TeachTime)
		file.SetCellStyle("Sheet1", "D"+strconv.FormatInt(int64(index+7), 10), "D"+strconv.FormatInt(int64(index+7), 10), contentStyle)
		file.SetCellValue("Sheet1", "E"+strconv.FormatInt(int64(index+7), 10), detail.Teacher)
		file.SetCellStyle("Sheet1", "E"+strconv.FormatInt(int64(index+7), 10), "E"+strconv.FormatInt(int64(index+7), 10), contentStyle)
		file.SetCellValue("Sheet1", "F"+strconv.FormatInt(int64(index+7), 10), detail.Num)
		file.SetCellStyle("Sheet1", "F"+strconv.FormatInt(int64(index+7), 10), "F"+strconv.FormatInt(int64(index+7), 10), contentStyle)
		file.SetCellValue("Sheet1", "G"+strconv.FormatInt(int64(index+7), 10), detail.SelectedPlaces)
		file.SetCellStyle("Sheet1", "G"+strconv.FormatInt(int64(index+7), 10), "G"+strconv.FormatInt(int64(index+7), 10), contentStyle)

		selectSql := "SELECT ss.id,ss.user_id,u.`name`,u.account,ss.class,ss.school_year,ss.ctime FROM stu_subject ss LEFT JOIN `user` u ON ss.user_id = u.id WHERE ss.cs_id = " + detail.Id + " AND ss.`enable` = " + strconv.FormatUint(comm.Enable, 10)
		var do []model.StuSubjectDo
		if err := s.mysql.SQL(selectSql).Find(&do); err != nil {
			zap.S().Error("管理员导出报表，查询课程选课错误,csId: %v", detail.Id, err)
			return nil, errors.New("导出异常")
		}
		excel.CreateFileExcelModel(file, subjectHead, detail.SubjectName, detail.SubjectName)
		for i, subject := range do {
			subscript := strconv.FormatInt(int64(i+3), 10)
			file.SetCellValue(detail.SubjectName, "A"+subscript, subject.Name)
			file.SetCellStyle("Sheet1", "A"+subscript, "A"+subscript, contentStyle)
			file.SetCellValue(detail.SubjectName, "B"+subscript, subject.Account)
			file.SetCellStyle("Sheet1", "B"+subscript, "B"+subscript, contentStyle)
			file.SetCellValue(detail.SubjectName, "C"+subscript, subject.Class)
			file.SetCellStyle("Sheet1", "C"+subscript, "C"+subscript, contentStyle)
			file.SetCellValue(detail.SubjectName, "D"+subscript, subject.SchoolYear)
			file.SetCellStyle("Sheet1", "D"+subscript, "D"+subscript, contentStyle)
			file.SetCellValue(detail.SubjectName, "E"+subscript, subject.Ctime)
			file.SetCellStyle("Sheet1", "E"+subscript, "E"+subscript, contentStyle)
		}
	}
	file.SetActiveSheet(0)
	return file, nil
}

var stuChooseSubjectLock sync.Mutex

// 管理员删除学生选课
func (s stuChooseSubjectService) PostDeleteStu(arg domain.ConfigurationAdminAppendArg, user *comm.SessionUSER) error {
	existSql := "select count(1) from stu_subject where user_id = ? and cs_id = ? and enable = ? and event_id = ?"
	updateStuSql := "UPDATE stu_subject SET enable = ? WHERE user_id = ? AND event_id = ? AND cs_id = ?"
	updateCsSql := "UPDATE configuration_subject SET selected_places = selected_places -1 WHERE id = ? AND event_id = ?"
	stuChooseSubjectLock.Lock()
	defer stuChooseSubjectLock.Unlock()
	var countNUm int
	if _, err := s.mysql.SQL(existSql, arg.Uid, arg.CsId, comm.Enable, arg.EventId).Get(&countNUm); err != nil {
		zap.S().Error("管理员删除学生选课,查询学生是否已经选择异常,eventId: %v ,csId: %v", arg.EventId, arg.CsId, err)
		return errors.New("删除异常")
	}

	if countNUm == 0 {
		return errors.New("查询不到课程存在")
	}

	session := s.mysql.NewSession()
	err := session.Begin()
	// 更新课程人数
	sql := "update configuration_subject set selected_places = selected_places+1 where id = ? and event_id = ?"
	if _, err := s.mysql.Exec(sql, arg.CsId, arg.EventId); err != nil {
		zap.S().Error(err)
		if err = session.Rollback(); err != nil {
			zap.S().Error(err)
		}
		return errors.New("删除异常")
	}

	if _, err := s.mysql.Exec(updateStuSql, comm.Disable, arg.Uid, arg.EventId, arg.CsId); err != nil {
		zap.S().Error(err)
		if err = session.Rollback(); err != nil {
			zap.S().Error(err)
		}
		return errors.New("删除异常")
	}

	if _, err := s.mysql.Exec(updateCsSql, arg.CsId, arg.EventId); err != nil {
		zap.S().Error(err)
		if err = session.Rollback(); err != nil {
			zap.S().Error(err)
		}
		return errors.New("删除异常")
	}

	count, err := db.RedisClient.HGet(comm.RedisStuSelectConfigSubjectRemainingPlacesKey+arg.EventId, arg.CsId).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			count = "0"
		} else {
			zap.S().Error("管理员删除学生选课，查询redis错误 ", err)
			if err = session.Rollback(); err != nil {
				zap.S().Error(err)
			}
			return errors.New("替换异常")
		}
	}
	if count != "0" {
		num := tools.StringToInt(count) - 1
		if _, err := db.RedisClient.HSet(comm.RedisStuSelectConfigSubjectRemainingPlacesKey+arg.EventId, arg.CsId, num).Result(); err != nil {
			zap.S().Error(err)
			if err = session.Rollback(); err != nil {
				zap.S().Error(err)
			}
			return errors.New("删除异常")
		}
	}
	if err = session.Commit(); err != nil {
		zap.S().Error(err)
		return errors.New("删除异常")
	}

	return nil
}

// 管理员追加学生到课堂
func (s stuChooseSubjectService) PostAppendStu(arg domain.ConfigurationAdminAppendArg, user *comm.SessionUSER) error {
	var student model.User
	if _, err := s.mysql.Id(arg.Uid).Where("enable = ?", comm.Enable).Get(&student); err != nil {
		zap.S().Error("管理员追加学生到课堂，查询不到学生")
		return errors.New("查询不到有效的学生")
	}

	eventService := *s.eventService
	event, err := eventService.PostExistEvent(user)
	if err != nil {
		return err
	}
	if tools.IsBlank(event.Id) || event.Id != arg.EventId {
		zap.S().Error("管理员追加学生到课堂，查询不到有效选课事件")
		return errors.New("查询不到有效的选课")
	}

	result, err := db.RedisClient.HGet(comm.RedisStuSelectConfigSubjectListKey+event.Id, arg.CsId).Result()
	if err != nil {
		zap.S().Error("管理员追加学生到课堂，查询不到课程，课程id： %v", arg.CsId)
		return errors.New("追加选课失败，查询不到课程")
	}
	var configSimpleDos model.ConfigurationSimpleDo
	err = json.Unmarshal([]byte(result), &configSimpleDos)
	if err != nil {
		zap.S().Error("管理员追加学生到课堂，序列号课程失败，课程id : "+arg.CsId, err)
		return errors.New("追加选课失败")
	}

	if configSimpleDos.Id == 0 {
		zap.S().Error("管理员追加学生到课堂，查询不到课程，课程id : "+arg.CsId, err)
		return errors.New("追加选课失败，查询不到课程")
	}

	existSql := "select count(1) from stu_subject where user_id = ? and cs_id = ? and enable = ? and event_id = ?"
	var countNUm int
	if _, err := s.mysql.SQL(existSql, user.Id, arg.CsId, comm.Enable, arg.EventId).Get(&countNUm); err != nil {
		zap.S().Error("管理员追加学生到课堂,查询学生是否已经选择异常,eventId: %v ,csId: %v", arg.EventId, arg.CsId, err)
		return errors.New("选课异常")
	}

	if countNUm > 0 {
		return errors.New("你已经选择了这门课程")
	}

	if _, err := db.RedisClient.HIncrBy(comm.RedisStuSelectConfigSubjectRemainingPlacesKey+arg.EventId, arg.CsId, 1).Result(); err != nil {
		zap.S().Error("管理员追加学生到课堂，初始化redis已选课程人数异常,eventId: v% ,csId: v%", arg.EventId, arg.CsId, err)
		return errors.New("选课异常")
	}

	lock := s.buildDateLock(tools.StringToInt64(arg.EventId), tools.StringToInt64(arg.CsId), event.Name, student.Class, student.SchoolYear, student.Id)
	session := s.mysql.NewSession()
	err = session.Begin()
	// 记录
	if _, err := s.mysql.Insert(lock); err != nil {
		zap.S().Error("管理员追加学生到课堂,保存数据库异常,eventId: v% ,csId: v%", arg.EventId, arg.CsId, err)
		if err = session.Rollback(); err != nil {
			zap.S().Error(err)
		}
		return errors.New("选课异常")
	}

	// 更新课程人数
	sql := "update configuration_subject set selected_places = selected_places+1 where id = ? and event_id = ?"
	if _, err := s.mysql.Exec(sql, arg.CsId, arg.EventId); err != nil {
		zap.S().Error(err)
		if err = session.Rollback(); err != nil {
			zap.S().Error(err)
		}
		return errors.New("选课异常")
	}
	if err = session.Commit(); err != nil {
		zap.S().Error(err)
		return errors.New("选课异常")
	}
	return nil
}

// 管理员替换学生到课堂
func (s stuChooseSubjectService) PostReplaceStu(arg domain.ConfigurationAdminReplaceArg, user *comm.SessionUSER) error {
	var student model.User
	if _, err := s.mysql.Id(arg.Uid).Where("enable = ?", comm.Enable).Get(&student); err != nil {
		zap.S().Error("管理员替换学生到课堂，查询不到学生")
		return errors.New("查询不到有效的学生")
	}
	eventService := *s.eventService
	event, err := eventService.PostExistEvent(user)
	if err != nil {
		return err
	}
	if tools.IsBlank(event.Id) || event.Id != arg.EventId {
		zap.S().Error("管理员替换学生到课堂，查询不到有效选课事件")
		return errors.New("查询不到有效的选课")
	}

	result, err := db.RedisClient.HGet(comm.RedisStuSelectConfigSubjectListKey+event.Id, arg.CsId).Result()
	if err != err {
		zap.S().Error("管理员替换学生到课堂，查询不到课程，课程id： %v", arg.CsId)
		return errors.New("替换选课失败，查询不到课程")
	}
	var configSimpleDos model.ConfigurationSimpleDo
	err = json.Unmarshal([]byte(result), &configSimpleDos)
	if err != nil {
		zap.S().Error("管理员替换学生到课堂，序列号课程失败，课程id : "+arg.CsId, err)
		return errors.New("替换选课失败")
	}

	if configSimpleDos.Id == 0 {
		zap.S().Error("管理员替换学生到课堂，查询不到课程，课程id : "+arg.CsId, err)
		return errors.New("替换选课失败，查询不到课程")
	}

	existSql := "select count(1) from stu_subject where user_id = ? and cs_id = ? and enable = ? and event_id = ?"
	updateStuSql := "UPDATE stu_subject SET enable = ? WHERE user_id = ? AND event_id = ? AND cs_id = ?"
	updateCsSql := "UPDATE configuration_subject SET selected_places = selected_places -1 WHERE id = ? AND event_id = ?"

	stuChooseSubjectLock.Lock()
	defer stuChooseSubjectLock.Unlock()

	var countNUm int
	if _, err := s.mysql.SQL(existSql, user.Id, arg.CsId, comm.Enable, arg.EventId).Get(&countNUm); err != nil {
		zap.S().Error("管理员替换学生到课堂,查询学生是否已经选择异常,eventId: %v ,csId: %v", arg.EventId, arg.CsId, err)
		return errors.New("替换异常")
	}

	if countNUm > 0 {
		return errors.New("你已经选择了这门课程")
	}

	if _, err := db.RedisClient.HIncrBy(comm.RedisStuSelectConfigSubjectRemainingPlacesKey+arg.EventId, arg.CsId, 1).Result(); err != nil {
		zap.S().Error("管理员替换学生到课堂，初始化redis已选课程人数异常,eventId: v% ,csId: v%", arg.EventId, arg.CsId, err)
		return errors.New("替换异常")
	}

	lock := s.buildDateLock(tools.StringToInt64(arg.EventId), tools.StringToInt64(arg.CsId), event.Name, student.Class, student.SchoolYear, student.Id)
	session := s.mysql.NewSession()
	err = session.Begin()
	// 记录
	if _, err := s.mysql.Insert(lock); err != nil {
		zap.S().Error("管理员替换学生到课堂,保存数据库异常,eventId: v% ,csId: v%", arg.EventId, arg.CsId, err)
		if err = session.Rollback(); err != nil {
			zap.S().Error(err)
		}
		return errors.New("替换异常")
	}

	// 更新课程人数
	sql := "update configuration_subject set selected_places = selected_places+1 where id = ? and event_id = ?"
	if _, err := s.mysql.Exec(sql, arg.CsId, arg.EventId); err != nil {
		zap.S().Error(err)
		if err = session.Rollback(); err != nil {
			zap.S().Error(err)
		}
		return errors.New("替换异常")
	}

	if _, err := s.mysql.Exec(updateStuSql, comm.Disable, arg.Uid, arg.EventId, arg.ReplaceCsId); err != nil {
		zap.S().Error(err)
		if err = session.Rollback(); err != nil {
			zap.S().Error(err)
		}
		return errors.New("替换异常")
	}

	if _, err := s.mysql.Exec(updateCsSql, arg.ReplaceCsId, arg.EventId); err != nil {
		zap.S().Error(err)
		if err = session.Rollback(); err != nil {
			zap.S().Error(err)
		}
		return errors.New("替换异常")
	}

	count, err := db.RedisClient.HGet(comm.RedisStuSelectConfigSubjectRemainingPlacesKey+arg.EventId, arg.ReplaceCsId).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			count = "0"
		} else {
			zap.S().Error("管理员替换学生到课堂，查询redis错误 ", err)
			if err = session.Rollback(); err != nil {
				zap.S().Error(err)
			}
			return errors.New("替换异常")
		}
	}
	if count != "0" {
		num := tools.StringToInt(count) - 1
		if _, err := db.RedisClient.HSet(comm.RedisStuSelectConfigSubjectRemainingPlacesKey+arg.EventId, arg.CsId, num).Result(); err != nil {
			zap.S().Error(err)
			if err = session.Rollback(); err != nil {
				zap.S().Error(err)
			}
			return errors.New("替换异常")
		}
	}
	if err = session.Commit(); err != nil {
		zap.S().Error(err)
		return errors.New("替换异常")
	}
	return nil
}

// 学生删除课程
func (s stuChooseSubjectService) PostDelete(arg domain.ConfigurationEventIdAndCsIdArg, user *comm.SessionUSER) error {
	eventService := *s.eventService
	event, err := eventService.PostExistEvent(user)
	if err != nil {
		return err
	}
	if tools.IsBlank(event.Id) || event.Id != arg.EventId {
		zap.S().Error("学生修改选课，查询不到有效选课事件,cs_id:%v,event_id:%v,userId: %v", arg.CsId, arg.EventId, user.Id)
		return errors.New("查询不到有效的选课")
	}

	if tools.StringToInt(event.CanUpdate) == comm.Disable {
		zap.S().Error("学生修改选课，不允许修改课程,cs_id:%v,event_id:%v,userId: %v", arg.CsId, arg.EventId, user.Id)
		return errors.New("当前选课不允许删除")
	}

	updateStuSql := "UPDATE stu_subject SET enable = ? WHERE user_id = ? AND event_id = ? AND cs_id = ?"
	updateCsSql := "UPDATE configuration_subject SET selected_places = selected_places-1 WHERE id = ? AND event_id = ?"

	stuChooseSubjectLock.Lock()
	defer stuChooseSubjectLock.Unlock()

	session := s.mysql.NewSession()
	err = session.Begin()
	if _, err := s.mysql.Exec(updateStuSql, comm.Disable, user.Id, arg.EventId, arg.CsId); err != nil {
		zap.S().Error(err)
		if err = session.Rollback(); err != nil {
			zap.S().Error(err)
		}
		return errors.New("删除异常")
	}

	if _, err := s.mysql.Exec(updateCsSql, arg.CsId, arg.EventId); err != nil {
		zap.S().Error(err)
		if err = session.Rollback(); err != nil {
			zap.S().Error(err)
		}
		return errors.New("删除异常")
	}

	count, err := db.RedisClient.HGet(comm.RedisStuSelectConfigSubjectRemainingPlacesKey+arg.EventId, arg.CsId).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			count = "0"
		} else {
			zap.S().Error("学生选课，查询redis错误 ", err)
			if err = session.Rollback(); err != nil {
				zap.S().Error(err)
			}
			return errors.New("删除异常")
		}
	}

	if count != "0" {
		num := tools.StringToInt(count) - 1
		if _, err := db.RedisClient.HSet(comm.RedisStuSelectConfigSubjectRemainingPlacesKey+arg.EventId, arg.CsId, num).Result(); err != nil {
			zap.S().Error(err)
			if err = session.Rollback(); err != nil {
				zap.S().Error(err)
			}
			return errors.New("删除异常")
		}
	}

	if err = session.Commit(); err != nil {
		zap.S().Error(err)
		return errors.New("删除异常")
	}
	return nil
}

// 学生选课
func (s stuChooseSubjectService) PostLock(eventId, csId string, user *comm.SessionUSER) error {
	eventService := *s.eventService
	event, err := eventService.PostExistEvent(user)
	if err != nil {
		return err
	}
	if tools.IsBlank(event.Id) || event.Id != eventId {
		zap.S().Error("学生选课，查询不到有效选课事件,userId: %v", user.Id)
		return errors.New("查询不到有效的选课")
	}

	result, err := db.RedisClient.HGet(comm.RedisStuSelectConfigSubjectListKey+event.Id, csId).Result()
	if err != err {
		zap.S().Error("学生选课，查询不到课程，课程id： %v ，学生id : %v", csId, user.Id)
		return errors.New("选课失败")
	}
	var configSimpleDos model.ConfigurationSimpleDo
	err = json.Unmarshal([]byte(result), &configSimpleDos)
	if err != nil {
		zap.S().Error("学生选课，序列号课程失败，课程id : "+csId+"学生id : "+strconv.FormatInt(user.Id, 10), err)
		return errors.New("选课失败")
	}

	if configSimpleDos.Id == 0 {
		zap.S().Error("学生选课，查询不到课程，课程id : "+csId+"学生id : "+strconv.FormatInt(user.Id, 10), err)
		return errors.New("课程不存在")
	}

	stuChooseSubjectLock.Lock()
	defer stuChooseSubjectLock.Unlock()

	count, err := db.RedisClient.HGet(comm.RedisStuSelectConfigSubjectRemainingPlacesKey+eventId, csId).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			count = "0"
		} else {
			zap.S().Error("学生选课，查询redis错误 ", err)
			return errors.New("选课失败")
		}
	}
	lock := s.buildDateLock(tools.StringToInt64(eventId), tools.StringToInt64(csId), event.Name, user.Class, user.SchoolYear, user.Id)

	if tools.IsBlank(count) || tools.StringToInt(count) < configSimpleDos.Num {
		existSql := "select cs_id from stu_subject where user_id = ? and  enable = ? and event_id = ?"
		var csIds []int64
		if err := s.mysql.SQL(existSql, user.Id, comm.Enable, eventId).Find(&csIds); err != nil {
			zap.S().Error("学生选课,查询自己是否已经选择异常,eventId: %v ,csId: %v", eventId, csId, err)
			return errors.New("选课异常")
		}

		if len(csIds) >= tools.StringToInt(event.Num) {
			return errors.New("当前选课已经选满啦")
		}

		for _, id := range csIds {
			if strconv.FormatInt(id, 10) == csId {
				return errors.New("你已经选择了这门课程")
			}
		}

		if _, err := db.RedisClient.HIncrBy(comm.RedisStuSelectConfigSubjectRemainingPlacesKey+eventId, csId, 1).Result(); err != nil {
			zap.S().Error("学生选课，初始化redis已选课程人数异常,eventId: v% ,csId: v%", eventId, csId, err)
			return errors.New("选课异常")
		}

		session := s.mysql.NewSession()
		err := session.Begin()
		// 记录
		if _, err := s.mysql.Insert(lock); err != nil {
			zap.S().Error("学生选课,保存数据库异常,eventId: v% ,csId: v%", eventId, csId, err)
			if err = session.Rollback(); err != nil {
				zap.S().Error(err)
			}
			return errors.New("选课异常")
		}

		// 更新课程人数
		sql := "update configuration_subject set selected_places = selected_places+1 where id = ? and event_id = ?"
		if _, err := s.mysql.Exec(sql, csId, eventId); err != nil {
			zap.S().Error(err)
			if err = session.Rollback(); err != nil {
				zap.S().Error(err)
			}
			return errors.New("选课异常")
		}
		if err = session.Commit(); err != nil {
			zap.S().Error(err)
			return errors.New("选课异常")
		}
		return nil

	}
	return errors.New("课程已经选完啦")
}

// 组装数据
func (s stuChooseSubjectService) buildDateLock(eventId, csId int64, eventName string, userClass int, userSchoolYear int, userId int64) model.StuSubject {
	var ss model.StuSubject
	ss.Id = tools.SnowFlake.Generate().Int64()
	ss.CsId = csId
	ss.EventId = eventId
	ss.Ctime = time.Now()
	ss.Enable = comm.Enable
	ss.EventName = eventName
	ss.SchoolYear = userSchoolYear
	ss.Class = userClass
	ss.UserId = userId
	return ss
}

type StuChooseSubjectService interface {
	// 学生选课
	PostLock(eventId, csId string, user *comm.SessionUSER) error

	// 学生删除课程
	PostDelete(arg domain.ConfigurationEventIdAndCsIdArg, user *comm.SessionUSER) error

	// 管理员追加学生到课堂
	PostAppendStu(arg domain.ConfigurationAdminAppendArg, user *comm.SessionUSER) error

	// 管理员替换学生到课堂
	PostReplaceStu(arg domain.ConfigurationAdminReplaceArg, user *comm.SessionUSER) error

	// 管理员删除学生选课
	PostDeleteStu(arg domain.ConfigurationAdminAppendArg, user *comm.SessionUSER) error

	// 管理员导出报表
	PostExport(arg domain.IdArg, user *comm.SessionUSER) (*excelize.File, error)
}

func NewStuChooseSubjectService(mysql *xorm.Engine, e *EventService, c *ConfigurationService) StuChooseSubjectService {
	return &stuChooseSubjectService{mysql: mysql, eventService: e, csService: c}
}
