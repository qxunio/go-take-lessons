package tools

import (
	"go-take-lessons/domain/comm"
	"regexp"
	"strings"
)

// 是否是数字的字符
func IsNumber(str string) bool {
	pattern := "\\d+"
	matched, err := regexp.MatchString(pattern, str)
	if err != nil {
		return false
	}
	return matched
}

// 是否是空白字符
func IsBlank(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}

// 是否是管理员
func IsAdmin(sessionUser *comm.SessionUSER) bool {
	return sessionUser.UserType == comm.Admin
}
