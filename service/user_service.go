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
	"go-take-lessons/tools/excel"
	"go.uber.org/zap"
	"mime/multipart"
	"strconv"
	"strings"
	"time"
)

type userService struct {
	mysql       *xorm.Engine
	roleService *RoleService
}

// 根据账号批量查询
func (u userService) ListByUserAccount(account []string) (*[]model.User, error) {
	users, err := u.selectInByAccount(account)
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}
	return &users, nil
}

// 根据UID更新TAG STU
func (u userService) UpdateTagStuByUid(user model.User) {
	if user.Id == 0 {
		return
	}

	var tg model.TagStu
	_, err := u.mysql.Where("uid = ?", user.Id).Get(&tg)
	if err != nil {
		zap.S().Error(err)
		return
	}
	if tg.Id == 0 {
		return
	}

	tg.Class = user.Class
	tg.Name = user.Name
	tg.SchoolYear = user.SchoolYear

	if _, err := u.mysql.Id(tg.Id).Update(tg); err != nil {
		zap.S().Error(err)
	}
}

// 根据入学年和班级查询学生列表
func (u userService) ListStuBySchoolYearClass(schoolYear, class string) *[]model.User {
	var users []model.User
	if tools.IsBlank(schoolYear) || tools.IsBlank(class) {
		return nil
	}
	err := u.mysql.Where("enable = ? and type = ? and school_year = ? and class = ?", comm.Enable, comm.Student, schoolYear, class).Find(&users)
	if err != nil {
		zap.S().Error(err)
	}
	return &users
}

// 导出模板
func (u userService) PostExportTemplate() *excelize.File {
	head := []string{"姓名", "账号", "班级"}
	file := excel.CreateExcelModel(head, "用户导入", "Sheet1")
	return file
}

// 查询所有教师
func (u userService) PostListSimpleTeacher(sessionUser *comm.SessionUSER) (*[]domain.UserSimpleVo, error) {
	if !tools.IsAdmin(sessionUser) {
		zap.S().Error("查询所有教师，当前用户不是管理员", sessionUser)
		return nil, errors.New("非法操作")
	}
	var users []model.User
	err := u.mysql.Where("enable = ? and type = ?", comm.Enable, comm.Teacher).Find(&users)
	if err != nil {
		zap.S().Error(err)
		return nil, ere.ErrorCommFindError
	}
	userSimpleVos := make([]domain.UserSimpleVo, 0)
	for _, user := range users {
		userSimpleVos = append(userSimpleVos, *domain.CvUserToSimpleVo(&user))
	}
	return &userSimpleVos, nil
}

// 重置密码
func (u userService) PostReset(userId string, sessionUser *comm.SessionUSER) error {
	if !tools.IsAdmin(sessionUser) {
		zap.S().Error("重置密码，当前用户不是管理员", sessionUser)
		return errors.New("非法操作")
	}
	if sessionUser.UserType != comm.Admin {
		return errors.New("非法操作")
	}
	var user model.User
	has, err := u.mysql.Id(userId).Get(&user)
	if err != nil {
		zap.S().Error(err)
		return errors.New(" 查询错误")
	}
	if !has {
		return errors.New("查询不到该用户")
	}

	password, salt, err := tools.GenCode("tl123456")
	if err != nil {
		zap.S().Error(err)
		return errors.New("重置失败,生成密码错误")
	}
	user.Password = password
	user.Salt = salt
	user.Utime = time.Now()

	if _, err = u.mysql.Id(user.Id).Update(user); err != nil {
		zap.S().Error(err)
		return errors.New("重置失败！")
	}
	return nil
}

