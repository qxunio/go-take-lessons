package controller

import (
	"github.com/kataras/iris/v12"
	"go-take-lessons/domain"
	"go-take-lessons/domain/comm"
	"go-take-lessons/domain/comm/ere"
	"go-take-lessons/service"
	"go-take-lessons/tools"
)

type StuFocusController struct {
	Ctx             iris.Context
	StuFocusService service.StuFocusService
}

// 创建
func (s *StuFocusController) PostCreate() *comm.ResultMsg {
	var arg domain.StuFocusCreateArg
	if err := tools.GetRequestJson(&s.Ctx, &arg); err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}

	err := s.StuFocusService.PostCreate(&arg, tools.GetSessionUser(&s.Ctx))
	if err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	return comm.SuccessResponse()
}

// 删除
func (s *StuFocusController) PostRemove() *comm.ResultMsg {
	var arg domain.StuFocusCreateArg
	if err := tools.GetRequestJson(&s.Ctx, &arg); err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}

	if err := s.StuFocusService.PostRemove(&arg, tools.GetSessionUser(&s.Ctx)); err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	return comm.SuccessResponse()
}

// 查询
func (s *StuFocusController) PostList() *comm.ResultMsg {
	eventId := s.Ctx.URLParam("eventId")
	if tools.IsBlank(eventId) {
		return comm.ErrorResponseMsg(ere.StrCommArgIsBlank)
	}
	return comm.SuccessResponseData(s.StuFocusService.PostList(eventId, tools.GetSessionUser(&s.Ctx)))
}
