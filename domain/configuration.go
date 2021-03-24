package domain

import (
	"go-take-lessons/domain/comm"
	"go-take-lessons/model"
	"go-take-lessons/tools"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
)

type ConfigurationSubjectVo struct {
	Id             string `json:"id"`
	EventId        string `json:"eventId"`
	SubjectId      string `json:"subjectId"`
	SubjectName    string `json:"subjectName"`
	ClassName      string `json:"className"`
	Num            string `json:"num"`
	SelectedPlaces string `json:"selectedPlaces"`
	TeachAddress   string `json:"teachAddress"`
	TeachTime      string `json:"teachTime"`
}

type ConfigurationTeacherVo struct {
	Id             string `json:"id"`
	EventId        string `json:"eventId"`
	CsId           string `json:"csId"`
	TeacherId      string `json:"teacherId"`
	TeacherName    string `json:"teacherName"`
	TeacherAccount string `json:"teacherAccount"`
	Ctime          string `json:"ctime"`
}

type ConfigurationVo struct {
	ConfigurationSubjectVo
	Teachers []ConfigurationTeacherVo `json:"teachers"`
}

type ConfigurationSimpleVo struct {
	Id              string `json:"id"`
	SubjectId       string `json:"subjectId"`
	SubjectName     string `json:"subjectName"`
	ClassName       string `json:"className"`
	Introduction    string `json:"introduction"`
	TeachAddress    string `json:"teachAddress"`
	TeachTime       string `json:"teachTime"`
	Teacher         string `json:"teacher"`
	RemainingPlaces string `json:"remainingPlaces"`
	Num             string `json:"num"`
}

type ConfigurationDetailsVo struct {
	ConfigurationSimpleVo
	SelectedPlaces string `json:"selectedPlaces"`
}

type ConfigurationCreateTeacherArg struct {
	TeacherId      string `json:"teacherId" valid:"required,isNotBlank~教师不能为空"`
	TeacherName    string `json:"teacherName" valid:"required,isNotBlank~教师姓名不能为空"`
	TeacherAccount string `json:"teacherAccount" valid:"required,isNotBlank~教师账号不能为空"`
}

type ConfigurationCreateSubjectArg struct {
	EventId      string `json:"eventId" valid:"required,isNotBlank~选课不能为空"`
	SubjectId    string `json:"subjectId" valid:"required,isNotBlank~学科不能为空"`
	Num          string `json:"num" valid:"required,isNotBlank~限制人数不能为空"`
	SubjectName  string `json:"subjectName" valid:"-"`
	ClassName    string `json:"className" valid:"required,isNotBlank~课堂名称不能为空"`
	TeachAddress string `json:"teachAddress" valid:"-"`
	TeachTime    string `json:"teachTime"  valid:"-"`
}

type ConfigurationCreateArg struct {
	ConfigurationCreateSubjectArg `valid:"required"`
	Teachers                      []ConfigurationCreateTeacherArg `json:"teachers" valid:"-"`
}

type ConfigurationCreateBatchSubjectArg struct {
	SubjectId   string   `json:"id" valid:"required,isNotBlank~学科不能为空"`
	SubjectName string   `json:"subjectName" valid:"required,isNotBlank~学科名称不能为空"`
	ClassName   string   `json:"className" valid:"required,isNotBlank~课堂名称不能为空"`
	Limit       string   `json:"limit" valid:"required,isNotBlank~限制人数不能为空"`
	Teacher     []string `json:"teacher" valid:"-"`
}

type ConfigurationCreateBatchArg struct {
	Subject []ConfigurationCreateBatchSubjectArg `json:"subject" valid:"required"`
	Address string                               `json:"address"`
	Time    string                               `json:"time"`
	EventId string                               `json:"eventId" valid:"required,isNotBlank~选课不能为空"`
}

type ConfigurationUpdateArg struct {
	ConfigurationUpdateSubjectArg `valid:"required"`
	Teachers                      []ConfigurationUpdateTeacherArg `json:"teachers" valid:"required"`
}

type ConfigurationUpdateTeacherArg struct {
	Id                            string `json:"id" valid:"-"`
	ConfigurationCreateTeacherArg `valid:"required"`
	Ctime                         string `json:"ctime"  valid:"-"`
}

type ConfigurationUpdateSubjectArg struct {
	Id                            string `json:"id" valid:"required,isNotBlank~id不能为空"`
	ConfigurationCreateSubjectArg `valid:"required"`
}

type ConfigurationEventIdAndCsIdArg struct {
	EventId string `json:"eventId" valid:"required,isNotBlank~id不能为空"`
	CsId    string `json:"csId" valid:"required,isNotBlank~Cs不能为空"`
}