// 导入用户
func (u userService) PostImport(fileReader multipart.File, userType string, sessionUser *comm.SessionUSER, schoolYear string) ([]domain.UserVo, error) {
	if !tools.IsAdmin(sessionUser) {
		zap.S().Error("导入用户，当前用户不是管理员", sessionUser)
		return nil, errors.New("非法操作")
	}
	file, err := excelize.OpenReader(fileReader)
	if err != nil {
		zap.S().Error(err)
		return nil, errors.New("错误，读取文件错误")
	}
	if file == nil {
		zap.S().Error("导入用户，读取文件为空")
		return nil, errors.New("错误，读取文件为空")
	}

	var simpleUsers []model.User
	var account []string
	rows := file.GetRows(file.GetSheetName(1))
	for r, row := range rows {
		if r > 1 {
			var user model.User
			for c, colCell := range row {
				if c == 0 {
					if tools.IsBlank(colCell) {
						return nil, errors.New("失败，有姓名为空")
					}
					user.Name = colCell
				}
				if c == 1 {
					if tools.IsBlank(colCell) {
						return nil, errors.New("失败，有账号为空")
					}
					user.Account = strings.TrimSpace(colCell)
					account = append(account, user.Account)
				}
				if c == 2 && tools.StringToInt(userType) == comm.Student {
					if !tools.IsBlank(colCell) {
						i, err := strconv.ParseInt(colCell, 10, 64)
						if err != nil {
							return nil, errors.New("失败，班级只能是数值")
						}
						user.Class = int(i)
					}
				}
			}
			simpleUsers = append(simpleUsers, user)
		}
	}

	// EXCEL中账号查重
	var repeat []string
	for _, s1 := range account {
		num := 0
		for _, s2 := range account {
			if s1 == s2 {
				num++
				if num >= 2 {
					repeat = append(repeat, s1)
					break
				}
			}
		}
	}

	if len(repeat) != 0 {
		return nil, errors.New("excel中账号重复：" + strings.Join(repeat, ","))
	}

	if account == nil {
		return nil, nil
	}
	existenceUser, err := u.selectBatchByAccount(account)
	if err != nil {
		return nil, err
	}
	var existenceVos []domain.UserVo
	if len(existenceUser) > 0 {
		for _, existenceUser := range existenceUser {
			existenceVos = append(existenceVos, *domain.CvUserToVo(&existenceUser))
		}
		return existenceVos, nil
	}

	userTypeUint := tools.StringToUint8(userType)
	nowTime := time.Now()

	var userRoles []model.UserRole
	var users []model.User

	isStudent := tools.StringToInt(userType) == comm.Student
	for _, user := range simpleUsers {
		password, salt, err := tools.GenCode("tl123456")
		if err != nil {
			zap.S().Error(err)
			return nil, errors.New("创建失败,生成密码错误")
		}
		user.Id = tools.SnowFlake.Generate().Int64()
		user.Type = userTypeUint
		user.Enable = comm.Enable
		user.Creator = sessionUser.Id
		user.Password = password
		user.Salt = salt
		user.Ctime = nowTime
		if isStudent {
			user.SchoolYear = tools.StringToInt(schoolYear)
		} else {
			user.SchoolYear = 0
		}
		users = append(users, user)
		role := domain.CvUserToUserRoleModel(&user)
		role.Ctime = nowTime
		role.Id = tools.SnowFlake.Generate().Int64()
		userRoles = append(userRoles, *role)
	}

	session := u.mysql.NewSession()
	defer session.Close()
	err = session.Begin()
	if err = u.insertBatchByUsers(users, session); err != nil {
		zap.S().Error(err)
		_ = session.Rollback()
		return nil, ere.ErrorCommSaveError
	}

	roleService := *u.roleService
	if err = roleService.InsertBatch(userRoles, session); err != nil {
		zap.S().Error(err)
		_ = session.Rollback()
		return nil, ere.ErrorCommSaveError
	}

	if err = session.Commit(); err != nil {
		zap.S().Error(err)
		return nil, ere.ErrorCommSaveError
	}
	return nil, err
}

