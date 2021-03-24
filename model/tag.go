package model

import "time"

type Tag struct {
	Id       int64     `xorm:"pk BigInt" json:"id"`
	Name     string    `xorm:"varchar(32)" json:"name"`
	TotalNum int64     `xorm:"Int" json:"total_num"`
	Ctime    time.Time `xorm:"DateTime" json:"ctime"`
	Utime    time.Time `xorm:"DateTime" json:"utime"`
	Enable   uint8     `xorm:"Int" json:"enable"`
	Creator  int64     `xorm:"BigInt" json:"creator"`
}

type TagStu struct {
	Id         int64     `xorm:"pk BigInt" json:"id"`
	Uid        int64     `xorm:"BigInt" json:"uid"`
	TagId      int64     `xorm:"BigInt" json:"tag_id"`
	Ctime      time.Time `xorm:"DateTime" json:"ctime"`
	Creator    int64     `xorm:"BigInt" json:"creator"`
	Name       string    `xorm:"varchar(32)" json:"name"`
	Class      int       `xorm:"Int" json:"class"`
	SchoolYear int       `xorm:"Int" json:"school_year"`
	Account    string    `xorm:"varchar(32)" json:"account"`
}
