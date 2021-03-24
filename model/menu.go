package model

import "time"

type Menu struct {
	Id       int64     `xorm:"pk BigInt" json:"id"`
	Name     string    `xorm:"varchar(32)" json:"name"`
	Router   string    `xorm:"varchar(32)" json:"router"`
	Priority uint8     `xorm:"Int" json:"priority"`
	ParentId int64     `xorm:"pk BigInt" json:"parentId"`
	Enable   uint8     `xorm:"Int" json:"enable"`
	Ctime    time.Time `xorm:"DateTime" json:"ctime"`
}