// 导出用户
func (u userService) PostExport(sessionUser *comm.SessionUSER) (*excelize.File, error) {
	if !tools.IsAdmin(sessionUser) {
		zap.S().Error("导出用户，当前用户不是管理员", sessionUser)
		return nil, errors.New("非法操作")
	}
	var users []model.User
	if err := u.mysql.Where("enable = ?", comm.Enable).Find(&users); err != nil {
		zap.S().Error(err)
		return nil, ere.ErrorCommFindError
	}

	head := []string{"姓名", "账号", "类型", "加入时间"}
	file := excel.CreateExcelModel(head, "用户导出", "Sheet1")
	contentStyle, _ := file.NewStyle(excel.GetContentCellStyle())

	for index := range users {
		user := users[index]
		subscript := strconv.FormatInt(int64(index+3), 10)
		userType := "教师"
		if user.Type == comm.Student {
			userType = "学生"
		}
		if user.Type == comm.Admin {
			userType = "管理员"
		}
		file.SetCellValue("Sheet1", "A"+subscript, user.Name)
		file.SetCellStyle("Sheet1", "A"+subscript, "A"+subscript, contentStyle)

		file.SetCellValue("Sheet1", "B"+subscript, user.Account)
		file.SetCellStyle("Sheet1", "B"+subscript, "B"+subscript, contentStyle)

		file.SetCellValue("Sheet1", "C"+subscript, userType)
		file.SetCellStyle("Sheet1", "C"+subscript, "C"+subscript, contentStyle)

		file.SetCellValue("Sheet1", "D"+subscript, user.Ctime.Format(comm.TimeFormatTime))
		file.SetCellStyle("Sheet1", "D"+subscript, "D"+subscript, contentStyle)
	}
	// 设置工作簿的默认工作表
	return file, nil
}

// 创建
func (u userService) PostCreate(arg *domain.UserCreateArg, sessionUser *comm.SessionUSER) (*domain.UserVo, error) {
	if !tools.IsAdmin(sessionUser) {
		zap.S().Error("创建用户，当前用户不是管理员", sessionUser)
		return nil, errors.New("非法操作")
	}
	used, err := u.accountUsed(arg.Account)
	if err != nil {
		return nil, ere.ErrorCommSaveError
	}
	if used {
		return nil, errors.New("账号被使用")
	}
	var user model.User
	user.Id = tools.SnowFlake.Generate().Int64()
	user.Account = arg.Account
	password, salt, err := tools.GenCode(arg.Password)
	if err != nil {
		zap.S().Error(err)
		return nil, ere.ErrorCommSaveError
	}
	nowTime := time.Now()
	user.Password = password
	user.Salt = salt
	user.Name = arg.Name
	user.Creator = sessionUser.Id
	user.Enable = comm.Enable
	user.Ctime = nowTime
	user.Type = tools.StringToUint8(arg.Type)
	if tools.StringToInt(arg.Type) == comm.Student {
		user.Class = tools.StringToInt(arg.Class)
		user.SchoolYear = tools.StringToInt(arg.SchoolYear)
	} else {
		user.SchoolYear = 0
		user.Class = 0
	}

	session := u.mysql.NewSession()
	defer session.Close()
	err = session.Begin()
	if _, err = session.Insert(user); err != nil {
		_ = session.Rollback()
		zap.S().Error(err)
		return nil, ere.ErrorCommSaveError
	}

	role := domain.CvUserToUserRoleModel(&user)
	role.Ctime = nowTime
	role.Id = tools.SnowFlake.Generate().Int64()
	roleService := *u.roleService
	if err := roleService.CreateUserRole(role, session); err != nil {
		_ = session.Rollback()
		zap.S().Error(err)
		return nil, ere.ErrorCommSaveError
	}
	if err = session.Commit(); err != nil {
		zap.S().Error(err)
		return nil, ere.ErrorCommSaveError
	}
	return domain.CvUserToVo(&user), nil
}

// 修改
func (u userService) PostUpdate(arg *domain.UserUpdateArg, sessionUser *comm.SessionUSER) error {
	if !tools.IsAdmin(sessionUser) {
		zap.S().Error("修改用户，当前用户不是管理员", sessionUser)
		return errors.New("非法操作")
	}
	var user model.User
	user.Id = tools.StringToInt64(arg.Id)
	user.Name = arg.Name
	user.Type = tools.StringToUint8(arg.Type)
	user.Utime = time.Now()
	if tools.StringToInt(arg.Type) == comm.Student {
		user.SchoolYear = tools.StringToInt(arg.SchoolYear)
		user.Class = tools.StringToInt(arg.Class)
	} else {
		user.SchoolYear = 0
		user.Class = 0
	}

	session := u.mysql.NewSession()
	err := session.Begin()
	if _, err = session.ID(user.Id).Update(user); err != nil {
		if err = session.Rollback(); err != nil {
			zap.S().Error(err)
		}
		return ere.ErrorCommUpdateError
	}
	role := domain.CvUserToUserRoleModel(&user)
	role.Utime = time.Now()
	roleService := *u.roleService
	if err = roleService.UpdateUserRoleByUserId(role, session); err != nil {
		_ = session.Rollback()
		zap.S().Error(err)
		return ere.ErrorCommUpdateError
	}
	if err = session.Commit(); err != nil {
		zap.S().Error(err)
		return ere.ErrorCommUpdateError
	}
	if user.Type == comm.Student {
		u.UpdateTagStuByUid(user)
	}
	return err
}

