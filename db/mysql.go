package db

import (
	"bytes"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"go-take-lessons/configs"
)

func InitMysql(conf *configs.MysqlConfig) *xorm.Engine {
	var dataSource bytes.Buffer
	dataSource.WriteString(conf.User)
	dataSource.WriteString(":")
	dataSource.WriteString(conf.Password)
	dataSource.WriteString("@(")
	dataSource.WriteString(conf.Host)
	dataSource.WriteString(")/")
	dataSource.WriteString(conf.DbName)
	dataSource.WriteString("?charset=utf8")
	engine, err := xorm.NewEngine("mysql", dataSource.String())
	if err != nil {
		panic("Mysql 连接失败:" + err.Error())
	}
	engine.ShowSQL(true)
	engine.SetMaxIdleConns(conf.MaxIdle)
	engine.SetMaxOpenConns(conf.MaxOpen)
	// 开启缓存会导致查询到NULL值错误
	//cacher := xorm.NewLRUCacher(xorm.NewMemoryStore(), 1000)
	//engine.SetDefaultCacher(cacher)
	return engine
}
