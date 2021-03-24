package configs

import "github.com/BurntSushi/toml"

var Conf QxunConfig

// configs
type QxunConfig struct {
	App   AppConfig
	Mysql MysqlConfig
	Redis RedisConfig
}

type AppConfig struct {
	Port      string
	LoginTime int64
}

type MysqlConfig struct {
	Host     string
	User     string
	Password string
	DbName   string
	MaxIdle  int
	MaxOpen  int
}

type RedisConfig struct {
	Host     string
	DbNumber int
}

// 读取配置文件
func InitConfigs() {
	if _, err := toml.DecodeFile("configs/iris.toml", &Conf); err != nil {
		panic("读取配置文件错误" + err.Error())
	}
}