// 分页
func (u userService) PostPage(arg *domain.UserPageArg, sessionUser *comm.SessionUSER) (interface{}, error) {
	if !tools.IsAdmin(sessionUser) {
		zap.S().Error("分页用户，当前用户不是管理员", sessionUser)
		return nil, errors.New("非法操作")
	}
	userSql := u.mysql.Where("enable = ?", comm.Enable)
	countSql := u.mysql.Where("enable = ?", comm.Enable)
	if !tools.IsBlank(arg.UserType) {
		sql1 := "type = "
		if tools.StringToInt(arg.UserType) == comm.Teacher {
			sql1 = sql1 + strconv.Itoa(comm.Teacher)
		} else {
			sql1 = sql1 + strconv.Itoa(comm.Student)
		}
		userSql.And(sql1)
		countSql.And(sql1)
	}

	if !tools.IsBlank(arg.Dest) {
		arg := "'%" + arg.Dest + "%'"
		sql2 := "account like " + arg + " or name like" + arg
		userSql.And(sql2)
		countSql.And(sql2)
	}

	if !tools.IsBlank(arg.Class) {
		sql3 := "class = '" + arg.Class + "'"
		userSql.And(sql3)
		countSql.And(sql3)
	}
	if !tools.IsBlank(arg.SchoolYear) {
		sql4 := "school_year = '" + arg.SchoolYear + "'"
		userSql.And(sql4)
		countSql.And(sql4)
	}

	var um model.User
	var total int64
	total, err := countSql.Count(um)
	if err != nil {
		zap.S().Error(err)
		return nil, ere.ErrorCommFindError
	}
	var pageVo comm.PageVo
	pageVo.TotalCount = total
	var users []model.User
	if err := userSql.Limit(arg.GetLimit(), arg.GetOffset()).Desc("ctime", "id").Find(&users); err != nil {
		zap.S().Error(err)
		return nil, ere.ErrorCommFindError
	}
	vos := make([]domain.UserVo, 0)
	for _, user := range users {
		vos = append(vos, *domain.CvUserToVo(&user))
	}
	pageVo.Data = vos
	return pageVo, nil
}

// 删除
func (u userService) PostRemove(userId string, sessionUser *comm.SessionUSER) error {
	if !tools.IsAdmin(sessionUser) {
		zap.S().Error("删除用户，当前用户不是管理员", sessionUser)
		return errors.New("非法操作")
	}
	var user model.User
	has, err := u.mysql.Id(userId).Get(&user)
	if err != nil {
		zap.S().Error(err)
		return ere.ErrorCommFindError
	}
	if !has {
		return errors.New("用户不存在")
	}

	user.Enable = comm.Disable
	user.Utime = time.Now()
	session := u.mysql.NewSession()
	defer session.Close()
	err = session.Begin()
	if _, err = session.ID(userId).Cols("enable", "utime").Update(&user); err != nil {
		_ = session.Rollback()
		zap.S().Error(err)
		return ere.ErrorCommUpdateError
	}
	roleService := *u.roleService
	if err = roleService.RemoveUserRoleByUserId(userId, session); err != nil {
		_ = session.Rollback()
		zap.S().Error(err)
		return ere.ErrorCommUpdateError
	}
	if err = session.Commit(); err != nil {
		zap.S().Error(err)
		return ere.ErrorCommDeleteError
	}
	return nil
}

