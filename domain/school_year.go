package domain

import (
	"go-take-lessons/model"
	"strconv"
)

type SchoolYearVo struct {
	Id         string `json:"id"`
	SchoolYear string `json:"schoolYear"`
}

func CvSchoolYearToVo(schoolYear *model.SchoolYear) *SchoolYearVo {
	return &SchoolYearVo{
		Id:         strconv.FormatInt(schoolYear.Id, 10),
		SchoolYear: strconv.Itoa(schoolYear.SchoolYear),
	}
}
