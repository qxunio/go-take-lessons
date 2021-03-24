package service

import (
	"errors"
	"github.com/go-xorm/xorm"
	"go-take-lessons/domain"
	"go-take-lessons/domain/comm"
	"go-take-lessons/domain/comm/ere"
	"go-take-lessons/model"
	"go.uber.org/zap"
)

type menuService struct {
	mysql *xorm.Engine
}

type MenuService interface {
	// 获取Router
	PostRouter(sessionUser *comm.SessionUSER) (error, interface{})
}

func (m menuService) PostRouter(sessionUser *comm.SessionUSER) (error, interface{}) {
	var roleIds []int64
	if err := m.mysql.Table("user_role").Cols("role_id").Where("user_id = ? ", sessionUser.Id).And("enable = ?", comm.Enable).Find(&roleIds); err != nil {
		zap.S().Error(err)
		return ere.ErrorCommFindError, nil
	}
	if roleIds == nil {
		return errors.New("当前用户没有身份认证"), nil
	}

	// select * from menu where id in (select menu_id from role_menu where role_id in (?))

	var menuIds []int64
	if err := m.mysql.Table("role_menu").Cols("menu_id").In("role_id", roleIds).And("enable = ?", comm.Enable).Find(&menuIds); err != nil {
		zap.S().Error(err)
		return ere.ErrorCommFindError, nil
	}

	var menus []model.Menu
	if err := m.mysql.In("id", menuIds).And("enable = ?", comm.Enable).Asc("priority").Find(&menus); err != nil {
		zap.S().Error(err)
		return ere.ErrorCommFindError, nil
	}

	return nil, domain.CvMenuToVo(&menus)
}

func NewMenuService(mysql *xorm.Engine) MenuService {
	return &menuService{
		mysql: mysql,
	}
}
