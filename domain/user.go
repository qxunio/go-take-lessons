package domain

import (
	"go-take-lessons/domain/comm"
	"go-take-lessons/model"
	"strconv"
)

type UserCreateArg struct {
	Name       string `json:"name" valid:"required,isNotBlank~姓名不能为空"`
	Account    string `json:"account" valid:"required,isNotBlank~账号不能为空"`
	Password   string `json:"password" valid:"required,isNotBlank~密码不能为空"`
	Type       string `json:"type" valid:"required,isNotBlank~类型不能为空"`
	Class      string `json:"class"  valid:"-"`
	Uid        string `json:"uid" valid:"required,isNotBlank~请刷新重试"`
	SchoolYear string `json:"schoolYear"  valid:"-"`
}

type UserUpdateArg struct {
	Id         string `json:"id" valid:"required,isNotBlank~ID不能为空"`
	Name       string `json:"name" valid:"required,isNotBlank~用户名不能为空"`
	Type       string `json:"type" valid:"required,isNotBlank~类型不能为空"`
	Class      string `json:"class" valid:"required,isNotBlank~班级编号不能为空"`
	SchoolYear string `json:"schoolYear"  valid:"-"`
}

type UserPageArg struct {
	comm.PageParam `valid:"required"`
	Dest           string `json:"dest" valid:"-"`
	UserType       string `json:"userType" valid:"-"`
	SchoolYear     string `json:"schoolYear" valid:"-"`
	Class          string `json:"class" valid:"-"`
}

type UserVo struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Account    string `json:"account"`
	Type       string `json:"type"`
	Class      string `json:"class"`
	Ctime      string `json:"ctime"`
	SchoolYear string `json:"schoolYear"`
}

type UserSimpleVo struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Account string `json:"account"`
	Tag     string `json:"tag"`
}

func CvUserToSimpleVo(user *model.User) *UserSimpleVo {
	return &UserSimpleVo{
		Id:      strconv.FormatInt(user.Id, 10),
		Name:    user.Name,
		Account: user.Account,
		Tag:     user.Name + "#" + user.Account,
	}
}

func CvUserToVo(user *model.User) *UserVo {
	var class string
	if user.Class == 0 {
		class = ""
	} else {
		class = strconv.FormatUint(uint64(user.Class), 10)
	}
	return &UserVo{
		Id:         strconv.FormatInt(user.Id, 10),
		Name:       user.Name,
		Account:    user.Account,
		Type:       strconv.FormatUint(uint64(user.Type), 10),
		Class:      class,
		Ctime:      user.Ctime.Format(comm.TimeFormatTime),
		SchoolYear: strconv.Itoa(user.SchoolYear),
	}
}

func CvUserToUserRoleModel(user *model.User) *model.UserRole {
	var roleId int64
	if user.Type == comm.Teacher {
		roleId = comm.TeacherId
	} else {
		roleId = comm.StudentId
	}
	return &model.UserRole{
		UserId: user.Id,
		RoleId: roleId,
		Enable: comm.Enable,
	}
}
