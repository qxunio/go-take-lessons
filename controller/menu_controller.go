package controller

import (
	"github.com/kataras/iris/v12"
	"go-take-lessons/domain/comm"
	"go-take-lessons/service"
	"go-take-lessons/tools"
)

type MenuController struct {
	Ctx         iris.Context
	MenuService service.MenuService
}

func (m *MenuController) PostRouter() *comm.ResultMsg {
	err, data := m.MenuService.PostRouter(tools.GetSessionUser(&m.Ctx))
	if err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	return comm.SuccessResponseData(data)
}
