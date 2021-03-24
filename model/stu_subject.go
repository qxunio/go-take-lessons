package model

import "time"

type StuSubject struct {
	Id         int64     `xorm:"pk BigInt" json:"id"`
	UserId     int64     `xorm:"BigInt" json:"user_id"`
	EventId    int64     `xorm:"BigInt" json:"event_id"`
	EventName  string    `xorm:"varchar(30)" json:"event_name"`
	CsId       int64     `xorm:"BigInt" json:"cs_id"`
	SchoolYear int       `xorm:"Int" json:"school_year"`
	Class      int       `xorm:"Int" json:"class"`
	Enable     uint8     `xorm:"Int" json:"enable"`
	Ctime      time.Time `xorm:"DateTime" json:"ctime"`
	Utime      time.Time `xorm:"DateTime" json:"utime"`
}

type StuSubjectDo struct {
	Id         int64     `json:"id"`
	UserId     int64     `json:"userId"`
	Name       string    `json:"name"`
	Account    string    `json:"account"`
	Class      int       `json:"class"`
	SchoolYear int       `json:"schoolYear"`
	Ctime      time.Time `json:"ctime"`
}

type StuSubjectAndBaseDo struct {
	Id          int64     `json:"id"`
	UserId      int64     `json:"userId"`
	CsId        int64     `json:"cs_id"`
	ClassName   string    `json:"class_name"`
	SubjectName string    `json:"subject_name"`
	Ctime       time.Time `json:"ctime"`
}
