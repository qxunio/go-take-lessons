package model

import "time"

type Subject struct {
	Id           int64     `xorm:"pk BigInt" json:"id"`
	Name         string    `xorm:"varchar(32)" json:"name"`
	Introduction string    `xorm:"varchar(912)" json:"introduction"`
	Enable       uint8     `xorm:"Int" json:"enable"`
	Ctime        time.Time `xorm:"DateTime" json:"ctime"`
	Utime        time.Time `xorm:"DateTime" json:"utime"`
	Creator      int64     `xorm:"BigInt" json:"creator"`
}
