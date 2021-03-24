package model

import "time"

type StuFocus struct {
	Id         int64     `xorm:"pk BigInt" json:"id"`
	UserId     int64     `xorm:"BigInt" json:"user_id"`
	EventId    int64     `xorm:"BigInt" json:"event_id"`
	EventName  string    `xorm:"BigInt" json:"event_name"`
	CsId       int64     `xorm:"BigInt" json:"cs_id"`
	SchoolYear int       `xorm:"Int" json:"school_year"`
	Enable     uint8     `xorm:"Int" json:"enable"`
	Ctime      time.Time `xorm:"DateTime" json:"ctime"`
	Utime      time.Time `xorm:"DateTime" json:"utime"`
}