type ConfigurationEventIdAndDestArg struct {
	EventId string `json:"eventId" valid:"required,isNotBlank~选课id不能为空"`
	Dest    string `json:"dest" valid:"required,isNotBlank~用户名或账号不能为空"`
}

type ConfigurationAdminAppendArg struct {
	EventId string `json:"eventId" valid:"required,isNotBlank~选课id不能为空"`
	CsId    string `json:"csId" valid:"required,isNotBlank~Cs不能为空"`
	Uid     string `json:"uid" valid:"required,isNotBlank~学生不能为空"`
}

type ConfigurationAdminReplaceArg struct {
	ConfigurationAdminAppendArg `valid:"required"`
	ReplaceCsId                 string `json:"replaceId" valid:"required,isNotBlank~需要替换的课程ID不能为空"`
}

type StudentSubjectSelectVo struct {
	Id           string `json:"id"`
	Class        string `json:"class"`
	SchoolYear   string `json:"schoolYear"`
	SelectTime   string `json:"selectTime"`
	ClassName    string `json:"className"`
	SubjectName  string `json:"subjectName"`
	TeachAddress string `json:"teachAddress"`
	TeachTime    string `json:"teachTime"`
	Introduction string `json:"introduction"`
	TeacherName  string `json:"teacherName"`
}

type StudentSubjectSelectInfoVo struct {
	UserId      string `json:"userId"`
	Name        string `json:"name"`
	Account     string `json:"account"`
	Class       string `json:"class"`
	SchoolYear  string `json:"schoolYear"`
	CsId        string `json:"csId"`
	SelectTime  string `json:"selectTime"`
	ClassName   string `json:"className"`
	SubjectName string `json:"subjectName"`
}

func CvConfigurationCreateArgToModel(arg *ConfigurationCreateArg, introduction string) (*model.ConfigurationSubject, *[]model.ConfigurationTeacher) {
	var configSubject model.ConfigurationSubject
	nowTime := time.Now()

	configSubject.Id = tools.SnowFlake.Generate().Int64()
	configSubject.EventId = tools.StringToInt64(arg.EventId)
	configSubject.SubjectId = tools.StringToInt64(arg.SubjectId)
	configSubject.SubjectName = strings.TrimSpace(arg.SubjectName)
	configSubject.Introduction = introduction
	if tools.IsBlank(arg.ClassName) {
		configSubject.ClassName = configSubject.SubjectName
	} else {
		configSubject.ClassName = strings.TrimSpace(arg.ClassName)
	}
	configSubject.Num = tools.StringToUint8(arg.Num)
	configSubject.TeachAddress = strings.TrimSpace(arg.TeachAddress)
	configSubject.TeachTime = strings.TrimSpace(arg.TeachTime)
	configSubject.Ctime = nowTime
	configSubject.Enable = comm.Enable

	configTeachers := make([]model.ConfigurationTeacher, 0)
	for _, configTeacherArg := range arg.Teachers {
		var configTeacher model.ConfigurationTeacher
		configTeacher.Id = tools.SnowFlake.Generate().Int64()
		configTeacher.EventId = tools.StringToInt64(arg.EventId)
		configTeacher.CsId = configSubject.Id
		configTeacher.TeacherId = tools.StringToInt64(configTeacherArg.TeacherId)
		configTeacher.TeacherName = strings.TrimSpace(configTeacherArg.TeacherName)
		configTeacher.TeacherAccount = strings.TrimSpace(configTeacherArg.TeacherAccount)
		configTeacher.Enable = comm.Enable
		configTeacher.Ctime = nowTime
		configTeachers = append(configTeachers, configTeacher)
	}
	return &configSubject, &configTeachers
}

