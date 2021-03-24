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
	"strings"
	"time"
)

type tagService struct {
	mysql       *xorm.Engine
	userService *UserService
}

// 导入用户
func (t tagService) PostImport(fileReader multipart.File, tagId string, user *comm.SessionUSER) error {
	if !tools.IsAdmin(user) {
		zap.S().Error("标签导入用户，当前用户不是管理员", user)
		return errors.New("非法操作")
	}
	file, err := excelize.OpenReader(fileReader)
	if err != nil {
		zap.S().Error(err)
		return errors.New("错误，读取文件错误")
	}
	if file == nil {
		zap.S().Error("标签 导入用户，文件为空")
		return errors.New("错误，读取文件为空")
	}

	var account []string
	rows := file.GetRows(file.GetSheetName(1))
	for r, row := range rows {
		if r > 1 {
			var user model.User
			for c, colCell := range row {
				if c == 1 {
					if tools.IsBlank(colCell) {
						return errors.New("失败，有账号为空")
					}
					user.Account = strings.TrimSpace(colCell)
					account = append(account, colCell)
				}
			}
		}
	}

	if account == nil || len(account) == 0 {
		return nil
	}

	userService := *t.userService

	userAll, err := userService.ListByUserAccount(account)
	if err != nil {
		return ere.ErrorCommFindError
	}

	var tagStu []model.TagStu

	if err := t.mysql.Where("tag_id = ?", tagId).Find(&tagStu); err != nil {
		zap.S().Error(err)
		return ere.ErrorCommSaveError
	}

	var stus = t.deduplication(userAll, &tagStu)
	now := time.Now()
	var newTagStu []model.TagStu
	for _, stu := range *stus {
		var t model.TagStu
		t.Id = tools.SnowFlake.Generate().Int64()
		t.Ctime = now
		t.Uid = stu.Id
		t.TagId = tools.StringToInt64(tagId)
		t.Creator = user.Id
		t.Name = stu.Name
		t.SchoolYear = stu.SchoolYear
		t.Class = stu.Class
		t.Account = stu.Account
		newTagStu = append(newTagStu, t)
	}
	if _, err = t.mysql.Insert(newTagStu); err != nil {
		return ere.ErrorCommSaveError
	}
	return nil
}

// 导出模板
func (t tagService) PostExportTemplate() *excelize.File {
	head := []string{"姓名", "账号"}
	file := excel.CreateExcelModel(head, "标签学生导入", "Sheet1")
	return file
}

// 批量创建
func (t tagService) PostCreateTagStuList(arg *domain.TagStuCreateListArg, user *comm.SessionUSER) error {
	if !tools.IsAdmin(user) {
		zap.S().Error("标签批量创建，当前用户不是管理员", user)
		return errors.New("非法操作")
	}
	var tagStu []model.TagStu

	if err := t.mysql.Where("tag_id = ?", arg.TagId).Find(&tagStu); err != nil {
		zap.S().Error(err)
		return ere.ErrorCommSaveError
	}

	var newTagStu []domain.TagStuCreateArg

	if tagStu == nil || len(tagStu) == 0 {
		newTagStu = arg.StuList
	} else {
		for _, a := range arg.StuList {
			exist := false
			for _, db := range tagStu {
				if db.Uid == tools.StringToInt64(a.Uid) {
					exist = true
				}
			}
			if !exist {
				newTagStu = append(newTagStu, a)
			}
		}
	}

	tgId := tools.StringToInt64(arg.TagId)
	var newTagStuDb []model.TagStu
	for _, data := range newTagStu {
		newTagStuDb = append(newTagStuDb, *domain.CvTagStuCreateArgToTagStu(&data, tgId, user.Id))
	}

	if _, err := t.mysql.Insert(newTagStuDb); err != nil {
		zap.S().Error(err)
		return ere.ErrorCommSaveError
	}
	return nil
}

// 移除学生
func (t tagService) PostRemoveTagStu(arg *domain.TagStuRemoveArg, user *comm.SessionUSER) error {
	if !tools.IsAdmin(user) {
		zap.S().Error("标签移除学生，当前用户不是管理员", user)
		return errors.New("非法操作")
	}
	var tg model.TagStu

	if _, err := t.mysql.Where("tag_id = ? and uid = ?", arg.TagId, arg.Uid).Delete(tg); err != nil {
		zap.S().Error(err)
		return ere.ErrorCommDeleteError
	}
	return nil
}

// 分页
func (t tagService) PostPage(arg *domain.TagStuPageArg, user *comm.SessionUSER) (interface{}, error) {
	if !tools.IsAdmin(user) {
		zap.S().Error("标签分页，当前用户不是管理员", user)
		return nil, errors.New("非法操作")
	}
	userSql := t.mysql.Where("tag_id = ?", arg.TagId)
	countSql := t.mysql.Where("tag_id = ?", arg.TagId)

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

	var um model.TagStu
	var total int64
	total, err := countSql.Count(um)
	if err != nil {
		zap.S().Error(err)
		return nil, ere.ErrorCommFindError
	}

	var pageVo comm.PageVo
	pageVo.TotalCount = total
	var users []model.TagStu
	if err := userSql.Limit(arg.GetLimit(), arg.GetOffset()).Desc("ctime", "id").Find(&users); err != nil {
		zap.S().Error(err)
		return nil, ere.ErrorCommFindError
	}

	vos := make([]domain.TagStuVo, 0)
	for _, user := range users {
		vos = append(vos, *domain.CvTagStuToTagStuVo(&user))
	}
	pageVo.Data = vos
	return pageVo, nil
}

