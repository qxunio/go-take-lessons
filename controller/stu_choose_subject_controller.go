package controller

import (
	"github.com/kataras/iris/v12"
	"go-take-lessons/domain"
	"go-take-lessons/domain/comm"
	"go-take-lessons/service"
	"go-take-lessons/tools"
	"go.uber.org/zap"
)

type StuChooseSubjectController struct {
	Ctx                     iris.Context
	StuChooseSubjectService service.StuChooseSubjectService
}

// 学生选课
func (sc *StuChooseSubjectController) PostLock() *comm.ResultMsg {
	var arg domain.ConfigurationEventIdAndCsIdArg

	if err := tools.GetRequestJson(&sc.Ctx, &arg); err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}

	err := sc.StuChooseSubjectService.PostLock(arg.EventId, arg.CsId, tools.GetSessionUser(&sc.Ctx))
	if err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	return comm.SuccessResponse()
}

// 学生删除课程
func (sc *StuChooseSubjectController) PostDelete() *comm.ResultMsg {
	var arg domain.ConfigurationEventIdAndCsIdArg

	if err := tools.GetRequestJson(&sc.Ctx, &arg); err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}

	err := sc.StuChooseSubjectService.PostDelete(arg, tools.GetSessionUser(&sc.Ctx))
	if err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	return comm.SuccessResponse()
}

// 管理员追加学生到课堂
func (sc *StuChooseSubjectController) PostAppendStu() *comm.ResultMsg {
	var arg domain.ConfigurationAdminAppendArg
	if err := tools.GetRequestJson(&sc.Ctx, &arg); err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	if err := sc.StuChooseSubjectService.PostAppendStu(arg, tools.GetSessionUser(&sc.Ctx)); err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	return comm.SuccessResponse()
}

// 管理员替换学生到课堂
func (sc *StuChooseSubjectController) PostReplaceStu() *comm.ResultMsg {
	var arg domain.ConfigurationAdminReplaceArg
	if err := tools.GetRequestJson(&sc.Ctx, &arg); err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	if err := sc.StuChooseSubjectService.PostReplaceStu(arg, tools.GetSessionUser(&sc.Ctx)); err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	return comm.SuccessResponse()
}

// 管理员删除学生选课
func (sc *StuChooseSubjectController) PostDeleteStu() *comm.ResultMsg {
	var arg domain.ConfigurationAdminAppendArg
	if err := tools.GetRequestJson(&sc.Ctx, &arg); err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	if err := sc.StuChooseSubjectService.PostDeleteStu(arg, tools.GetSessionUser(&sc.Ctx)); err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	return comm.SuccessResponse()
}

// 管理员导出报表
func (sc *StuChooseSubjectController) PostExport() {
	var arg domain.IdArg
	if err := tools.GetRequestJson(&sc.Ctx, &arg); err != nil {
		zap.S().Error("导出选课报表,: 获取参数错误", err.Error())
		_, _ = sc.Ctx.JSON("入参错误")
	}
	file, err := sc.StuChooseSubjectService.PostExport(arg, tools.GetSessionUser(&sc.Ctx))
	if err != nil {
		_, _ = sc.Ctx.JSON(comm.ErrorResponseMsg(err.Error()))
		return
	}
	if err = file.Write(sc.Ctx.ResponseWriter()); err != nil {
		zap.S().Error("用户CTL,导出用户: 写出响应错误", err.Error())
		_, _ = sc.Ctx.JSON("导出失败")
		return
	}

}
