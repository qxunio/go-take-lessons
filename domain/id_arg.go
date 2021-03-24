package domain

type IdArg struct {
	Id string `json:"id" valid:"required,isNotBlank~ID不能为空"`
}
