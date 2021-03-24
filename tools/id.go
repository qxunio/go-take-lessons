package tools

import (
	"go-take-lessons/third_party"
	"go.uber.org/zap"
)

var SnowFlake *third_party.Node

func InitSnowFlakeId() {
	var err error
	SnowFlake, err = third_party.NewNode(1)
	if err != nil {
		zap.S().Error("Init SnowFlake Fail ", err)
		return
	}
}
