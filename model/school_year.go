package model

import "time"

type SchoolYear struct {
	Id         int64     `xorm:"pk BigInt" json:"id"`
	SchoolYear int       `xorm:"Int" json:"school_year"`
	Enable     uint8     `xorm:"Int" json:"enable"`
	Ctime      time.Time `xorm:"DateTime" json:"ctime"`
	Utime      time.Time `xorm:"DateTime" json:"utime"`
}
