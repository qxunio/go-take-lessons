package domain

import (
	"go-take-lessons/domain/comm"
	"go-take-lessons/model"
	"strconv"
)

type StuSubjectPageArg struct {
	comm.PageParam `valid:"required"`
	Dest           string `json:"dest" valid:"-"`
	Class          string `json:"class" valid:"-"`
	CsId           string `json:"csId"  valid:"required,isNotBlank~课堂不能为空"`
}

type StuSubjectVo struct {
	Id         string `json:"id"`
	UserId     string `json:"userId"`
	Name       string `json:"name"`
	Account    string `json:"account"`
	Class      string `json:"class"`
	SchoolYear string `json:"schoolYear"`
	SelectTime string `json:"selectTime"`
}

func CvStuSubjectDoToStuSubjectVo(model model.StuSubjectDo) StuSubjectVo {
	return StuSubjectVo{
		Id:         strconv.FormatInt(model.Id, 10),
		UserId:     strconv.FormatInt(model.UserId, 10),
		Name:       model.Name,
		Account:    model.Account,
		Class:      strconv.FormatInt(int64(model.Class), 10),
		SchoolYear: strconv.FormatInt(int64(model.SchoolYear), 10),
		SelectTime: model.Ctime.Format(comm.TimeFormatTime),
	}
}
