package model

import "time"

type UserRole struct {
	Id     int64     `xorm:"pk BigInt" json:"id"`
	UserId int64     `xorm:"BigInt" json:"userId"`
	RoleId int64     `xorm:"BigInt" json:"roleId"`
	Enable uint8     `xorm:"Int" json:"enable"`
	Ctime  time.Time `xorm:"DateTime" json:"ctime"`
	Utime  time.Time `xorm:"DateTime" json:"utime"`
}
