package service

import (
	"github.com/go-xorm/xorm"
	"go-take-lessons/domain/comm"
	"go-take-lessons/domain/comm/ere"
	"go-take-lessons/model"
	"go-take-lessons/tools"
	"go.uber.org/zap"
	"time"
)

type roleService struct {
	mysql *xorm.Engine
}

type RoleService interface {
	// 创建用户对应的角色
	CreateUserRole(userRole *model.UserRole, session *xorm.Session) error

	// 更新用户对应的角色
	UpdateUserRoleByUserId(userRole *model.UserRole, session *xorm.Session) error

	// 根据userId删除
	RemoveUserRoleByUserId(userId string, session *xorm.Session) error

	// 批量插入
	InsertBatch(userRole []model.UserRole, session *xorm.Session) error
}

// 批量插入
func (r *roleService) InsertBatch(userRole []model.UserRole, session *xorm.Session) error {
	length := len(userRole)
	if length > comm.BatchFactor {
		c := length / comm.BatchFactor
		m := length % comm.BatchFactor
		index := 0

		for i := 1; i <= c; i++ {
			if _, err := session.Insert(userRole[index : comm.BatchFactor*i]); err != nil {
				zap.S().Error(err)
				return ere.ErrorCommSaveError
			}
			index = comm.BatchFactor * i
		}

		if m != 0 {
			if _, err := session.Insert(userRole[index:]); err != nil {
				zap.S().Error(err)
				return ere.ErrorCommSaveError
			}
		}
	} else {
		if _, err := session.Insert(userRole); err != nil {
			return ere.ErrorCommSaveError
		}
		return nil
	}
	return nil
}

// 创建用户对应的角色
func (r *roleService) CreateUserRole(userRole *model.UserRole, session *xorm.Session) error {
	if session == nil {
		if _, err := r.mysql.Insert(userRole); err != nil {
			zap.S().Error(err)
			return ere.ErrorCommSaveError
		}
	} else {
		if _, err := session.Insert(userRole); err != nil {
			zap.S().Error(err)
			return ere.ErrorCommSaveError
		}
	}
	return nil
}

// 更新用户对应的角色
func (r *roleService) UpdateUserRoleByUserId(userRole *model.UserRole, session *xorm.Session) error {
	if session == nil {
		if _, err := r.mysql.Update(userRole, &model.UserRole{UserId: userRole.UserId}); err != nil {
			zap.S().Error(err)
			return ere.ErrorCommUpdateError
		}
	} else {
		if _, err := session.Update(userRole, &model.UserRole{UserId: userRole.UserId}); err != nil {
			zap.S().Error(err)
			return ere.ErrorCommUpdateError
		}
	}
	return nil
}

// 移除用户对应的角色
func (r *roleService) RemoveUserRoleByUserId(userId string, session *xorm.Session) error {
	var userRole model.UserRole
	userRole.Enable = comm.Disable
	userRole.UserId = tools.StringToInt64(userId)
	userRole.Utime = time.Now()

	if session == nil {
		if _, err := r.mysql.Cols("enable", "utime").Update(&userRole, &model.UserRole{UserId: userRole.UserId}); err != nil {
			zap.S().Error(err)
			return ere.ErrorCommUpdateError
		}
	} else {
		if _, err := session.Cols("enable", "utime").Update(&userRole, &model.UserRole{UserId: userRole.UserId}); err != nil {
			zap.S().Error(err)
			return ere.ErrorCommUpdateError
		}
	}
	return nil
}

func NewRoleService(mysql *xorm.Engine) RoleService {
	return &roleService{mysql: mysql}
}
