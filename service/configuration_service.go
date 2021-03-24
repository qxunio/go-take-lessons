package service

import (
	"bytes"
	"encoding/json"
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

type configurationService struct {
	userService *UserService
	mysql       *xorm.Engine
}

// 教师追加学生到课堂 搜索学生
func (c configurationService) PostAppendStuSearch(arg domain.ConfigurationEventIdAndDestArg, user *comm.SessionUSER) (*[]domain.StudentSubjectSelectInfoVo, error) {
	selectArg := "%" + arg.Dest + "%"
	selectStuSql := "SELECT id,name,account,class,school_year FROM `user` WHERE type = ? AND enable = ? AND (name LIKE ? OR account LIKE ?)"

	var stu []model.StuDo

	if err := c.mysql.SQL(selectStuSql, comm.Student, comm.Enable, selectArg, selectArg).Find(&stu); err != nil {
		zap.S().Error(err)
		return nil, ere.ErrorCommFindError
	}

	var vos []domain.StudentSubjectSelectInfoVo
	if len(stu) != 0 {
		var selectInfo []model.StuSubjectAndBaseDo
		var studentIds []int64
		for _, do := range stu {
			studentIds = append(studentIds, do.Id)
		}

		if err := c.mysql.Table("stu_subject").
			Select("stu_subject.user_id,stu_subject.cs_id,stu_subject.ctime,configuration_subject.class_name,configuration_subject.subject_name").
			Join("LEFT", "configuration_subject", "configuration_subject.id = stu_subject.cs_id").
			Where("stu_subject.event_id = ? AND stu_subject.`enable` = ?", arg.EventId, comm.Enable).In("stu_subject.user_id", studentIds).
			Find(&selectInfo); err != nil {
			zap.S().Error(err)
			return nil, ere.ErrorCommFindError
		}

		for _, do := range stu {
			var vo domain.StudentSubjectSelectInfoVo
			vo.UserId = strconv.FormatInt(do.Id, 10)
			vo.Class = strconv.FormatInt(int64(do.Class), 10)
			vo.SchoolYear = strconv.FormatInt(int64(do.SchoolYear), 10)
			vo.Name = do.Name
			vo.Account = do.Account
			for _, baseDo := range selectInfo {
				if do.Id == baseDo.UserId {
					vo.CsId = strconv.FormatInt(baseDo.CsId, 10)
					vo.SelectTime = baseDo.Ctime.Format(comm.TimeFormatTime)
					vo.ClassName = baseDo.ClassName
					vo.SubjectName = baseDo.SubjectName
				}
			}
			vos = append(vos, vo)
		}
	}

	return &vos, nil
}

// 查询选课下的课程详情的学生列表
func (c configurationService) PostPageCsDetailsStu(arg domain.StuSubjectPageArg, user *comm.SessionUSER) (interface{}, error) {
	selectSql := "SELECT ss.id,ss.user_id,u.`name`,u.account,ss.class,ss.school_year,ss.ctime FROM stu_subject ss LEFT JOIN `user` u ON ss.user_id = u.id WHERE ss.cs_id = " + arg.CsId + " AND ss.`enable` = " + strconv.FormatUint(comm.Enable, 10)
	countSql := "SELECT count(1) FROM stu_subject ss LEFT JOIN `user` u ON ss.user_id = u.id WHERE ss.cs_id = " + arg.CsId + " AND ss.`enable` = " + strconv.FormatUint(comm.Enable, 10)

	if !tools.IsBlank(arg.Dest) {
		arg := "'%" + arg.Dest + "%'"
		sql2 := " AND u.account like " + arg + " or u.name like" + arg
		selectSql += sql2
		countSql += sql2
	}

	if !tools.IsBlank(arg.Class) {
		sql3 := " AND ss.class = '" + arg.Class + "'"
		selectSql += sql3
		countSql += sql3
	}

	limitSql := " ORDER BY ss.ctime DESC limit " + strconv.FormatUint(uint64(arg.GetOffset()), 10) + "," + strconv.FormatUint(uint64(arg.GetLimit()), 10)
	selectSql += limitSql

	var do []model.StuSubjectDo
	if err := c.mysql.SQL(selectSql).Find(&do); err != nil {
		zap.S().Error(err)
		return nil, ere.ErrorCommFindError
	}

	var um model.StuSubject
	var total int64
	total, err := c.mysql.SQL(countSql).Count(&um)
	if err != nil {
		zap.S().Error(err)
		return nil, ere.ErrorCommFindError
	}
	var pageVo comm.PageVo
	pageVo.TotalCount = total

	vos := make([]domain.StuSubjectVo, 0)

	for _, subjectDo := range do {
		vos = append(vos, domain.CvStuSubjectDoToStuSubjectVo(subjectDo))
	}
	pageVo.Data = vos
	return pageVo, nil
}

// 查询选课下的课程详情
func (c configurationService) PostCsDetails(eventId string, user *comm.SessionUSER) ([]domain.ConfigurationDetailsVo, error) {
	var subjects []model.ConfigurationSubject
	var teachers []model.ConfigurationTeacher

	if err := c.mysql.Where("enable = ?  and event_id = ?", comm.Enable, eventId).Find(&subjects); err != nil {
		zap.S().Error("查询选课下的课程详情错误， event_id : %v ", eventId, err)
		return nil, ere.ErrorCommFindError
	}
	var csIds []int64
	for _, s := range subjects {
		csIds = append(csIds, s.Id)
	}
	if err := c.mysql.Where("enable = ?", comm.Enable).In("cs_id", csIds).Find(&teachers); err != nil {
		zap.S().Error("查询选课下的课程详情下的教师错误， event_id : %v ", eventId, err)
		return nil, ere.ErrorCommFindError
	}

	var vos []domain.ConfigurationDetailsVo
	for _, sub := range subjects {
		var tens []string
		for _, teacher := range teachers {
			if teacher.CsId == sub.Id {
				tens = append(tens, teacher.TeacherName)
			}
		}
		teacher := ""
		if tens != nil {
			teacher = strings.Join(tens, ",")
		}
		vos = append(vos, domain.CvConfigurationToConfigurationDetailsVo(sub, teacher))
	}
	return vos, nil
}

// 学生查询自己选课记录
func (c configurationService) PostStuHistoryStuSelected(eventId string, user *comm.SessionUSER) (*[]domain.StudentSubjectSelectVo, error) {
	var studentSubject []model.StuSubject

	if err := c.mysql.Where("user_id = ? and event_id = ? and enable = ?", user.Id, eventId, comm.Enable).Find(&studentSubject); err != nil {
		zap.S().Error("学生查询自己选课记录错误， event_id : %v , user_id : %v", eventId, user.Id, err)
		return nil, ere.ErrorCommFindError
	}

	if studentSubject == nil {
		return nil, nil
	}
	var vos []domain.StudentSubjectSelectVo

	var csIds []int64
	for _, stuSubject := range studentSubject {
		csIds = append(csIds, stuSubject.CsId)
	}

	var cs []model.ConfigurationSubject
	if err := c.mysql.Where("enable = ?", comm.Enable).In("id", csIds).Find(&cs); err != nil {
		zap.S().Error("学生查询自己选课记录错误， event_id : %v , user_id : %v", eventId, user.Id, err)
		return nil, ere.ErrorCommFindError
	}
	var ts []model.ConfigurationTeacher
	if err := c.mysql.Where("enable = ?", comm.Enable).In("cs_id", csIds).Find(&ts); err != nil {
		zap.S().Error("学生查询自己选课记录错误， event_id : %v , user_id : %v", eventId, user.Id, err)
		return nil, ere.ErrorCommFindError
	}

	for _, stuSub := range studentSubject {
		for _, subject := range cs {
			if stuSub.CsId == subject.Id {
				var vo domain.StudentSubjectSelectVo
				vo.Id = strconv.FormatInt(stuSub.CsId, 10)
				vo.Class = strconv.FormatInt(int64(stuSub.Class), 10)
				vo.SchoolYear = strconv.FormatInt(int64(stuSub.SchoolYear), 10)
				vo.SelectTime = stuSub.Ctime.Format(comm.TimeFormatTime)
				vo.SubjectName = subject.SubjectName
				vo.ClassName = subject.ClassName
				vo.Introduction = subject.Introduction
				vo.TeachTime = subject.TeachTime
				vo.TeachAddress = subject.TeachAddress

				var tens []string
				for _, teacher := range ts {
					tens = append(tens, teacher.TeacherName)
				}
				if tens != nil {
					vo.TeacherName = strings.Join(tens, ",")
				}
				vos = append(vos, vo)
			}
		}
	}
	return &vos, nil
}

// 批量创建
func (c configurationService) PostCreateBatch(arg domain.ConfigurationCreateBatchArg, user *comm.SessionUSER) error {
	var accounts []string
	for _, subject := range arg.Subject {
		for _, teacher := range subject.Teacher {
			accounts = append(accounts, teacher)
		}
	}

	var us *[]model.User
	var err error
	userService := *c.userService
	if len(accounts) != 0 {
		us, err = userService.ListByUserAccount(accounts)
		if err != nil {
			return err
		}
	}

	var newSubjects []*model.ConfigurationSubject
	var newTeachers []*model.ConfigurationTeacher
	var newSimpleDo []model.ConfigurationSimpleDo

	now := time.Now()

	for _, subject := range arg.Subject {
		var cSubject model.ConfigurationSubject
		cSubject.Id = tools.SnowFlake.Generate().Int64()
		cSubject.Enable = comm.Enable
		cSubject.Utime = now
		cSubject.Ctime = now
		cSubject.ClassName = subject.ClassName
		cSubject.SubjectId = tools.StringToInt64(subject.SubjectId)
		cSubject.SubjectName = subject.SubjectName
		cSubject.Introduction = c.getIntroduction(subject.SubjectId)
		cSubject.SelectedPlaces = 0
		cSubject.EventId = tools.StringToInt64(arg.EventId)
		cSubject.Num = tools.StringToUint8(subject.Limit)
		cSubject.TeachAddress = arg.Address
		cSubject.TeachTime = arg.Time
		newSubjects = append(newSubjects, &cSubject)

		var do model.ConfigurationSimpleDo
		do.Id = cSubject.Id
		do.SubjectId = cSubject.SubjectId
		do.SubjectName = cSubject.SubjectName
		do.ClassName = subject.ClassName
		do.Introduction = cSubject.Introduction
		do.Num = int(cSubject.Num)
		do.TeachAddress = cSubject.TeachAddress
		do.TeachTime = cSubject.TeachTime

		var teacherName []string
		if subject.Teacher != nil && len(subject.Teacher) > 0 {
			for _, t := range subject.Teacher {
				for _, u := range *us {
					if u.Account == t {
						var cTeacher model.ConfigurationTeacher
						cTeacher.Id = tools.SnowFlake.Generate().Int64()
						cTeacher.EventId = tools.StringToInt64(arg.EventId)
						cTeacher.Enable = comm.Enable
						cTeacher.Ctime = now
						cTeacher.Utime = now
						cTeacher.CsId = cSubject.Id
						cTeacher.TeacherAccount = u.Account
						cTeacher.TeacherId = u.Id
						cTeacher.TeacherName = u.Name
						newTeachers = append(newTeachers, &cTeacher)
						teacherName = append(teacherName, u.Name)
						break
					}
				}
			}
			do.Teacher = strings.Join(teacherName, ",")
			newSimpleDo = append(newSimpleDo, do)
		}
	}

	session := c.mysql.NewSession()
	defer session.Close()
	err = session.Begin()
	if _, err = session.Insert(newSubjects); err != nil {
		zap.S().Error(err)
		_ = session.Rollback()
		return ere.ErrorCommSaveError
	}
	if _, err = session.Insert(newTeachers); err != nil {
		zap.S().Error(err)
		_ = session.Rollback()
		return ere.ErrorCommSaveError
	}
	if err = session.Commit(); err != nil {
		zap.S().Error(err)
		return ere.ErrorCommSaveError
	}
	c.updateCacheConfigSimpleDo(arg.EventId, &newSimpleDo)
	return nil
}

var updateConfigSimpleCacheLock sync.Mutex

// 根据选课事件ID查询选课课程 学生
func (c configurationService) PostListStu(eventId string) *[]domain.ConfigurationSimpleVo {
	var configSimpleDO *[]model.ConfigurationSimpleDo
	configSimpleDO = c.stuSelectConfigSubject(eventId)
	if *configSimpleDO != nil {
		c.setCacheConfigSimpleDo(eventId, configSimpleDO)
	}
	// 先从缓存读不变的数据
	configSimpleDO = c.getCacheStuSelectConfigSubject(eventId)
	if *configSimpleDO == nil {
		updateConfigSimpleCacheLock.Lock()
		// 这里从缓存拿是，防止下一个线程去查询数据表并更新缓存
		configSimpleDO = c.getCacheStuSelectConfigSubject(eventId)
		if *configSimpleDO == nil {
			// 缓存没有，再去数据库查并同步缓存
			configSimpleDO = c.stuSelectConfigSubject(eventId)
			if *configSimpleDO != nil {
				c.setCacheConfigSimpleDo(eventId, configSimpleDO)
			}
		}
		defer updateConfigSimpleCacheLock.Unlock()
	}
	if *configSimpleDO == nil {
		return nil
	}
	//处理剩余的名额
	remainingPlaces := c.getCacheStuSelectConfigSubjectRemainingPlaces(eventId)
	// TODO 排序
	return domain.CvConfigSimpleDOToVo(configSimpleDO, remainingPlaces)
}

// 删除选课教师配置
func (c configurationService) PostRemoveTeacher(configTeacherId string, sessionUser *comm.SessionUSER) error {
	var configTeacher model.ConfigurationTeacher

	_, err := c.mysql.Where("id = ? and enable = ?", configTeacherId, comm.Enable).Get(&configTeacher)
	if err != nil {
		zap.S().Error("选课配置SERVICE,删除选课教师配置 查询教师配置错误", err)
		return ere.ErrorCommDeleteError
	}
	if configTeacher.Id == 0 {
		zap.S().Error("选课配置SERVICE,删除选课教师配置 查询教师配置不存在", err)
		return ere.ErrorCommDeleteError
	}
	eventId := strconv.FormatInt(configTeacher.EventId, 10)
	csId := strconv.FormatInt(configTeacher.CsId, 10)
	do := c.GetCacheCSubjectByEventAndCsId(eventId, csId)

	configTeacher.Id = tools.StringToInt64(configTeacherId)
	configTeacher.Enable = comm.Disable
	configTeacher.Utime = time.Now()
	if _, err := c.mysql.Id(configTeacherId).Cols("enable", "utime").Update(&configTeacher); err != nil {
		zap.S().Error("选课配置SERVICE,删除选课教师配置错误", err)
		return ere.ErrorCommDeleteError
	}

	if do != nil {
		// 更新
		var teachers []model.ConfigurationTeacher
		if err := c.mysql.Where("cs_id = ? and enable = ?", csId, comm.Enable).Find(&teachers); err != nil {
			zap.S().Error("选课配置SERVICE,删除选课教师配置 同步缓存 查询教师配置错误", err)
		}
		var teacherName []string
		for _, teacher := range teachers {
			teacherName = append(teacherName, teacher.TeacherName)
		}
		do.Teacher = strings.Join(teacherName, ",")
		var arr []model.ConfigurationSimpleDo
		arr = append(arr, *do)
		c.updateCacheConfigSimpleDo(eventId, &arr)
	}
	return nil
}

// 删除选课课程,级联删除该课程下的教师
func (c configurationService) PostRemove(eventId, configSubjectId string, sessionUser *comm.SessionUSER) error {
	var configSubject model.ConfigurationSubject
	configSubject.Id = tools.StringToInt64(configSubjectId)
	configSubject.Enable = comm.Disable
	configSubject.Utime = time.Now()

	var configTeacher model.ConfigurationTeacher
	configTeacher.CsId = tools.StringToInt64(configSubjectId)
	configTeacher.Enable = comm.Disable
	configTeacher.Utime = time.Now()

	session := c.mysql.NewSession()
	defer session.Close()
	err := session.Begin()
	if _, err = session.ID(configSubject.Id).Cols("enable", "utime").Update(&configSubject); err != nil {
		_ = session.Rollback()
		zap.S().Error("选课配置SERVICE,删除选课课程错误，回滚数据", err)
		return ere.ErrorCommUpdateError
	}
	if _, err = session.Cols("enable", "utime").Where("cs_id = ?", configSubject.Id).Update(&configTeacher); err != nil {
		_ = session.Rollback()
		zap.S().Error("选课配置SERVICE,删除选课课程错误，回滚数据", err)
		return ere.ErrorCommUpdateError
	}
	if err = session.Commit(); err != nil {
		zap.S().Error("选课配置SERVICE,删除选课课程错误，提交事务失败", err)
		return ere.ErrorCommDeleteError
	}
	c.deleteCacheConfigSimpleDo(eventId, configSubjectId)
	return nil
}

// 更新
func (c configurationService) PostUpdate(arg *domain.ConfigurationUpdateArg) error {
	configSubject, configTeachers := domain.CvConfigurationUpdateArgToModel(arg, c.getIntroduction(arg.SubjectId))

	session := c.mysql.NewSession()
	defer session.Close()
	err := session.Begin()
	if _, err = session.ID(configSubject.Id).Update(configSubject); err != nil {
		zap.S().Error(err)
		_ = session.Rollback()
		return ere.ErrorCommSaveError
	}
	sql := c.sqlUpdateBatchConfigTeachers(*configTeachers)
	if !tools.IsBlank(sql) {
		if _, err = session.Exec(sql); err != nil {
			zap.S().Error(err)
			_ = session.Rollback()
			return ere.ErrorCommSaveError
		}
	}
	if err = session.Commit(); err != nil {
		zap.S().Error(err)
		return ere.ErrorCommSaveError
	}
	// 同步更新缓存
	c.buildDataUpdateCache(configSubject, configTeachers)
	return nil
}

// 创建
func (c configurationService) PostCreate(arg *domain.ConfigurationCreateArg) error {
	configSubject, configTeachers := domain.CvConfigurationCreateArgToModel(arg, c.getIntroduction(arg.SubjectId))

	session := c.mysql.NewSession()
	defer session.Close()
	err := session.Begin()
	if _, err = session.Insert(configSubject); err != nil {
		zap.S().Error(err)
		_ = session.Rollback()
		return ere.ErrorCommSaveError
	}
	if _, err = session.Insert(configTeachers); err != nil {
		zap.S().Error(err)
		_ = session.Rollback()
		return ere.ErrorCommSaveError
	}
	if err = session.Commit(); err != nil {
		zap.S().Error(err)
		return ere.ErrorCommSaveError
	}
	// 同步更新缓存
	c.buildDataUpdateCache(configSubject, configTeachers)
	return nil
}

// 根据选课事件ID查询所有的选课课程
func (c configurationService) PostList(eventId string, sessionUser *comm.SessionUSER) (*[]domain.ConfigurationVo, error) {
	var configSubjects []model.ConfigurationSubject
	err := c.mysql.Where("enable = ? and event_id = ?", comm.Enable, eventId).Find(&configSubjects)
	if err != nil {
		zap.S().Error(err)
		return nil, ere.ErrorCommFindError
	}

	var configTeachers []model.ConfigurationTeacher
	err = c.mysql.Where("enable = ? and event_id = ?", comm.Enable, eventId).Find(&configTeachers)
	if err != nil {
		zap.S().Error(err)
		return nil, ere.ErrorCommFindError
	}

	var configListVo []domain.ConfigurationVo
	cutMap, err := tools.SliceCut(configTeachers, "CsId")
	if err != nil {
		zap.S().Error(err)
		return nil, ere.ErrorCommFindError
	}
	for _, configSubject := range configSubjects {
		var configVo domain.ConfigurationVo
		configVo.Id = strconv.FormatInt(configSubject.Id, 10)
		configVo.EventId = strconv.FormatInt(configSubject.EventId, 10)
		configVo.SubjectId = strconv.FormatInt(configSubject.SubjectId, 10)
		configVo.SubjectName = configSubject.SubjectName
		configVo.ClassName = configSubject.ClassName
		configVo.Num = strconv.FormatUint(uint64(configSubject.Num), 10)
		configVo.TeachAddress = configSubject.TeachAddress
		configVo.TeachTime = configSubject.TeachTime
		configVo.SelectedPlaces = strconv.FormatUint(uint64(configSubject.SelectedPlaces), 10)
		specialConfigTeachers := cutMap[configVo.Id]
		if specialConfigTeachers != nil {
			var configTeachersVo []domain.ConfigurationTeacherVo
			for _, specialConfigTeacher := range specialConfigTeachers {
				var configTeacherVo domain.ConfigurationTeacherVo
				configTeacherVo.Id = specialConfigTeacher["Id"]
				configTeacherVo.EventId = specialConfigTeacher["EventId"]
				configTeacherVo.TeacherId = specialConfigTeacher["TeacherId"]
				configTeacherVo.TeacherAccount = specialConfigTeacher["TeacherAccount"]
				configTeacherVo.TeacherName = specialConfigTeacher["TeacherName"]
				configTeacherVo.Ctime = specialConfigTeacher["Ctime"]
				configTeacherVo.CsId = configVo.Id
				configTeachersVo = append(configTeachersVo, configTeacherVo)
			}
			configVo.Teachers = configTeachersVo
		}
		configListVo = append(configListVo, configVo)
	}
	return &configListVo, nil
}

// 更新选课教师的SQL
func (c configurationService) sqlUpdateBatchConfigTeachers(teachers []model.ConfigurationTeacher) string {
	lang := len(teachers)
	if lang == 0 {
		return ""
	}
	var sql bytes.Buffer
	sql.WriteString("replace into configuration_teacher (id,event_id,cs_id,teacher_id,teacher_name,teacher_account,enable,ctime,utime) values")
	for index, teacher := range teachers {
		sql.WriteString("(")
		sql.WriteString("'" + strconv.FormatInt(teacher.Id, 10) + "',")
		sql.WriteString("'" + strconv.FormatInt(teacher.EventId, 10) + "',")
		sql.WriteString("'" + strconv.FormatInt(teacher.CsId, 10) + "',")
		sql.WriteString("'" + strconv.FormatInt(teacher.TeacherId, 10) + "',")
		sql.WriteString("'" + teacher.TeacherName + "',")
		sql.WriteString("'" + teacher.TeacherAccount + "',")
		sql.WriteString("'" + strconv.Itoa(int(teacher.Enable)) + "',")
		sql.WriteString("'" + teacher.Ctime.Format(comm.TimeFormatTime) + "',")
		if teacher.Utime.IsZero() {
			sql.WriteString("NULL")
		} else {
			sql.WriteString("'" + teacher.Utime.Format(comm.TimeFormatTime) + "'")
		}
		sql.WriteString(")")
		if index < lang-1 {
			sql.WriteString(",")
		}
	}
	return sql.String()
}

func (c configurationService) stuSelectConfigSubject(eventId string) *[]model.ConfigurationSimpleDo {
	var configSimpleDO []model.ConfigurationSimpleDo
	sql := "SELECT " +
		"s.id, " +
		"s.subject_id, " +
		"s.subject_name, " +
		"s.class_name, " +
		"s.introduction, " +
		"s.num, " +
		"s.teach_address, " +
		"s.teach_time, " +
		"tt.teacher " +
		"FROM configuration_subject s " +
		"LEFT JOIN " +
		"(SELECT " +
		"GROUP_CONCAT( t.teacher_name ) AS teacher , " +
		"t.cs_id as cs_id " +
		"FROM " +
		"configuration_teacher t " +
		"WHERE t.`enable` = ?  " +
		"AND t.event_id = ? " +
		"GROUP BY t.cs_id ) tt " +
		"ON " +
		"tt.cs_id = s.id " +
		"WHERE " +
		"s.event_id = ? AND " +
		"s.`enable` = ? " +
		"ORDER BY s.ctime DESC"
	if err := c.mysql.SQL(sql, comm.Enable, eventId, eventId, comm.Enable).Find(&configSimpleDO); err != nil {
		zap.S().Error(err)
	}
	return &configSimpleDO
}

func (c configurationService) getIntroduction(subjectId string) string {

	selectSql := "SELECT introduction FROM subject where id = ? AND enable = ?"
	var introduction string
	if _, err := c.mysql.SQL(selectSql, subjectId, comm.Enable).Get(&introduction); err != nil {
		zap.S().Error(err)
	}
	return introduction
}

func (c configurationService) buildDataUpdateCache(configSubject *model.ConfigurationSubject, configTeachers *[]model.ConfigurationTeacher) {
	do := c.GetCacheCSubjectByEventAndCsId(strconv.FormatInt(configSubject.EventId, 10), strconv.FormatInt(configSubject.Id, 10))

	if do != nil {
		do.SubjectId = configSubject.SubjectId
		do.SubjectName = configSubject.SubjectName
		do.ClassName = configSubject.ClassName
		do.Introduction = configSubject.Introduction
		do.Num = int(configSubject.Num)
		do.TeachAddress = configSubject.TeachAddress
		do.TeachTime = configSubject.TeachTime

		var teacherName []string
		for _, teacher := range *configTeachers {
			teacherName = append(teacherName, teacher.TeacherName)
		}
		do.Teacher = strings.Join(teacherName, ",")
		var newSimpleDo []model.ConfigurationSimpleDo
		newSimpleDo = append(newSimpleDo, *do)
		c.updateCacheConfigSimpleDo(strconv.FormatInt(configSubject.EventId, 10), &newSimpleDo)
	}
}

// 删除缓存
func (c configurationService) deleteCacheConfigSimpleDo(eventId, csId string) {
	result, err := db.RedisClient.HLen(comm.RedisStuSelectConfigSubjectListKey + eventId).Result()
	if err != nil {
		zap.S().Error("不需要更新缓存,查询redis错误，eventId : v%", eventId, err)
		return
	}
	if result == 0 {
		zap.S().Error("不需要更新缓存，eventId : v%", eventId)
		return
	}

	if _, err := db.RedisClient.HDel(comm.RedisStuSelectConfigSubjectListKey+eventId, csId).Result(); err != nil {
		zap.S().Error("删除缓存失败,eventId: "+eventId+"csId: "+csId, err)
	}
}

// 更新缓存
func (c configurationService) updateCacheConfigSimpleDo(eventId string, configSimpleDo *[]model.ConfigurationSimpleDo) {
	result, err := db.RedisClient.HLen(comm.RedisStuSelectConfigSubjectListKey + eventId).Result()
	if err != nil {
		zap.S().Error("不需要更新缓存,查询redis错误，eventId: "+eventId, err.Error())
		return
	}
	if result == 0 {
		zap.S().Error("不需要更新缓存,eventId: " + eventId)
		return
	}

	for _, do := range *configSimpleDo {
		marshal, _ := json.Marshal(do)
		if _, err := db.RedisClient.HSet(comm.RedisStuSelectConfigSubjectListKey+eventId, strconv.FormatInt(do.Id, 10), marshal).Result(); err != nil {
			zap.S().Error("更新缓存失败,eventId: "+eventId+"arg: "+string(marshal), err.Error())
		}
	}
}

// 设置缓存，学生查询选课课程列表的缓存
func (c configurationService) setCacheConfigSimpleDo(eventId string, configSimpleDo *[]model.
	ConfigurationSimpleDo) {

	m := make(map[string]interface{})
	for _, do := range *configSimpleDo {
		marshal, _ := json.Marshal(do)
		m[strconv.FormatInt(do.Id, 10)] = marshal
	}
	if _, err := db.RedisClient.HMSet(comm.RedisStuSelectConfigSubjectListKey+eventId, m).Result(); err != nil {
		zap.S().Error(err)
	}
}

// 获取缓存，学生查询选课课程列表的缓存
func (c configurationService) getCacheStuSelectConfigSubject(eventId string) *[]model.ConfigurationSimpleDo {
	if result, err := db.RedisClient.HGetAll(comm.RedisStuSelectConfigSubjectListKey + eventId).Result(); err == err {
		var configSimpleDos []model.ConfigurationSimpleDo
		for _, element := range result {
			var m model.ConfigurationSimpleDo
			if err = json.Unmarshal([]byte(element), &m); err != nil {
				zap.S().Error("Json Unmarshal Error ", err)
				return nil
			}
			configSimpleDos = append(configSimpleDos, m)
		}
		return &configSimpleDos
	}
	return nil
}

// 根据选课ID 和 CSID 获取缓存
func (c configurationService) GetCacheCSubjectByEventAndCsId(eventId, csId string) *model.ConfigurationSimpleDo {
	result, err := db.RedisClient.HGet(comm.RedisStuSelectConfigSubjectListKey+eventId, csId).Result()
	if err != nil {
		zap.S().Error("查询redis错误，eventId: "+eventId, err)
		return nil
	}
	var do *model.ConfigurationSimpleDo
	if err = json.Unmarshal([]byte(result), &do); err != nil {
		zap.S().Error("Json Unmarshal Error ", err)
		return nil
	}
	return do
}

// 获取缓存，学生选课的课程剩余名额情况
func (c configurationService) getCacheStuSelectConfigSubjectRemainingPlaces(eventId string) map[string]string {
	//{
	//	test := make(map[string]interface{})
	//	test["1261600976368832512"] = "10"
	//	test["1261601135634944000"] = "9"
	//	test["1261601254707040256"] = "8"
	//	test["1261601416653312000"] = "7"
	//	test["1261601723697336320"] = "6"
	//	_, err := db.RedisClient.HMSet(comm.RedisStuSelectConfigSubjectRemainingPlacesKey+eventId, test).Result()
	//	if err != nil {
	//		iris.New().Logger().Error("redis set hash error ", err)
	//	}
	//}
	result, err := db.RedisClient.HGetAll(comm.RedisStuSelectConfigSubjectRemainingPlacesKey + eventId).Result()
	if err == nil {
		return result
	}
	return nil
}

type ConfigurationService interface {

	// 根据选课事件ID查询所有的选课课程
	PostList(eventId string, sessionUser *comm.SessionUSER) (*[]domain.ConfigurationVo, error)

	// 创建
	PostCreate(arg *domain.ConfigurationCreateArg) error

	// 更新
	PostUpdate(arg *domain.ConfigurationUpdateArg) error

	// 删除选课课程,级联删除该课程下的教师
	PostRemove(eventId, configSubjectId string, sessionUser *comm.SessionUSER) error

	// 删除选课教师配置
	PostRemoveTeacher(configTeacherId string, sessionUser *comm.SessionUSER) error

	// 根据选课事件ID查询选课课程 学生
	PostListStu(eventId string) *[]domain.ConfigurationSimpleVo

	// 批量创建
	PostCreateBatch(arg domain.ConfigurationCreateBatchArg, user *comm.SessionUSER) error

	// 根据选课ID 和 CSID 获取缓存
	GetCacheCSubjectByEventAndCsId(event, csId string) *model.ConfigurationSimpleDo

	// 学生查询自己选课记录
	PostStuHistoryStuSelected(id string, user *comm.SessionUSER) (*[]domain.StudentSubjectSelectVo, error)

	// 查询选课下的课程详情
	PostCsDetails(eventId string, user *comm.SessionUSER) ([]domain.ConfigurationDetailsVo, error)

	//查询选课下的课程详情的学生列表
	PostPageCsDetailsStu(arg domain.StuSubjectPageArg, user *comm.SessionUSER) (interface{}, error)

	// 教师追加学生到课堂 搜索学生
	PostAppendStuSearch(arg domain.ConfigurationEventIdAndDestArg, user *comm.SessionUSER) (*[]domain.StudentSubjectSelectInfoVo, error)
}

func NewConfigurationService(mysql *xorm.Engine, userService *UserService) ConfigurationService {
	return &configurationService{mysql: mysql, userService: userService}
}
