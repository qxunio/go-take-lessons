package model

import "time"

type Event struct {
	Id         int64     `xorm:"pk BigInt" json:"id"`
	Name       string    `xorm:"varchar(32)" json:"name"`
	CanUpdate  int       `xorm:"Int" json:"can_update"`
	Num        uint8     `xorm:"Int" json:"num"`
	SchoolYear string    `xorm:"varchar(32)" json:"school_year"`
	TagIds     string    `xorm:"varchar(32)" json:"tag_ids"`
	Stime      time.Time `xorm:"DateTime" json:"stime"`
	Etime      time.Time `xorm:"DateTime" json:"etime"`
	Status     uint8     `xorm:"Int" json:"status"`
	Enable     uint8     `xorm:"Int" json:"enable"`
	Ctime      time.Time `xorm:"DateTime" json:"ctime"`
	Utime      time.Time `xorm:"DateTime" json:"utime"`
	Creator    int64     `xorm:"BigInt" json:"creator"`
}
