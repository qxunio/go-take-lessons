package controller

import (
	"encoding/base64"
	"github.com/kataras/iris/v12"
	"go-take-lessons/db"
	"go-take-lessons/domain"
	"go-take-lessons/domain/comm"
	"go-take-lessons/domain/comm/ere"
	"go-take-lessons/service"
	"go-take-lessons/tools"
	"go.uber.org/zap"
)

// 认证 controller
type AuthController struct {
	Ctx         iris.Context
	AuthService service.AuthService
}

// 刷新TOKEN
func (au *AuthController) PostReload() *comm.ResultMsg {
	if token, err := au.AuthService.PostReload(au.Ctx.GetHeader(comm.AuthorizationHeader),
		tools.GetSessionUser(&au.Ctx)); err == nil {
		return comm.SuccessResponseData(token)
	} else {
		return comm.ErrorResponseMsg(err.Error())
	}
}

// Logout
func (au *AuthController) PostLogout() *comm.ResultMsg {
	if user := tools.GetSessionUser(&au.Ctx); user != nil {
		au.AuthService.PostLogout(user)
	}
	return comm.SuccessResponse()
}

// Login
func (au *AuthController) PostLogin() *comm.ResultMsg {
	var arg domain.AuthLoginarg
	if err := tools.GetRequestJson(&au.Ctx, &arg); err != nil {
		return comm.ErrorResponseMsg(err.Error())

	}
	UUID := arg.A3
	privateKey, err := db.RedisClient.Get(comm.RedisAuthEncryptionKey + UUID).Result()
	if err != nil {
		zap.S().Error("认证CTR,登录: 缓存获取私钥失败", err.Error())
		return comm.ErrorResponseMsg(ere.StrLoginFail)
	}

	bAccount, err := base64.StdEncoding.DecodeString(arg.A1)
	if err != nil {
		zap.S().Error("认证CTR,登录: base64账号解码失败", err.Error())
		return comm.ErrorResponseMsg(ere.StrLoginFail)
	}
	bPassword, err := base64.StdEncoding.DecodeString(arg.A2)
	if err != nil {
		zap.S().Error("认证CTR,登录: base64密码解码失败", err.Error())
		return comm.ErrorResponseMsg(ere.StrLoginFail)
	}

	account := tools.RSADecrypt(bAccount, []byte(privateKey))
	password := tools.RSADecrypt(bPassword, []byte(privateKey))

	vo, err := au.AuthService.PostLogin(string(account), string(password))
	if err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	return comm.SuccessResponseData(vo)
}
