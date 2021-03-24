package domain

import (
	"go-take-lessons/domain/comm"
	"go-take-lessons/model"
	"go-take-lessons/tools"
	"strconv"
	"strings"
)

type EventCreateArg struct {
	Name       string   `json:"name" valid:"required,isNotBlank~选课名称不能为空"`
	SchoolYear []string `json:"schoolYear"  valid:"required"`
	TagIds     []string `json:"tagIds"  `
	CanUpdate  string   `json:"canUpdate" valid:"required,isNotBlank~修改选课不能为空"`
	Num        string   `json:"num" valid:"required,isNotBlank~限制人数不能为空"`
	Stime      string   `json:"stime"  valid:"required,isNotBlank~开始时间不能为空"`
	Etime      string   `json:"etime"  valid:"required,isNotBlank~结束时间不能为空"`
}

type EventUpdateArg struct {
	EventCreateArg `valid:"required"`
	Id             string `json:"id" valid:"required,isNotBlank~ID不能为空"`
}

type EventReactivationArg struct {
	Id    string `json:"id" valid:"required,isNotBlank~ID不能为空"`
	Etime string `json:"etime" valid:"required,isNotBlank~结束时间不能为空"`
}

type EventVo struct {
	Id         string   `json:"id"`
	Name       string   `json:"name"`
	Num        string   `json:"num"`
	Stime      string   `json:"stime"`
	Etime      string   `json:"etime"`
	Status     string   `json:"status"`
	Ctime      string   `json:"ctime"`
	SchoolYear []string `json:"schoolYear"`
	CanUpdate  string   `json:"canUpdate"`
	TagIds     []string `json:"tagIds"`
}

type EventSimpleVo struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Num       string `json:"num"`
	Stime     string `json:"stime"`
	Etime     string `json:"etime"`
	Status    string `json:"status"`
	CanUpdate string `json:"canUpdate"`
}

type EventUpdateCanUpdateArg struct {
	Id     string `json:"id" valid:"required,isNotBlank~ID不能为空"`
	Status string `json:"status" valid:"required,isNotBlank~状态不能为空"`
}

func CvEventSimpleToVo(event model.Event) EventSimpleVo {
	return EventSimpleVo{
		Id:        strconv.FormatInt(event.Id, 10),
		Name:      event.Name,
		Num:       strconv.FormatUint(uint64(event.Num), 10),
		Stime:     event.Stime.Format(comm.TimeFormatTime),
		Etime:     event.Etime.Format(comm.TimeFormatTime),
		Status:    strconv.FormatUint(uint64(event.Status), 10),
		CanUpdate: strconv.FormatInt(int64(event.CanUpdate), 10),
	}
}

func CvEventToVo(event *model.Event) *EventVo {
	e := &EventVo{
		Id:         strconv.FormatInt(event.Id, 10),
		Name:       event.Name,
		Num:        strconv.FormatUint(uint64(event.Num), 10),
		Stime:      event.Stime.Format(comm.TimeFormatTime),
		Etime:      event.Etime.Format(comm.TimeFormatTime),
		Status:     strconv.FormatUint(uint64(event.Status), 10),
		Ctime:      event.Ctime.Format(comm.TimeFormatTime),
		SchoolYear: strings.Split(event.SchoolYear, ","),
		TagIds:     strings.Split(event.TagIds, ","),
		CanUpdate:  strconv.FormatInt(int64(event.CanUpdate), 10),
	}

	if tools.IsBlank(event.SchoolYear) {
		e.SchoolYear = nil
	} else {
		e.SchoolYear = strings.Split(event.SchoolYear, ",")
	}
	if tools.IsBlank(event.TagIds) {
		e.TagIds = nil
	} else {
		e.TagIds = strings.Split(event.TagIds, ",")
	}
	return e
}

func CvEventToStuVo(event *model.Event) *EventVo {
	return &EventVo{
		Id:        strconv.FormatInt(event.Id, 10),
		Name:      event.Name,
		Num:       strconv.FormatUint(uint64(event.Num), 10),
		Stime:     event.Stime.Format(comm.TimeFormatTime),
		Etime:     event.Etime.Format(comm.TimeFormatTime),
		Status:    strconv.FormatUint(uint64(event.Status), 10),
		Ctime:     event.Ctime.Format(comm.TimeFormatTime),
		CanUpdate: strconv.FormatInt(int64(event.CanUpdate), 10),
	}
}
