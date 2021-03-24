package service

import (
	"github.com/go-xorm/xorm"
	uuid "github.com/satori/go.uuid"
	"go-take-lessons/db"
	"go-take-lessons/domain"
	"go-take-lessons/domain/comm"
	"go-take-lessons/model"
	"go-take-lessons/tools"
	"go.uber.org/zap"
	"strings"
	"time"
)

type appService struct {
	mysql *xorm.Engine
}

// 管理员Index数据
func (app appService) PostAdminIndex() *domain.AdminIndexVo {
	selectUser := "SELECT type AS type,count( 1 ) AS count FROM `user` WHERE `enable` = ?  AND type != 1 GROUP BY type"
	var userCountDo []model.UserCountDo
	if err := app.mysql.SQL(selectUser, comm.Enable).Find(&userCountDo); err != nil {
		zap.S().Error("管理员查询用户统计错误", err)
		return domain.CvDefaultAdminIndexVo()
	}

	var vo domain.AdminIndexVo
	if len(userCountDo) == 1 {
		count := userCountDo[0].Count
		if userCountDo[0].Type == 2 {
			vo.Teacher = count
			vo.Student = 0
		} else {
			vo.Teacher = 0
			vo.Student = count
		}
	} else {
		for _, dto := range userCountDo {
			if dto.Type == 2 {
				vo.Teacher = dto.Count
			} else {
				vo.Student = dto.Count
			}
		}
	}
	vo.Account = vo.Student + vo.Teacher

	selectSubject := "SELECT count(1) FROM subject WHERE `enable` = ?"
	subjectCount := 0
	if _, err := app.mysql.SQL(selectSubject, comm.Enable).Get(&subjectCount); err != nil {
		zap.S().Error("管理员查询学科库错误", err)
	}
	vo.SubjectNum = subjectCount

	selectTag := "SELECT count(1) FROM tag WHERE `enable` = ?"
	tagCount := 0
	if _, err := app.mysql.SQL(selectTag, comm.Enable).Get(&tagCount); err != nil {
		zap.S().Error("管理员查询学科库错误", err)
	}
	vo.TagNum = tagCount

	var hotSubject []domain.HotSubjectVo
	selectHotSubject := "SELECT cs.subject_name,tt.num FROM( SELECT cs_id, COUNT( 1 ) AS num FROM stu_subject WHERE `enable` = ? GROUP BY cs_id ) tt LEFT JOIN configuration_subject cs ON cs.id = tt.cs_id ORDER BY tt.num DESC LIMIT 0 ,15"
	if err := app.mysql.SQL(selectHotSubject, comm.Enable).Find(&hotSubject); err != nil {
		zap.S().Error("管理员查询热门学科错误", err)
	}
	vo.HotSubject = hotSubject

	selectEventNum := "SELECT count(1) FROM event WHERE `enable` = ?"
	eventCount := 0
	if _, err := app.mysql.SQL(selectEventNum, comm.Enable).Get(&eventCount); err != nil {
		zap.S().Error("管理员查询学科库错误", err)
	}
	vo.EventNum = eventCount
	return &vo
}

type AppService interface {
	// Post 获取密钥
	Post() domain.AppCodeVo

	// 管理员Index数据
	PostAdminIndex() *domain.AdminIndexVo
}

// Post 获取密钥
func (app appService) Post() domain.AppCodeVo {
	privateKey, publicKey, _ := tools.GenRSAEncrypt()
	public := strings.Replace(strings.Split(string(publicKey), "-----")[2], "\n", "", -1)
	UUID := uuid.NewV4().String()
	db.RedisClient.Set(comm.RedisAuthEncryptionKey+UUID, privateKey, time.Second*600)
	var appCode domain.AppCodeVo
	appCode.Uid = UUID
	appCode.Pk = public
	return appCode
}

func NewAppService(mysql *xorm.Engine) AppService {
	return &appService{mysql: mysql}
}
