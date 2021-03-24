package domain

import (
	"go-take-lessons/domain/comm"
	"go-take-lessons/model"
	"go-take-lessons/tools"
	"strconv"
	"time"
)

type TagVo struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type TagStuVo struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Account    string `json:"account"`
	Class      string `json:"class"`
	Ctime      string `json:"ctime"`
	SchoolYear string `json:"schoolYear"`
}

type TagStuCreateByClassArg struct {
	SchoolYear string `json:"schoolYear" valid:"required,isNotBlank~入学年不能为空"`
	Class      string `json:"class" valid:"required,isNotBlank~班级不能为空"`
	TagId      string `json:"tagId" valid:"required,isNotBlank~标签ID不能为空"`
}

type TagStuRemoveArg struct {
	Uid   string `json:"uid" valid:"required,isNotBlank~学生不能为空"`
	TagId string `json:"tagId" valid:"required,isNotBlank~标签ID不能为空"`
}

type TagStuPageArg struct {
	UserPageArg
	TagId string `json:"tagId" valid:"required,isNotBlank~标签ID不能为空"`
}

type TagStuCreateListArg struct {
	StuList []TagStuCreateArg `json:"stuList" valid:"required"`
	TagId   string            `json:"tagId" valid:"required,isNotBlank~标签ID不能为空"`
}

type TagStuCreateArg struct {
	Uid        string `json:"id" valid:"required,isNotBlank~学生不能为空"`
	Name       string `json:"name" valid:"required,isNotBlank~参数缺省"`
	Class      string `json:"class" valid:"required,isNotBlank~参数缺省"`
	SchoolYear string `json:"schoolYear" valid:"required,isNotBlank~参数缺省"`
	Account    string `json:"account" valid:"required,isNotBlank~参数缺省"`
}

func CvTagToTagVo(tag *model.Tag) *TagVo {
	return &TagVo{
		Id:   strconv.FormatInt(tag.Id, 10),
		Name: tag.Name,
	}
}

func CvTagStuToTagStuVo(tagStu *model.TagStu) *TagStuVo {
	return &TagStuVo{
		Id:         strconv.FormatInt(tagStu.Uid, 10),
		Name:       tagStu.Name,
		Account:    tagStu.Account,
		Class:      strconv.FormatInt(int64(tagStu.Class), 10),
		Ctime:      tagStu.Ctime.Format(comm.TimeFormatTime),
		SchoolYear: strconv.Itoa(tagStu.SchoolYear),
	}
}

func CvTagStuCreateArgToTagStu(arg *TagStuCreateArg, tagId, creator int64) *model.TagStu {
	return &model.TagStu{
		Id:         tools.SnowFlake.Generate().Int64(),
		Uid:        tools.StringToInt64(arg.Uid),
		TagId:      tagId,
		Ctime:      time.Now(),
		Creator:    creator,
		Name:       arg.Name,
		Class:      tools.StringToInt(arg.Class),
		SchoolYear: tools.StringToInt(arg.SchoolYear),
		Account:    arg.Account,
	}
}
