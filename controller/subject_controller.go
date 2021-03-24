package controller

import (
	"github.com/kataras/iris/v12"
	"go-take-lessons/domain"
	"go-take-lessons/domain/comm"
	"go-take-lessons/domain/comm/ere"
	"go-take-lessons/service"
	"go-take-lessons/tools"
	"go.uber.org/zap"
)

type SubjectController struct {
	Ctx            iris.Context
	SubjectService service.SubjectService
}

// 创建
func (s *SubjectController) PostCreate() *comm.ResultMsg {
	var arg domain.SubjectCreateArg
	if err := tools.GetRequestJson(&s.Ctx, &arg); err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	vo, err := s.SubjectService.PostCreate(&arg, tools.GetSessionUser(&s.Ctx))
	if err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	return comm.SuccessResponseData(vo)
}

// 分页
func (s *SubjectController) PostPage() *comm.ResultMsg {
	var arg comm.PageParam
	if err := tools.GetRequestJson(&s.Ctx, &arg); err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	vo, err := s.SubjectService.PostPage(&arg, tools.GetSessionUser(&s.Ctx))
	if err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	return comm.SuccessResponseData(vo)
}

// 修改
func (s *SubjectController) PostUpdate() *comm.ResultMsg {
	var arg domain.SubjectUpdateArg
	if err := tools.GetRequestJson(&s.Ctx, &arg); err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	if err := s.SubjectService.PostUpdate(&arg); err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	return comm.SuccessResponse()

}

// 删除
func (s *SubjectController) PostRemove() *comm.ResultMsg {
	subjectId := s.Ctx.URLParam("subjectId")
	if tools.IsBlank(subjectId) {
		return comm.ErrorResponseMsg(ere.StrCommArgIsBlank)
	}
	if err := s.SubjectService.PostRemove(subjectId, tools.GetSessionUser(&s.Ctx)); err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	return comm.SuccessResponse()
}

// 导出学科
func (s *SubjectController) PostExport() {
	file, err := s.SubjectService.PostExport()
	if err != nil {
		_, _ = s.Ctx.JSON(comm.ErrorResponseMsg(err.Error()))
		return
	}
	if err = file.Write(s.Ctx.ResponseWriter()); err != nil {
		zap.S().Error("学科CTL,导出表格: 写出响应错误", err.Error())
		_, _ = s.Ctx.JSON(comm.ErrorResponseMsg("导出失败"))
		return
	}
}

// 查询所有学科
func (s *SubjectController) PostListSimpleSubject() *comm.ResultMsg {
	vo, err := s.SubjectService.PostListSimpleSubject(tools.GetSessionUser(&s.Ctx))
	if err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	return comm.SuccessResponseData(vo)
}
