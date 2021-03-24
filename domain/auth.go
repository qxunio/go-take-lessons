package domain

type AuthLoginarg struct {
	A1 string `json:"a1" valid:"required,isNotBlank~账号不能为空"`
	A2 string `json:"a2" valid:"required,isNotBlank~密码不能为空"`
	A3 string `json:"a3"valid:"required,isNotBlank~请刷新重试"`
}

type AuthLoginVo struct {
	Uid       string `json:"uid"`
	Token     string `json:"token"`
	Username  string `json:"username"`
	Type      string `json:"type"`
	LoginTime string `json:"loginTime"`
	Device    string `json:"device"`
}
