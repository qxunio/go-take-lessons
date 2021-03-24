package controller

import (
	"github.com/kataras/iris/v12"
	"go-take-lessons/domain/comm"
	"go-take-lessons/domain/comm/ere"
	"go-take-lessons/service"
	"go-take-lessons/tools"
	"go.uber.org/zap"
	"strconv"
)

type SchoolYearController struct {
	Ctx               iris.Context
	SchoolYearService service.SchoolYearService
}

// 查询所有
func (s *SchoolYearController) PostList() *comm.ResultMsg {
	vo, err := s.SchoolYearService.PostList()
	if err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	return comm.SuccessResponseData(vo)
}

// 创建
func (s *SchoolYearController) PostCreate() *comm.ResultMsg {
	schoolYear := s.Ctx.URLParam("schoolYear")
	if tools.IsBlank(schoolYear) {
		return comm.ErrorResponseMsg(ere.StrCommArgIsBlank)
	}
	i, err := strconv.ParseInt(schoolYear, 10, 64)
	if err != nil {
		zap.S().Error("入学年CTL,创建：字符串转数值错误", err.Error())
		return comm.ErrorResponseMsg("入学年只能是数值")
	}
	vo, err := s.SchoolYearService.PostCreate(int(i))
	if err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	return comm.SuccessResponseData(vo)
}

// 删除
func (s *SchoolYearController) PostRemove() *comm.ResultMsg {
	schoolYear := s.Ctx.URLParam("schoolYear")
	if tools.IsBlank(schoolYear) {
		return comm.ErrorResponseMsg(ere.StrCommArgIsBlank)
	}
	if err := s.SchoolYearService.PostRemove(schoolYear); err != nil {
		return comm.ErrorResponseMsg(err.Error())
	}
	return comm.SuccessResponse()
}
