package model

import "time"

type User struct {
	Id         int64     `xorm:"pk BigInt" json:"id"`
	Name       string    `xorm:"varchar(32)" json:"name"`
	Salt       string    `xorm:"varchar(12)" json:"salt"`
	Account    string    `xorm:"varchar(15)" json:"account"`
	Password   string    `xorm:"varchar(32)" json:"password"`
	SchoolYear int       `xorm:"Int" json:"school_year"`
	Ctime      time.Time `xorm:"DateTime" json:"ctime"`
	Utime      time.Time `xorm:"DateTime" json:"utime"`
	Enable     uint8     `xorm:"Int" json:"enable"`
	Type       uint8     `xorm:"Int" json:"super"`
	Class      int       `xorm:"Int" json:"class"`
	Creator    int64     `xorm:"BigInt" json:"creator"`
}

type UserCountDo struct {
	Type  int `json:"type"`
	Count int `json:"count"`
}

type StuDo struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	Account    string `json:"account"`
	Class      int    `json:"type"`
	SchoolYear int    `json:"schoolYear"`
}
