package db

import (
	"fmt"
	"github.com/go-redis/redis"
	"go-take-lessons/configs"
)

var RedisClient *redis.Client

func InitRedis(conf *configs.RedisConfig) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: conf.Host,
		DB:   conf.DbNumber,
	})
	pong, err := RedisClient.Ping().Result()
	if err != nil {
		panic("初始化Redis失败" + err.Error())
		return
	}
	fmt.Printf("Redis Response %s\n", pong)
}
