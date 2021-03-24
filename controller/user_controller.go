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
	"strings"
)

type UserController struct {
	Ctx         iris.Context
	UserService service.UserService
}

// 创建
func (u *UserController) PostCreate() *comm.ResultMsg {
	var arg domain.UserCreateArg
	if err := tools.GetRequestJson(&u.Ctx, &arg); err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	if tools.StringToInt(arg.Type) == comm.Student {
		if tools.IsBlank(arg.SchoolYear) {
			return comm.ErrorResponseMsg("学生必须选择入学年")
		}
	}
	UUID := arg.Uid
	privateKey, err := db.RedisClient.Get(comm.RedisAuthEncryptionKey + UUID).Result()

	if err != nil {
		zap.S().Error("用户CTR,创建: 缓存获取私钥失败", err.Error())
		return comm.ErrorResponseMsg("超时，请重试")
	}
	byteAccount, err := base64.StdEncoding.DecodeString(arg.Account)
	if err != nil {
		zap.S().Error("base64 StdEncoding DecodeString fail Account ", err.Error())
		return comm.ErrorResponseMsg(ere.StrParseArgFail)
	}
	bytePassword, err := base64.StdEncoding.DecodeString(arg.Password)
	if err != nil {
		zap.S().Error("用户CTR,创建: base64密码解码失败", err.Error())
		return comm.ErrorResponseMsg(ere.StrParseArgFail)
	}
	account := tools.RSADecrypt(byteAccount, []byte(privateKey))
	password := tools.RSADecrypt(bytePassword, []byte(privateKey))
	arg.Account = string(account)
	arg.Password = string(password)
	vo, err := u.UserService.PostCreate(&arg, tools.GetSessionUser(&u.Ctx))
	if err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	return comm.SuccessResponseData(vo)
}

// 修改
func (u *UserController) PostUpdate() *comm.ResultMsg {
	var arg domain.UserUpdateArg
	if err := tools.GetRequestJson(&u.Ctx, &arg); err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	if tools.StringToInt(arg.Type) == comm.Student {
		if tools.IsBlank(arg.SchoolYear) {
			return comm.ErrorResponseMsg("学生必须选择入学年")
		}
	}
	if err := u.UserService.PostUpdate(&arg, tools.GetSessionUser(&u.Ctx)); err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	return comm.SuccessResponse()
}

// 删除
func (u *UserController) PostRemove() *comm.ResultMsg {
	userId := u.Ctx.URLParam("userId")
	if tools.IsBlank(userId) {
		return comm.ErrorResponseMsg(ere.StrCommArgIsBlank)
	}
	if err := u.UserService.PostRemove(userId, tools.GetSessionUser(&u.Ctx)); err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	return comm.SuccessResponse()
}

// 分页
func (u *UserController) PostPage() *comm.ResultMsg {
	var arg domain.UserPageArg
	if err := tools.GetRequestJson(&u.Ctx, &arg); err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}

	vo, err := u.UserService.PostPage(&arg, tools.GetSessionUser(&u.Ctx))
	if err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	return comm.SuccessResponseData(vo)
}

// 导出用户
func (u *UserController) PostExport() {
	file, err := u.UserService.PostExport(tools.GetSessionUser(&u.Ctx))
	if err != nil {
		_, _ = u.Ctx.JSON(comm.ErrorResponseMsg(err.Error()))
		return
	}
	if err = file.Write(u.Ctx.ResponseWriter()); err != nil {
		zap.S().Error("用户CTL,导出用户: 写出响应错误", err.Error())
		_, _ = u.Ctx.JSON("导出失败")
		return
	}
}

// 导入用户
func (u *UserController) PostImport() *comm.ResultMsg {
	fileReader, fileInfo, err := u.Ctx.FormFile("file")
	if err != nil {
		return comm.ErrorResponseMsg("上传失败")
	}

	userType := u.Ctx.FormValue("userType")
	if tools.IsBlank(userType) {
		return comm.ErrorResponseMsg("上传的用户类型未选择")
	}

	schoolYear := u.Ctx.FormValue("schoolYear")
	if tools.StringToInt(userType) == comm.Student {
		if tools.IsBlank(schoolYear) {
			return comm.ErrorResponseMsg("学生必须选择入学年")
		}
	}

	if fileInfo.Size > 3145728 {
		return comm.ErrorResponseMsg("上传文件大小不能超过3M")
	}
	if strings.HasSuffix(fileInfo.Filename, ".xlxs") {
		return comm.ErrorResponseMsg("上传的文件类型错误")
	}
	existenceUserVo, err := u.UserService.PostImport(fileReader, userType, tools.GetSessionUser(&u.Ctx), schoolYear)

	if err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	if len(existenceUserVo) != 0 {
		return comm.SuccessResponseData(existenceUserVo)
	}
	return comm.SuccessResponseMsg("导入成功")
}

// 重置密码
func (u *UserController) PostReset() *comm.ResultMsg {
	userId := u.Ctx.URLParam("userId")
	if tools.IsBlank(userId) {
		return comm.ErrorResponseMsg(ere.StrCommArgIsBlank)
	}
	if err := u.UserService.PostReset(userId, tools.GetSessionUser(&u.Ctx)); err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	return comm.SuccessResponse()
}

// 查询所有教师
func (u *UserController) PostListSimpleTeacher() *comm.ResultMsg {
	vo, err := u.UserService.PostListSimpleTeacher(tools.GetSessionUser(&u.Ctx))
	if err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	return comm.SuccessResponseData(vo)
}

// 导出模板
func (u *UserController) PostExportTemplate() {
	file := u.UserService.PostExportTemplate()
	if err := file.Write(u.Ctx.ResponseWriter()); err != nil {
		zap.S().Error("用户CTL,导出模板: 写出响应错误", err.Error())
		_, _ = u.Ctx.JSON("导出失败")
		return
	}
}
