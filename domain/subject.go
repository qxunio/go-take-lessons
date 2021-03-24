package domain

import (
	"go-take-lessons/domain/comm"
	"go-take-lessons/model"
	"strconv"
)

type SubjectCreateArg struct {
	Name         string `json:"name" valid:"required,isNotBlank~学科名称不能为空"`
	Introduction string `json:"introduction"  valid:"-"`
}

type SubjectUpdateArg struct {
	Id           string `json:"id" valid:"required,isNotBlank~ID不能为空"`
	Name         string `json:"name" valid:"required,isNotBlank~学科名称不能为空"`
	Introduction string `json:"introduction"  valid:"-"`
}

type SubjectPageArg struct {
	comm.PageParam `valid:"required"`
	Name           string `json:"name" valid:"-"`
}

type SubjectVo struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Introduction string `json:"introduction"`
	Ctime        string `json:"ctime"`
}

type SubjectSimpleVo struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func CvSubjectToSimpleVo(subject *model.Subject) *SubjectSimpleVo {
	return &SubjectSimpleVo{
		Id:   strconv.FormatInt(subject.Id, 10),
		Name: subject.Name,
	}
}

func CvSubjectToVo(subject *model.Subject) *SubjectVo {
	return &SubjectVo{
		Id:           strconv.FormatInt(subject.Id, 10),
		Name:         subject.Name,
		Introduction: subject.Introduction,
		Ctime:        subject.Ctime.Format(comm.TimeFormatTime),
	}
}
