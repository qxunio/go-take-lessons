package tools

import (
	"github.com/asaskevich/govalidator"
	"github.com/kataras/iris/v12"
	"go-take-lessons/domain/comm"
	"go.uber.org/zap"
)

// 获取请求JSON参数
func GetRequestJson(ctx *iris.Context, d interface{}) error {
	c := *ctx
	var err error
	if err = c.ReadJSON(d); err != nil {
		zap.S().Error(err)
		return err
	}
	if _, err = govalidator.ValidateStruct(d); err != nil {
		zap.S().Error(err)
		return err
	}
	return nil
}

// 获取请求User参数
func GetSessionUser(ctx *iris.Context) *comm.SessionUSER {
	context := *ctx
	if sessionUser, ok := context.Values().Get(comm.ContextSessionUserKey).(*comm.SessionUSER); ok {
		return sessionUser
	}
	return nil
}
