package service

import (
	"errors"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/go-xorm/xorm"
	"go-take-lessons/domain"
	"go-take-lessons/domain/comm"
	"go-take-lessons/domain/comm/ere"
	"go-take-lessons/model"
	"go-take-lessons/tools"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
)

type subjectService struct {
	mysql       *xorm.Engine
	sessionUser *comm.SessionUSER
}

type SubjectService interface {
	// 创建
	PostCreate(arg *domain.SubjectCreateArg, sessionUser *comm.SessionUSER) (*domain.SubjectVo, error)

	// 分页
	PostPage(arg *comm.PageParam, sessionUser *comm.SessionUSER) (interface{}, error)

	// 修改
	PostUpdate(arg *domain.SubjectUpdateArg) error

	// 删除
	PostRemove(subjectId string, sessionUser *comm.SessionUSER) error

	// 导出学科
	PostExport() (*excelize.File, error)

	// 查询所有学科
	PostListSimpleSubject(sessionUser *comm.SessionUSER) (*[]domain.SubjectSimpleVo, error)
}

// 查询所有学科
func (s subjectService) PostListSimpleSubject(sessionUser *comm.SessionUSER) (*[]domain.SubjectSimpleVo, error) {
	var subjects []model.Subject
	err := s.mysql.Where("enable = ?", comm.Enable).Find(&subjects)
	if err != nil {
		zap.S().Error(err)
		return nil, ere.ErrorCommFindError
	}
	subjectSimpleVo := make([]domain.SubjectSimpleVo, 0)
	for _, subject := range subjects {
		subjectSimpleVo = append(subjectSimpleVo, *domain.CvSubjectToSimpleVo(&subject))
	}
	return &subjectSimpleVo, nil
}

// 导出学科
func (s subjectService) PostExport() (*excelize.File, error) {
	var subjects []model.Subject
	if err := s.mysql.Where("enable = ?", comm.Enable).Find(&subjects); err != nil {
		zap.S().Error(err)
		return nil, ere.ErrorCommFindError
	}

	file := excelize.NewFile()
	// 创建一个工作表
	index := file.NewSheet("Sheet1")
	file.SetCellValue("Sheet1", "A1", "学科")
	file.SetCellValue("Sheet1", "B1", "介绍")
	file.SetCellValue("Sheet1", "C1", "创建时间")

	for index := range subjects {
		subject := subjects[index]
		subscript := strconv.FormatInt(int64(index+2), 10)
		file.SetCellValue("Sheet1", "A"+subscript, subject.Name)
		file.SetCellValue("Sheet1", "B"+subscript, subject.Introduction)
		file.SetCellValue("Sheet1", "C"+subscript, subject.Ctime)
	}
	// 设置工作簿的默认工作表
	file.SetActiveSheet(index)
	return file, nil
}

// 删除
func (s subjectService) PostRemove(subjectId string, sessionUser *comm.SessionUSER) error {
	if sessionUser.UserType != comm.Admin {
		return errors.New("非法操作")
	}

	var subject model.Subject
	has, err := s.mysql.Id(subjectId).Get(&subject)
	if err != nil {
		zap.S().Error(err)
		return ere.ErrorCommFindError
	}
	if !has {
		return errors.New("学科不存在")
	}
	subject.Enable = comm.Disable
	subject.Utime = time.Now()
	if _, err := s.mysql.Id(subjectId).Cols("enable", "utime").Update(&subject); err != nil {
		zap.S().Error(err)
		return ere.ErrorCommDeleteError
	}
	return nil
}

// 修改
func (s subjectService) PostUpdate(arg *domain.SubjectUpdateArg) error {
	var subject model.Subject
	subject.Id = tools.StringToInt64(arg.Id)
	subject.Name = arg.Name
	subject.Introduction = arg.Introduction
	subject.Utime = time.Now()
	if _, err := s.mysql.Id(subject.Id).Update(subject); err != nil {
		zap.S().Error(err)
		return ere.ErrorCommUpdateError
	}
	return nil
}

// 分页
func (s subjectService) PostPage(arg *comm.PageParam, sessionUser *comm.SessionUSER) (interface{}, error) {
	subjectSql := s.mysql.Where("enable = ?", comm.Enable)
	countSql := s.mysql.Where("enable = ?", comm.Enable)

	var subjects []model.Subject
	if err := subjectSql.Limit(arg.GetLimit(), arg.GetOffset()).Desc("ctime", "id").Find(&subjects); err != nil {
		zap.S().Error(err)
		return nil, ere.ErrorCommFindError
	}

	var selectSubject model.Subject
	total, err := countSql.Count(selectSubject)
	if err != nil {
		zap.S().Error(err)
		return nil, ere.ErrorCommFindError
	}

	var pageVo comm.PageVo
	pageVo.TotalCount = total

	subjectsVo := make([]domain.SubjectVo, 0)
	for _, subject := range subjects {
		subjectsVo = append(subjectsVo, *domain.CvSubjectToVo(&subject))
	}
	pageVo.Data = subjectsVo
	return &pageVo, nil
}

// 创建
func (s subjectService) PostCreate(arg *domain.SubjectCreateArg, session *comm.SessionUSER) (*domain.SubjectVo, error) {
	used, err := s.nameUsed(arg.Name)
	if err != nil {
		return nil, ere.ErrorCommSaveError
	}
	if used {
		return nil, errors.New("学科名被使用")
	}
	var subject model.Subject
	subject.Id = tools.SnowFlake.Generate().Int64()
	subject.Name = strings.TrimSpace(arg.Name)
	subject.Introduction = strings.TrimSpace(arg.Introduction)
	subject.Ctime = time.Now()
	subject.Creator = session.Id
	subject.Enable = comm.Enable
	if _, err := s.mysql.Insert(subject); err != nil {
		zap.S().Error(err)
		return nil, ere.ErrorCommSaveError
	}
	return domain.CvSubjectToVo(&subject), nil
}

// 检查学科名是否被使用
func (s subjectService) nameUsed(name string) (bool, error) {
	subject := new(model.Subject)
	count, err := s.mysql.Where("name = ?", name).And("enable = ?", comm.Enable).Count(subject)
	if err != nil {
		zap.S().Error(err)
		return false, ere.ErrorCommFindError
	}
	return count > 0, nil
}

func NewSubjectService(mysql *xorm.Engine) SubjectService {
	return &subjectService{mysql: mysql}
}
