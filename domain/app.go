package domain

type AppCodeVo struct {
	Uid string `json:"uid"`
	Pk  string `json:"pk"`
}
type AdminIndexVo struct {
	Account    int            `json:"account"`
	Student    int            `json:"student"`
	Teacher    int            `json:"teacher"`
	EventNum   int            `json:"eventNum"`
	SubjectNum int            `json:"subjectNum"`
	TagNum     int            `json:"tagNum"`
	HotSubject []HotSubjectVo `json:"hotSubject"`
}

type HotSubjectVo struct {
	SubjectName string `json:"subjectName"`
	Num         string `json:"num"`
}

func CvDefaultAdminIndexVo() *AdminIndexVo {
	return &AdminIndexVo{
		Account:    0,
		Student:    0,
		Teacher:    0,
		EventNum:   0,
		SubjectNum: 0,
		TagNum:     0,
		HotSubject: nil,
	}
}
