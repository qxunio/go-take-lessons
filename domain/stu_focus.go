package domain

import (
	"go-take-lessons/domain/comm"
	"go-take-lessons/model"
	"go-take-lessons/tools"
	"time"
)

type StuFocusCreateArg struct {
	EventId         string `json:"eventId" valid:"required,isNotBlank~选课不能为空"`
	ConfigSubjectId string `json:"configSubjectId" valid:"required,isNotBlank~课程不能为空"`
}

type StuFocusRemoveArg struct {
	EventId         string `json:"eventId" valid:"required,isNotBlank~选课不能为空"`
	ConfigSubjectId string `json:"configSubjectId" valid:"required,isNotBlank~课程不能为空"`
}

func CvStuFocusCreateArgToModel(arg *StuFocusCreateArg, user *comm.SessionUSER) *model.StuFocus {
	var focus model.StuFocus
	focus.Id = tools.SnowFlake.Generate().Int64()
	focus.Enable = comm.Enable
	focus.EventId = tools.StringToInt64(arg.EventId)
	focus.CsId = tools.StringToInt64(arg.ConfigSubjectId)
	focus.UserId = user.Id
	focus.Ctime = time.Now()
	return &focus
}
