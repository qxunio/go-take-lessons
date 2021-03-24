package comm

import (
	"go-take-lessons/model"
	"time"
)

// 当前登录用户信息
type SessionUSER struct {
	Id         int64     `json:"id"`
	Name       string    `json:"name"`
	Account    string    `json:"account"`
	Super      byte      `json:"super"`
	LoginTime  time.Time `json:"login_time"`
	SchoolYear int       `json:"school_year"`
	UserType   byte      `json:"user_type"`
	Class      int       `json:"class"`
}

func ConversionSessionUSER(user model.User) *SessionUSER {
	return &SessionUSER{
		Id:         user.Id,
		Name:       user.Name,
		Account:    user.Account,
		UserType:   user.Type,
		LoginTime:  time.Now(),
		SchoolYear: user.SchoolYear,
		Class:      user.Class,
	}
}
