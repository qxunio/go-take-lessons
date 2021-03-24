package service

import (
	"encoding/json"
	"errors"
	"github.com/go-xorm/xorm"
	"go-take-lessons/configs"
	"go-take-lessons/db"
	"go-take-lessons/domain"
	"go-take-lessons/domain/comm"
	"go-take-lessons/domain/comm/ere"
	"go-take-lessons/model"
	"go-take-lessons/tools"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type authService struct {
	mysql *xorm.Engine
}

// 刷新TOKEN
func (au authService) PostReload(nowToken string, sessionUser *comm.SessionUSER) (*domain.AuthLoginVo, error) {
	userId, err := tools.ValidateToken(nowToken)
	if err != nil {
		zap.S().Error("认证SERVICE，刷新TOKEN：TOKEN验证失败", err.Error())
		return nil, errors.New("登录过期")
	}
	if tools.StringToInt64(userId) != sessionUser.Id {
		return nil, ere.ErrorCommConvertError
	}
	if _, err := db.RedisClient.Rename(comm.RedisAuthTokenKey+userId, comm.RedisPreAuthTokenKey+userId).Result(); err == nil {
		user, err := au.getByAccount(sessionUser.Account)
		if err != nil {
			return nil, err
		}
		token, err := tools.GenToken(strconv.FormatInt(user.Id, 10))
		if err != nil {
			zap.S().Error("认证SERVICE，刷新TOKEN：生成令牌错误", err.Error())
			return nil, errors.New("生成令牌错误")
		}
		userJsonByte, err := json.Marshal(comm.ConversionSessionUSER(*user))
		if err != nil {
			zap.S().Error("认证SERVICE，刷新TOKEN： JSON序列化登录用户失败", err.Error())
			return nil, ere.ErrorCommUnknown
		}
		_, err = db.RedisClient.Set(comm.RedisAuthTokenKey+strconv.FormatInt(user.Id, 10), string(userJsonByte), time.Second*1800).Result()
		if err != nil {
			zap.S().Error("认证SERVICE，刷新TOKEN： 设置TOKEN缓存失败", err.Error())
			return nil, ere.ErrorCommUnknown
		}
		loginVo := domain.AuthLoginVo{
			Uid:       strconv.FormatInt(user.Id, 10),
			Username:  user.Name,
			Token:     token,
			LoginTime: time.Now().Format(comm.TimeFormatTime),
			Device:    "pc",
		}
		return &loginVo, nil
	}
	zap.S().Error("认证SERVICE，刷新TOKEN：REDIS RENAME 失败", err)
	return nil, errors.New("重载失败")
}

// 登出
func (au authService) PostLogout(sessionUser *comm.SessionUSER) {
	userId := strconv.FormatInt(sessionUser.Id, 10)
	if !tools.IsBlank(userId) {
		_, err := db.RedisClient.Del(comm.RedisAuthTokenKey + strconv.FormatInt(sessionUser.Id, 10)).Result()
		if err != nil {
			zap.S().Error("认证SERVICE，登出：删除登录用户令牌缓存错误", err.Error())
		}
		// 忽略PreToken的处理
	}
}

// 登录
func (au authService) PostLogin(account, password string) (*domain.AuthLoginVo, error) {
	user, err := au.getByAccount(account)
	if err != nil {
		return nil, err
	}
	if tools.Compare(password, user.Salt, user.Password) {
		token, err := tools.GenToken(strconv.FormatInt(user.Id, 10))
		if err != nil {
			zap.S().Error("认证SERVICE，登录：生成用户令牌错误", err.Error())
			return nil, errors.New("生成用户令牌错误")
		}

		userJsonByte, err := json.Marshal(comm.ConversionSessionUSER(*user))
		if err != nil {
			zap.S().Error("认证SERVICE，登录：序列化登录对象错误", err.Error())
			return nil, ere.ErrorCommUnknown
		}
		_, err = db.RedisClient.Set(comm.RedisAuthTokenKey+strconv.FormatInt(user.Id, 10), string(userJsonByte),
			time.Duration(1000*1000*1000*configs.Conf.App.LoginTime)).Result()
		if err != nil {
			zap.S().Error("认证SERVICE，登录：保存令牌到缓存错误", err.Error())
			return nil, ere.ErrorCommUnknown
		}
		loginVo := domain.AuthLoginVo{
			Uid:       strconv.FormatInt(user.Id, 10),
			Username:  user.Name,
			Token:     token,
			LoginTime: time.Now().Format(comm.TimeFormatTime),
			Device:    "pc",
		}
		if user.Type == comm.Student {
			loginVo.Type = "s"
		} else {
			if user.Type == comm.Teacher {
				loginVo.Type = "t"
			} else {
				loginVo.Type = "a"
			}
		}
		return &loginVo, nil
	}
	return nil, errors.New("账号或密码错误")
}

// 根据账号来查询用户
func (au authService) getByAccount(account string) (*model.User, error) {
	var user model.User
	has, err := au.mysql.Where("account = ?", account).And("enable = ?", comm.Enable).Get(&user)
	if err != nil {
		zap.S().Error("认证SERVICE：根据账号来查询用户错误", err.Error())
		return nil, ere.ErrorCommFindError
	}
	if !has {
		zap.S().Error("认证SERVICE：查询不到用户", account)
		return nil, ere.ErrorCommNotFond
	}
	return &user, nil
}

type AuthService interface {
	// 登录
	PostLogin(account, password string) (*domain.AuthLoginVo, error)

	// 登出
	PostLogout(sessionUser *comm.SessionUSER)

	// 刷新TOKEN
	PostReload(nowToken string, sessionUser *comm.SessionUSER) (*domain.AuthLoginVo, error)
}

func NewAuthService(mysql *xorm.Engine) AuthService {
	return &authService{
		mysql: mysql,
	}
}
