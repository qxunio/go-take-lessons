package model

import "time"

type Notice struct {
	Id         int64     `xorm:"pk BigInt" json:"id"`
	Content    string    `xorm:"varchar(330)" json:"content"`
	Title      string    `xorm:"varchar(60)" json:"title"`
	Type       uint8     `xorm:"Int" json:"type"`
	ExpireDate time.Time `xorm:"DateTime" json:"expire_date"`
	Status     uint8     `xorm:"Int" json:"status"`
	Enable     uint8     `xorm:"Int" json:"enable"`
	Ctime      time.Time `xorm:"DateTime" json:"ctime"`
	Utime      time.Time `xorm:"DateTime" json:"utime"`
	Creator    int64     `xorm:"BigInt" json:"creator"`
}
