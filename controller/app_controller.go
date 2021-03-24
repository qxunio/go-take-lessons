package controller

import (
	"github.com/kataras/iris/v12"
	"go-take-lessons/domain/comm"
	"go-take-lessons/service"
)

// app controller
type AppController struct {
	Ctx        iris.Context
	AppService service.AppService
}

// Post 获取密钥
func (app *AppController) Post() *comm.ResultMsg {
	return comm.SuccessResponseData(app.AppService.Post())
}

// 管理员Index数据
func (app *AppController) PostAdminIndex() *comm.ResultMsg {
	return comm.SuccessResponseData(app.AppService.PostAdminIndex())
}
