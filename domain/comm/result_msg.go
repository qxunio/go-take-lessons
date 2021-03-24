package comm

type ResultMsg struct {
	Code   string      `json:"code"`
	Msg    string      `json:"msg"`
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

// 错误响应 通用Result
func ErrorUnKnownResponse() *ResultMsg {
	data := &ResultMsg{
		Code:   "500",
		Msg:    "UnKnow Error",
		Data:   nil,
		Status: "fail",
	}
	return data
}

// 错误响应 带提示信息
func ErrorResponseMsg(msg string) *ResultMsg {
	data := &ResultMsg{
		Code:   "500",
		Msg:    msg,
		Data:   nil,
		Status: "fail",
	}
	return data
}

// 错误响应 带提示信息
func ErrorResponseCodeMsg(code, msg string) *ResultMsg {
	data := &ResultMsg{
		Code:   code,
		Msg:    msg,
		Data:   nil,
		Status: "fail",
	}
	return data
}

// 错误响应 带数据
func ErrorResponseData(d interface{}) *ResultMsg {
	data := &ResultMsg{
		Code:   "500",
		Msg:    "Error",
		Data:   d,
		Status: "fail",
	}
	return data
}

// 成功响应  通用Result
func SuccessResponse() *ResultMsg {
	data := &ResultMsg{
		Code:   "200",
		Msg:    "OK",
		Data:   nil,
		Status: "success",
	}
	return data
}

// 成功响应  带数据
func SuccessResponseData(d interface{}) *ResultMsg {
	data := &ResultMsg{
		Code:   "200",
		Msg:    "OK",
		Data:   d,
		Status: "success",
	}
	return data
}

// 成功响应  带消息
func SuccessResponseMsg(msg string) *ResultMsg {
	data := &ResultMsg{
		Code:   "200",
		Msg:    msg,
		Data:   nil,
		Status: "success",
	}
	return data
}