func CvConfigurationUpdateArgToModel(arg *ConfigurationUpdateArg, introduction string) (*model.ConfigurationSubject, *[]model.ConfigurationTeacher) {
	var configSubject model.ConfigurationSubject
	nowTime := time.Now()

	configSubject.Id = tools.StringToInt64(arg.Id)
	configSubject.SubjectId = tools.StringToInt64(arg.SubjectId)
	configSubject.SubjectName = strings.TrimSpace(arg.SubjectName)
	configSubject.Introduction = introduction
	if tools.IsBlank(arg.ClassName) {
		configSubject.ClassName = configSubject.SubjectName
	} else {
		configSubject.ClassName = strings.TrimSpace(arg.ClassName)
	}
	configSubject.Num = tools.StringToUint8(arg.Num)
	configSubject.TeachAddress = strings.TrimSpace(arg.TeachAddress)
	configSubject.TeachTime = strings.TrimSpace(arg.TeachTime)
	configSubject.Utime = nowTime
	configSubject.EventId = tools.StringToInt64(arg.EventId)

	configTeachers := make([]model.ConfigurationTeacher, 0)
	for _, configTeacherArg := range arg.Teachers {
		var configTeacher model.ConfigurationTeacher
		if tools.IsBlank(configTeacherArg.Id) {
			configTeacher.Id = tools.SnowFlake.Generate().Int64()
			configTeacher.Ctime = nowTime
		} else {
			configTeacher.Id = tools.StringToInt64(configTeacherArg.Id)
			ctime, err := time.ParseInLocation(comm.TimeFormatTime, configTeacherArg.Ctime, time.Local)
			if err != nil {
				zap.S().Error(err)
			}
			configTeacher.Ctime = ctime
			configTeacher.Utime = nowTime
		}
		configTeacher.EventId = configSubject.EventId
		configTeacher.CsId = configSubject.Id
		configTeacher.Enable = comm.Enable
		configTeacher.TeacherId = tools.StringToInt64(configTeacherArg.TeacherId)
		configTeacher.TeacherName = strings.TrimSpace(configTeacherArg.TeacherName)
		configTeacher.TeacherAccount = strings.TrimSpace(configTeacherArg.TeacherAccount)
		configTeachers = append(configTeachers, configTeacher)
	}
	return &configSubject, &configTeachers
}

func CvConfigSimpleDOToVo(args *[]model.ConfigurationSimpleDo,
	remainingPlaces map[string]string) *[]ConfigurationSimpleVo {
	var configurationSimpleVos []ConfigurationSimpleVo

	if remainingPlaces == nil {
		for _, arg := range *args {
			var configSimpleVo ConfigurationSimpleVo
			configSimpleVo.Id = strconv.FormatInt(arg.Id, 10)
			configSimpleVo.SubjectId = strconv.FormatInt(arg.SubjectId, 10)
			configSimpleVo.SubjectName = arg.SubjectName
			configSimpleVo.ClassName = arg.ClassName
			configSimpleVo.Introduction = arg.Introduction
			configSimpleVo.TeachTime = arg.TeachTime
			configSimpleVo.TeachAddress = arg.TeachAddress
			configSimpleVo.Teacher = arg.Teacher
			configSimpleVo.RemainingPlaces = strconv.FormatInt(int64(arg.Num), 10)
			configSimpleVo.Num = strconv.FormatInt(int64(arg.Num), 10)
			configurationSimpleVos = append(configurationSimpleVos, configSimpleVo)
		}
	} else {
		for _, arg := range *args {
			var configSimpleVo ConfigurationSimpleVo
			configSimpleVo.Id = strconv.FormatInt(arg.Id, 10)
			configSimpleVo.SubjectId = strconv.FormatInt(arg.SubjectId, 10)
			configSimpleVo.SubjectName = arg.SubjectName
			configSimpleVo.ClassName = arg.ClassName
			configSimpleVo.Introduction = arg.Introduction
			configSimpleVo.TeachTime = arg.TeachTime
			configSimpleVo.TeachAddress = arg.TeachAddress
			configSimpleVo.Teacher = arg.Teacher
			if tools.IsBlank(remainingPlaces[configSimpleVo.Id]) {
				configSimpleVo.RemainingPlaces = strconv.FormatInt(int64(arg.Num), 10)
			} else {
				configSimpleVo.RemainingPlaces = strconv.FormatInt(int64(arg.Num)-tools.StringToInt64(remainingPlaces[configSimpleVo.Id]), 10)
			}
			configSimpleVo.Num = strconv.FormatInt(int64(arg.Num), 10)
			configurationSimpleVos = append(configurationSimpleVos, configSimpleVo)
		}
	}
	return &configurationSimpleVos
}

func CvConfigurationToConfigurationDetailsVo(arg model.ConfigurationSubject, teacher string) ConfigurationDetailsVo {
	var data ConfigurationDetailsVo
	data.Id = strconv.FormatInt(arg.Id, 10)
	data.Teacher = teacher
	data.TeachAddress = arg.TeachAddress
	data.TeachTime = arg.TeachTime
	data.Introduction = arg.Introduction
	data.ClassName = arg.ClassName
	data.Num = strconv.FormatUint(uint64(arg.Num), 10)
	data.SubjectId = strconv.FormatInt(arg.SubjectId, 10)
	data.SubjectName = arg.SubjectName
	rp := 0
	u := arg.Num - arg.SelectedPlaces
	if u > 0 {
		rp = int(u)
	}
	data.RemainingPlaces = strconv.FormatUint(uint64(rp), 10)
	data.SelectedPlaces = strconv.FormatUint(uint64(arg.SelectedPlaces), 10)
	return data
}
