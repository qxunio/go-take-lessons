package controller

import (
	"github.com/asaskevich/govalidator"
	"github.com/kataras/iris/v12"
	"go-take-lessons/domain"
	"go-take-lessons/domain/comm"
	"go-take-lessons/domain/comm/ere"
	"go-take-lessons/service"
	"go-take-lessons/tools"
	"go.uber.org/zap"
	"time"
)

type EventController struct {
	Ctx          iris.Context
	EventService service.EventService
}

// 学生查询历史
func (e *EventController) PostStuHistoryEvent() *comm.ResultMsg {
	vo, err := e.EventService.PostStuHistoryEvent(tools.GetSessionUser(&e.Ctx))
	if err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	return comm.SuccessResponseData(vo)
}

// 创建
func (e *EventController) PostCreate() *comm.ResultMsg {
	var arg domain.EventCreateArg
	if err := tools.GetRequestJson(&e.Ctx, &arg); err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	vo, err := e.EventService.PostCreate(&arg, tools.GetSessionUser(&e.Ctx))
	if err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	return comm.SuccessResponseData(vo)
}

// 修改
func (e *EventController) PostUpdate() *comm.ResultMsg {
	var arg domain.EventUpdateArg
	if err := tools.GetRequestJson(&e.Ctx, &arg); err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	if err := e.EventService.PostUpdate(&arg, tools.GetSessionUser(&e.Ctx)); err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	return comm.SuccessResponse()
}

// 分页
func (e *EventController) PostPage() *comm.ResultMsg {
	var arg comm.PageParam
	if err := tools.GetRequestJson(&e.Ctx, &arg); err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	vo, err := e.EventService.PostPage(&arg, tools.GetSessionUser(&e.Ctx))
	if err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	return comm.SuccessResponseData(vo)
}

// 删除
func (e *EventController) PostRemove() *comm.ResultMsg {
	eventId := e.Ctx.URLParam("eventId")
	if govalidator.IsNull(eventId) {
		return comm.ErrorResponseMsg(ere.StrCommArgIsBlank)
	}
	if err := e.EventService.PostRemove(eventId, tools.GetSessionUser(&e.Ctx)); err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	return comm.SuccessResponse()
}

// 修改选课状态
func (e *EventController) PostChange() *comm.ResultMsg {
	eventId := e.Ctx.URLParam("eventId")
	if tools.IsBlank(eventId) {
		return comm.ErrorResponseMsg(ere.StrCommArgIsBlank)
	}
	if err := e.EventService.PostChange(eventId, tools.GetSessionUser(&e.Ctx)); err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	return comm.SuccessResponse()
}

// 查询存在的有效选课 学生
func (e *EventController) PostExistEvent() *comm.ResultMsg {
	vo, err := e.EventService.PostExistEvent(tools.GetSessionUser(&e.Ctx))
	if err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	return comm.SuccessResponseData(vo)
}

// 重新激活选课
func (e *EventController) PostReactivation() *comm.ResultMsg {
	var arg domain.EventReactivationArg
	if err := tools.GetRequestJson(&e.Ctx, &arg); err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	etime, err := time.ParseInLocation(comm.TimeFormatTime, arg.Etime, time.Local)
	if err != nil {
		zap.S().Error(err)
		return comm.ErrorResponseMsg("格式化结束时间失败！")
	}
	if etime.Before(time.Now()) {
		return comm.ErrorResponseMsg("结束时间不能小于当前时间")
	}
	if err := e.EventService.PostReactivation(arg.Id, etime); err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	return comm.SuccessResponse()
}

// 修改学生是否可以修改选课
func (e *EventController) PostUpdateCanUpdate() *comm.ResultMsg {
	var arg domain.EventUpdateCanUpdateArg
	if err := tools.GetRequestJson(&e.Ctx, &arg); err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	if err := e.EventService.PostUpdateCanUpdate(arg, tools.GetSessionUser(&e.Ctx)); err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	return comm.SuccessResponse()
}

func (e *EventController) GetRedisGet() *comm.ResultMsg {
	eventId := e.Ctx.URLParam("key")
	return comm.SuccessResponseData(e.EventService.GetRedisGet(eventId))
}

func (e *EventController) GetRedisSet() *comm.ResultMsg {
	key := e.Ctx.URLParam("key")
	value := e.Ctx.URLParam("value")
	return comm.SuccessResponseData(e.EventService.GetRedisSet(key, value))
}

func (e *EventController) GetRedisDel() *comm.ResultMsg {
	key := e.Ctx.URLParam("key")
	return comm.SuccessResponseData(e.EventService.GetRedisDel(key))
}
