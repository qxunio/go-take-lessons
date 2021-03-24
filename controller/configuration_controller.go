package controller

import (
	"github.com/kataras/iris/v12"
	"go-take-lessons/domain"
	"go-take-lessons/domain/comm"
	"go-take-lessons/domain/comm/ere"
	"go-take-lessons/service"
	"go-take-lessons/tools"
)

type ConfigurationController struct {
	Ctx                  iris.Context
	ConfigurationService service.ConfigurationService
}

// 学生查询自己历史选课记录
func (c *ConfigurationController) PostStuHistoryStuSelected() *comm.ResultMsg {
	var arg domain.IdArg
	if err := tools.GetRequestJson(&c.Ctx, &arg); err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	vo, err := c.ConfigurationService.PostStuHistoryStuSelected(arg.Id, tools.GetSessionUser(&c.Ctx))
	if err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	return comm.SuccessResponseData(vo)
}

// 根据选课事件ID查询选课课程 学生
func (c *ConfigurationController) PostListStu() *comm.ResultMsg {
	eventId := c.Ctx.URLParam("eventId")
	if tools.IsBlank(eventId) {
		return comm.ErrorResponseMsg(ere.StrCommArgIsBlank)
	}
	return comm.SuccessResponseData(c.ConfigurationService.PostListStu(eventId))
}

// 根据选课事件ID查询所有的选课课程
func (c *ConfigurationController) PostList() *comm.ResultMsg {
	eventId := c.Ctx.URLParam("eId")
	if tools.IsBlank(eventId) {
		return comm.ErrorResponseMsg(ere.StrCommArgIsBlank)
	}
	vo, err := c.ConfigurationService.PostList(eventId, tools.GetSessionUser(&c.Ctx))
	if err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	return comm.SuccessResponseData(vo)
}

// 创建选课课程
func (c *ConfigurationController) PostCreate() *comm.ResultMsg {
	var arg domain.ConfigurationCreateArg
	if err := tools.GetRequestJson(&c.Ctx, &arg); err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	if err := c.ConfigurationService.PostCreate(&arg); err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	return comm.SuccessResponse()
}

// 更新选课课程
func (c *ConfigurationController) PostUpdate() *comm.ResultMsg {
	var arg domain.ConfigurationUpdateArg
	if err := tools.GetRequestJson(&c.Ctx, &arg); err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	if err := c.ConfigurationService.PostUpdate(&arg); err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	return comm.SuccessResponse()
}

// 删除选课课程,级联删除该课程下的教师
func (c *ConfigurationController) PostRemove() *comm.ResultMsg {
	var arg domain.ConfigurationEventIdAndCsIdArg
	if err := tools.GetRequestJson(&c.Ctx, &arg); err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	if err := c.ConfigurationService.PostRemove(arg.EventId, arg.CsId, tools.GetSessionUser(&c.Ctx)); err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	return comm.SuccessResponse()
}

// 删除选课教师配置
func (c *ConfigurationController) PostRemoveTeacher() *comm.ResultMsg {
	configTeacherId := c.Ctx.URLParam("id")
	if tools.IsBlank(configTeacherId) {
		return comm.ErrorResponseMsg(ere.StrCommArgIsBlank)
	}
	if err := c.ConfigurationService.PostRemoveTeacher(configTeacherId, tools.GetSessionUser(&c.Ctx)); err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	return comm.SuccessResponse()
}

// 批量创建
func (c *ConfigurationController) PostCreateBatch() *comm.ResultMsg {
	var arg domain.ConfigurationCreateBatchArg
	if err := tools.GetRequestJson(&c.Ctx, &arg); err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	if err := c.ConfigurationService.PostCreateBatch(arg, tools.GetSessionUser(&c.Ctx)); err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	return comm.SuccessResponse()
}

//查询选课下的课程详情
func (c *ConfigurationController) PostCsDetails() *comm.ResultMsg {
	var arg domain.IdArg
	if err := tools.GetRequestJson(&c.Ctx, &arg); err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	vo, err := c.ConfigurationService.PostCsDetails(arg.Id, tools.GetSessionUser(&c.Ctx))
	if err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	return comm.SuccessResponseData(vo)
}

//查询选课下的课程详情的学生列表
func (c *ConfigurationController) PostPageCsDetailsStu() *comm.ResultMsg {
	var arg domain.StuSubjectPageArg
	if err := tools.GetRequestJson(&c.Ctx, &arg); err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	vo, err := c.ConfigurationService.PostPageCsDetailsStu(arg, tools.GetSessionUser(&c.Ctx))
	if err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	return comm.SuccessResponseData(vo)
}

// 管理员追加学生到课堂 搜索学生
func (c *ConfigurationController) PostAppendStuSearch() *comm.ResultMsg {
	var arg domain.ConfigurationEventIdAndDestArg
	if err := tools.GetRequestJson(&c.Ctx, &arg); err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	vo, err := c.ConfigurationService.PostAppendStuSearch(arg, tools.GetSessionUser(&c.Ctx))
	if err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	return comm.SuccessResponseData(vo)
}