// 根据班级创建
func (t tagService) PostCreateStuByClass(arg domain.TagStuCreateByClassArg, user *comm.SessionUSER) error {
	if !tools.IsAdmin(user) {
		zap.S().Error("标签根据班级创建，当前用户不是管理员", user)
		return errors.New("非法操作")
	}
	userService := *t.userService

	stuList := userService.ListStuBySchoolYearClass(arg.SchoolYear, arg.Class)
	if stuList == nil {
		return errors.New("查询不到学生")
	}

	var tagStu []model.TagStu
	err := t.mysql.Where("tag_id = ?", arg.TagId).Find(&tagStu)
	if err != nil {
		return ere.ErrorCommFindError
	}

	var stus = t.deduplication(stuList, &tagStu)

	now := time.Now()
	var newTagStu []model.TagStu
	for _, stu := range *stus {
		var t model.TagStu
		t.Id = tools.SnowFlake.Generate().Int64()
		t.Ctime = now
		t.Uid = stu.Id
		t.TagId = tools.StringToInt64(arg.TagId)
		t.Creator = user.Id
		t.Name = stu.Name
		t.SchoolYear = stu.SchoolYear
		t.Class = stu.Class
		t.Account = stu.Account
		newTagStu = append(newTagStu, t)
	}
	if _, err = t.mysql.Insert(newTagStu); err != nil {
		return ere.ErrorCommSaveError
	}
	return nil
}

// 删除
func (t tagService) PostRemove(tag string, user *comm.SessionUSER) error {
	if !tools.IsAdmin(user) {
		zap.S().Error("标签删除，当前用户不是管理员", user)
		return errors.New("非法操作")
	}
	var tg model.Tag
	has, err := t.mysql.Where("id = ?", tag).And("enable = ?", comm.Enable).Get(&tg)
	if err != nil {
		zap.S().Error(err)
		return ere.ErrorCommDeleteError
	}
	if !has {
		return errors.New("标签不存在")
	}
	tg.Enable = comm.Disable
	tg.Utime = time.Now()
	if _, err := t.mysql.Id(tg.Id).Cols("enable", "utime").Update(tg); err != nil {
		zap.S().Error(err)
		return ere.ErrorCommDeleteError
	}
	return nil
}

// 创建
func (t tagService) PostCreate(tag string, sessionUser *comm.SessionUSER) (*domain.TagVo, error) {
	if !tools.IsAdmin(sessionUser) {
		zap.S().Error("标签创建，当前用户不是管理员", sessionUser)
		return nil, errors.New("非法操作")
	}
	var tg model.Tag
	tg.Name = tag
	tg.Id = tools.SnowFlake.Generate().Int64()
	tg.TotalNum = 0
	tg.Ctime = time.Now()
	tg.Creator = sessionUser.Id
	tg.Utime = time.Now()
	tg.Enable = comm.Enable
	if _, err := t.mysql.Insert(tg); err != nil {
		zap.S().Error(err)
		return nil, ere.ErrorCommSaveError
	}
	return domain.CvTagToTagVo(&tg), nil
}

// 查询所有
func (t tagService) PostList() (*[]domain.TagVo, error) {
	var tags []model.Tag
	err := t.mysql.Where("enable = ?", comm.Enable).Find(&tags)
	if err != nil {
		zap.S().Error(err)
		return nil, ere.ErrorCommFindError
	}
	tagVo := make([]domain.TagVo, 0)

	for _, t := range tags {
		tagVo = append(tagVo, *domain.CvTagToTagVo(&t))
	}
	return &tagVo, nil
}

// 去重
func (t tagService) deduplication(source *[]model.User, compare *[]model.TagStu) *[]model.User {
	var stus []model.User
	if compare == nil {
		stus = *source
	} else {
		for _, a := range *source {
			exist := false
			for _, b := range *compare {
				if a.Id == b.Uid {
					exist = true
				}
			}
			if !exist {
				stus = append(stus, a)
			}
		}
	}
	return &stus
}

type TagService interface {
	// 查询所有
	PostList() (*[]domain.TagVo, error)

	// 创建
	PostCreate(tag string, sessionUser *comm.SessionUSER) (*domain.TagVo, error)

	// 删除
	PostRemove(tag string, user *comm.SessionUSER) error

	// 根据班级创建
	PostCreateStuByClass(arg domain.TagStuCreateByClassArg, user *comm.SessionUSER) error

	// 分页
	PostPage(d *domain.TagStuPageArg, user *comm.SessionUSER) (interface{}, error)

	// 移除学生
	PostRemoveTagStu(d *domain.TagStuRemoveArg, user *comm.SessionUSER) error

	// 批量创建
	PostCreateTagStuList(i *domain.TagStuCreateListArg, user *comm.SessionUSER) error

	// 到处模板
	PostExportTemplate() *excelize.File

	// 导入用户
	PostImport(reader multipart.File, tagId string, user *comm.SessionUSER) error
}

func NewTagService(mysql *xorm.Engine, userService *UserService) TagService {
	return &tagService{
		mysql:       mysql,
		userService: userService,
	}
}
