package model

import "time"

type ConfigurationSubject struct {
	Id             int64     `xorm:"pk BigInt" json:"id"`
	EventId        int64     `xorm:"BigInt" json:"event_id"`
	SubjectId      int64     `xorm:"BigInt" json:"subject_id"`
	SubjectName    string    `xorm:"varchar(32)" json:"subject_name"`
	Introduction   string    `xorm:"varchar(912)" json:"introduction"`
	ClassName      string    `xorm:"varchar(32)" json:"class_name"`
	Num            uint8     `xorm:"Int" json:"num"`
	SelectedPlaces uint8     `xorm:"Int" json:"selected_places"`
	TeachAddress   string    `xorm:"varchar(60)" json:"teach_address"`
	TeachTime      string    `xorm:"varchar(60)" json:"teach_time"`
	Enable         uint8     `xorm:"Int" json:"enable"`
	Ctime          time.Time `xorm:"DateTime" json:"ctime"`
	Utime          time.Time `xorm:"DateTime" json:"utime"`
}

type ConfigurationTeacher struct {
	Id             int64     `xorm:"pk BigInt" json:"id"`
	EventId        int64     `xorm:"BigInt" json:"event_id"`
	CsId           int64     `xorm:"BigInt" json:"cs_id"`
	TeacherId      int64     `xorm:"BigInt" json:"teacher_id"`
	TeacherName    string    `xorm:"varchar(32)" json:"teacher_name"`
	TeacherAccount string    `xorm:"varchar(32)" json:"teacher_account"`
	Enable         uint8     `xorm:"Int" json:"enable"`
	Ctime          time.Time `xorm:"DateTime" json:"ctime"`
	Utime          time.Time `xorm:"DateTime" json:"utime"`
}

type ConfigurationSimpleDo struct {
	Id           int64  `json:"id"`
	SubjectId    int64  `json:"subject_id"`
	SubjectName  string `json:"subject_name"`
	ClassName    string `json:"class_name"`
	Introduction string `json:"introduction"`
	Num          int    `json:"num"`
	TeachAddress string `json:"teach_address"`
	TeachTime    string `json:"teach_time"`
	Teacher      string `json:"teacher"`
}
