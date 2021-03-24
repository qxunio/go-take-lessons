package controller

import (
	"github.com/kataras/iris/v12"
	"go-take-lessons/domain"
	"go-take-lessons/domain/comm"
	"go-take-lessons/domain/comm/ere"
	"go-take-lessons/service"
	"go-take-lessons/tools"
	"go.uber.org/zap"
	"strings"
)

type TagController struct {
	Ctx        iris.Context
	TagService service.TagService
}

// 查询
func (t *TagController) PostList() *comm.ResultMsg {
	vo, err := t.TagService.PostList()
	if err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	return comm.SuccessResponseData(vo)
}

// 创建
func (t *TagController) PostCreate() *comm.ResultMsg {
	tag := t.Ctx.URLParam("tag")
	if tools.IsBlank(tag) {
		return comm.ErrorResponseMsg(ere.StrCommArgIsBlank)
	}
	vo, err := t.TagService.PostCreate(tag, tools.GetSessionUser(&t.Ctx))
	if err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	return comm.SuccessResponseData(vo)
}

// 删除
func (t *TagController) PostRemove() *comm.ResultMsg {
	tag := t.Ctx.URLParam("tag")
	if tools.IsBlank(tag) {
		return comm.ErrorResponseMsg(ere.StrCommArgIsBlank)
	}
	err := t.TagService.PostRemove(tag, tools.GetSessionUser(&t.Ctx))
	if err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	return comm.SuccessResponse()
}

// 根据班级创建
func (t *TagController) PostCreateStuByClass() *comm.ResultMsg {
	var arg domain.TagStuCreateByClassArg
	if err := tools.GetRequestJson(&t.Ctx, &arg); err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	if err := t.TagService.PostCreateStuByClass(arg, tools.GetSessionUser(&t.Ctx)); err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	return comm.SuccessResponse()
}

// 分页
func (t *TagController) PostPage() *comm.ResultMsg {
	var arg domain.TagStuPageArg
	if err := tools.GetRequestJson(&t.Ctx, &arg); err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}

	vo, err := t.TagService.PostPage(&arg, tools.GetSessionUser(&t.Ctx))
	if err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	return comm.SuccessResponseData(vo)
}

// 删除学生
func (t *TagController) PostRemoveTagStu() *comm.ResultMsg {
	var arg domain.TagStuRemoveArg
	if err := tools.GetRequestJson(&t.Ctx, &arg); err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	err := t.TagService.PostRemoveTagStu(&arg, tools.GetSessionUser(&t.Ctx))
	if err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	return comm.SuccessResponse()
}

// 批量创建
func (t *TagController) PostCreateTagStuList() *comm.ResultMsg {
	var arg domain.TagStuCreateListArg
	if err := tools.GetRequestJson(&t.Ctx, &arg); err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	err := t.TagService.PostCreateTagStuList(&arg, tools.GetSessionUser(&t.Ctx))
	if err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	return comm.SuccessResponse()
}

// 导出模板
func (t *TagController) PostExportTemplate() {
	file := t.TagService.PostExportTemplate()
	if err := file.Write(t.Ctx.ResponseWriter()); err != nil {
		zap.S().Error("标签CTL,导出模板: 写出响应错误", err.Error())
		_, _ = t.Ctx.JSON("导出失败")
		return
	}
}

// 导入用户
func (t *TagController) PostImport() *comm.ResultMsg {
	fileReader, fileInfo, err := t.Ctx.FormFile("file")
	if err != nil {
		return comm.ErrorResponseMsg("上传失败")
	}

	tagId := t.Ctx.FormValue("tagId")
	if tools.IsBlank(tagId) {
		return comm.ErrorResponseMsg("上传的用户标签未选择")
	}

	if fileInfo.Size > 3145728 {
		return comm.ErrorResponseMsg("上传文件大小不能超过3M")
	}
	if strings.HasSuffix(fileInfo.Filename, ".xlxs") {
		return comm.ErrorResponseMsg("上传的文件类型错误")
	}
	err = t.TagService.PostImport(fileReader, tagId, tools.GetSessionUser(&t.Ctx))
	if err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	return comm.SuccessResponseMsg("导入成功")
}
