package main

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/go-xorm/xorm"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/mvc"
	"go-take-lessons/configs"
	"go-take-lessons/controller"
	"go-take-lessons/db"
	"go-take-lessons/domain/comm"
	"go-take-lessons/service"
	"go-take-lessons/tools"
	"go.uber.org/zap"
	"strings"
)

func main() {
	// 初始化日志
	logger, err := zap.NewDevelopment()
	if err != nil {
		println(err.Error())
	}
	zap.ReplaceGlobals(logger)

	app := iris.New()
	app.Logger().SetLevel("debug")
	// 将请求记录到终端
	app.UseGlobal(beforeHandler)
	app.DoneGlobal(afterHandler)

	registerValidator()
	resetHandError(app)

	configs.InitConfigs()
	mysql := initBaseServer(configs.Conf)
	registerRouter(app, mysql)
	_ = app.Run(iris.Addr(":"+configs.Conf.App.Port), iris.WithoutServerError(iris.ErrServerClosed), iris.WithOptimizations)
}

// 注册数据校验
func registerValidator() {
	govalidator.CustomTypeTagMap.Set("isNotBlank", func(i interface{}, context interface{}) bool {
		return govalidator.IsNotNull(i.(string))
	})
}

// 注册路由
func registerRouter(app *iris.Application, mysql *xorm.Engine) {
	// service
	roleService := service.NewRoleService(mysql)
	userService := service.NewUserService(mysql, &roleService)
	configurationService := service.NewConfigurationService(mysql, &userService)
	eventService := service.NewEventService(mysql)

	// APP
	appMvc := mvc.New(app.Party("/app"))
	appMvc.Register(service.NewAppService(mysql))
	appMvc.Handle(new(controller.AppController))

	// 用户
	userMvc := mvc.New(app.Party("/user"))
	userMvc.Register(userService)
	userMvc.Handle(new(controller.UserController))

	// 认证
	authMvc := mvc.New(app.Party("/auth"))
	authMvc.Register(service.NewAuthService(mysql))
	authMvc.Handle(new(controller.AuthController))

	// 路由
	menuMvc := mvc.New(app.Party("/menu"))
	menuMvc.Register(service.NewMenuService(mysql))
	menuMvc.Handle(new(controller.MenuController))

	// 学科
	subjectMvc := mvc.New(app.Party("/subject"))
	subjectMvc.Register(service.NewSubjectService(mysql))
	subjectMvc.Handle(new(controller.SubjectController))

	// 选课
	eventMvc := mvc.New(app.Party("/event"))
	eventMvc.Register(eventService)
	eventMvc.Handle(new(controller.EventController))

	// 配置
	configurationMvc := mvc.New(app.Party("/configuration"))
	configurationMvc.Register(configurationService)
	configurationMvc.Handle(new(controller.ConfigurationController))

	// 入学年
	schoolYearMvc := mvc.New(app.Party("/sy"))
	schoolYearMvc.Register(service.NewSchoolYearService(mysql))
	schoolYearMvc.Handle(new(controller.SchoolYearController))

	// 关注 （学生）
	stuFocusMvc := mvc.New(app.Party("/focus"))
	stuFocusMvc.Register(service.NewStuFocusService(mysql, &configurationService))
	stuFocusMvc.Handle(new(controller.StuFocusController))

	// 选课
	stuChooseSubjectMvc := mvc.New(app.Party("/choose-subject"))
	stuChooseSubjectMvc.Register(service.NewStuChooseSubjectService(mysql, &eventService, &configurationService))
	stuChooseSubjectMvc.Handle(new(controller.StuChooseSubjectController))

	// 标签
	tagMvc := mvc.New(app.Party("/tag"))
	tagMvc.Register(service.NewTagService(mysql, &userService))
	tagMvc.Handle(new(controller.TagController))
}

// 初始化一些基础服务
func initBaseServer(conf configs.QxunConfig) (engine *xorm.Engine) {
	engine = db.InitMysql(&conf.Mysql)
	db.InitRedis(&conf.Redis)
	tools.InitSnowFlakeId()
	return engine
}

// 前置处理器
func beforeHandler(ctx iris.Context) {
	if strings.Compare(ctx.Path(), "/auth/login") != 0 && strings.Compare(ctx.Path(), "/app") != 0 {
		header := ctx.GetHeader("Authorization")
		if len(header) == 0 {
			_, _ = ctx.JSON(comm.ErrorResponseCodeMsg("555", "未登录"))
			return
		}
		userId, err := tools.ValidateToken(header)
		if err != nil {
			_, _ = ctx.JSON(comm.ErrorResponseCodeMsg("555", "登录过期"))
			return
		}
		sessionUserJson, err := db.RedisClient.Get(comm.RedisAuthTokenKey + userId).Result()
		if err != nil {
			sessionUserJson, err = db.RedisClient.Get(comm.RedisPreAuthTokenKey + userId).Result()
			if err != nil {
				zap.S().Error("前处理器,登录拦截: json序列化失败")
				_, _ = ctx.JSON(comm.ErrorResponseCodeMsg("555", "登录过期"))
			}
			return
		}
		sessionUser := new(comm.SessionUSER)
		if err = json.Unmarshal([]byte(sessionUserJson), sessionUser); err != nil {
			zap.S().Error("前处理器,登录拦截: json序列化失败")
			_, _ = ctx.JSON(comm.ErrorResponseCodeMsg("555", "获取登录用户失败"))
			return
		}
		ctx.Values().Set(comm.ContextSessionUserKey, sessionUser)
	}
	ctx.Next()
}

// 后置处理器
func afterHandler(ctx iris.Context) {
	ctx.Next()
}

// 错误处理
func resetHandError(app *iris.Application) {
	app.OnErrorCode(iris.StatusNotFound, func(context context.Context) {
		_, _ = context.JSON(comm.ErrorResponseMsg("资源不存在"))
		return
	})
	app.OnErrorCode(iris.StatusInternalServerError, func(context context.Context) {
		_, _ = context.JSON(comm.ErrorResponseMsg("内部错误"))
		return
	})
}