// 检查账号是否被使用
func (u userService) accountUsed(account string) (bool, error) {
	user := new(model.User)
	i, e := u.mysql.Where("account = ?", account).And("enable = ?", comm.Enable).Count(user)
	if e != nil {
		zap.S().Error(e)
		return false, ere.ErrorCommFindError
	}
	return i > 0, nil
}

// 分批插入
func (u userService) insertBatchByUsers(users []model.User, session *xorm.Session) error {
	length := len(users)
	if length > comm.BatchFactor {
		c := length / comm.BatchFactor
		m := length % comm.BatchFactor
		index := 0

		for i := 1; i <= c; i++ {
			if _, err := session.Insert(users[index : comm.BatchFactor*i]); err != nil {
				zap.S().Error(err)
				return ere.ErrorCommSaveError
			}
			index = comm.BatchFactor * i
		}
		if m != 0 {
			if _, err := session.Insert(users[index:]); err != nil {
				zap.S().Error(err)
				return ere.ErrorCommSaveError
			}
		}
	} else {
		if _, err := session.Insert(users); err != nil {
			return ere.ErrorCommSaveError
		}
		return nil
	}
	return nil
}

// 分批账号查询
func (u userService) selectBatchByAccount(accounts []string) ([]model.User, error) {
	var existence []model.User
	var err error
	length := len(accounts)
	if length > comm.BatchFactor {
		c := length / comm.BatchFactor
		m := length % comm.BatchFactor
		index := 0
		for i := 1; i <= c; i++ {
			dbUser, err := u.selectInByAccount(accounts[index : comm.BatchFactor*i])
			if err != nil {
				zap.S().Error(err)
				return nil, err
			}
			for _, v := range dbUser {
				existence = append(existence, v)
			}
			index = comm.BatchFactor * i
		}
		if m != 0 {
			dbUser, err := u.selectInByAccount(accounts[index:])
			if err != nil {
				zap.S().Error(err)
				return nil, err
			}
			for _, v := range dbUser {
				existence = append(existence, v)
			}
		}
	} else {
		existence, err = u.selectInByAccount(accounts)
		if err != nil {
			zap.S().Error(err)
			return nil, err
		}
	}
	return existence, nil
}

// In查询
func (u userService) selectInByAccount(accounts []string) ([]model.User, error) {
	var dbUser []model.User
	if err := u.mysql.Cols("id", "name", "account", "school_year", "class").Where("enable = ?", comm.Enable).In("account", accounts).Find(&dbUser); err != nil {
		zap.S().Error(err)
		return nil, ere.ErrorCommFindError
	}
	return dbUser, nil
}

type UserService interface {
	// 创建
	PostCreate(arg *domain.UserCreateArg, sessionUser *comm.SessionUSER) (*domain.UserVo, error)

	//修改
	PostUpdate(arg *domain.UserUpdateArg, sessionUser *comm.SessionUSER) error

	// 分页
	PostPage(arg *domain.UserPageArg, sessionUser *comm.SessionUSER) (interface{}, error)

	// 删除
	PostRemove(userId string, sessionUser *comm.SessionUSER) error

	// 导出用户
	PostExport(sessionUser *comm.SessionUSER) (*excelize.File, error)

	// 导入用户
	PostImport(fileReader multipart.File, userType string, sessionUser *comm.SessionUSER, schoolYear string) ([]domain.UserVo, error)

	// 重置密码
	PostReset(userId string, sessionUser *comm.SessionUSER) error

	// 查询所有教师
	PostListSimpleTeacher(sessionUser *comm.SessionUSER) (*[]domain.UserSimpleVo, error)

	// 导出模板
	PostExportTemplate() *excelize.File

	// 根据入学年和班级查询学生列表
	ListStuBySchoolYearClass(schoolYear, class string) *[]model.User

	// 根据UID更新TAG STU
	UpdateTagStuByUid(user model.User)

	// 根据账号批量查询
	ListByUserAccount(account []string) (*[]model.User, error)
}

func NewUserService(mysql *xorm.Engine, roleService *RoleService) UserService {
	return &userService{
		mysql:       mysql,
		roleService: roleService,
	}
}
