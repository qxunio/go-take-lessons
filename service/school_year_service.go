package service

import (
	"errors"
	"github.com/go-xorm/xorm"
	"go-take-lessons/domain"
	"go-take-lessons/domain/comm"
	"go-take-lessons/domain/comm/ere"
	"go-take-lessons/model"
	"go-take-lessons/tools"
	"go.uber.org/zap"
	"time"
)

type schoolYearService struct {
	mysql *xorm.Engine
}

type SchoolYearService interface {

	// 查询所有
	PostList() (*[]domain.SchoolYearVo, error)

	// 创建
	PostCreate(schoolYear int) (*domain.SchoolYearVo, error)

	// 删除
	PostRemove(schoolYear string) error
}

// 创建
func (s schoolYearService) PostCreate(schoolYear int) (*domain.SchoolYearVo, error) {
	schoolYearModel := new(model.SchoolYear)
	count, err := s.mysql.Where("school_year = ?", schoolYear).And("enable = ?", comm.Enable).Count(schoolYearModel)
	if err != nil {
		zap.S().Error(err)
		return nil, ere.ErrorCommFindError
	}
	if count > 0 {
		return nil, errors.New("入学年已存在")
	}

	schoolYearModel.Id = tools.SnowFlake.Generate().Int64()
	schoolYearModel.SchoolYear = schoolYear
	schoolYearModel.Ctime = time.Now()
	schoolYearModel.Enable = comm.Enable
	if _, err := s.mysql.Insert(schoolYearModel); err != nil {
		zap.S().Error(err)
		return nil, ere.ErrorCommSaveError
	}
	return domain.CvSchoolYearToVo(schoolYearModel), nil
}

// 删除
func (s schoolYearService) PostRemove(schoolYear string) error {
	var schoolYearModel model.SchoolYear
	has, err := s.mysql.Where("school_year = ?", schoolYear).And("enable = ?", comm.Enable).Get(&schoolYearModel)
	if err != nil {
		zap.S().Error(err)
		return ere.ErrorCommDeleteError
	}
	if !has {
		return errors.New("入学年不存在")
	}
	schoolYearModel.Enable = comm.Disable
	schoolYearModel.Utime = time.Now()
	if _, err := s.mysql.Id(schoolYearModel.Id).Cols("enable", "utime").Update(schoolYearModel); err != nil {
		zap.S().Error(err)
		return ere.ErrorCommDeleteError
	}
	return nil
}

// 查询所有
func (s schoolYearService) PostList() (*[]domain.SchoolYearVo, error) {
	var schoolYearModels []model.SchoolYear
	e := s.mysql.Where("enable = ?", comm.Enable).OrderBy("school_year").Find(&schoolYearModels)
	if e != nil {
		zap.S().Error(e)
		return nil, ere.ErrorCommFindError
	}
	vos := make([]domain.SchoolYearVo, 0)
	for _, schoolYearModel := range schoolYearModels {
		vos = append(vos, *domain.CvSchoolYearToVo(&schoolYearModel))
	}
	return &vos, nil
}

func NewSchoolYearService(mysql *xorm.Engine) SchoolYearService {
	return &schoolYearService{
		mysql: mysql,
	}
}
