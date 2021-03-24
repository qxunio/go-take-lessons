package ere

import "errors"

const StrCommArgIsBlank = "参数为空"

const StrLoginFail = "登录失败！请稍后重试"

const StrParseArgFail = "解析参数失败"

var ErrorCommUnknown = errors.New("未知错误")

var ErrorCommFindError = errors.New("查询错误")

var ErrorCommNotFond = errors.New("未查询到")

var ErrorCommDeleteError = errors.New("删除错误")

var ErrorCommSaveError = errors.New("保存错误")

var ErrorCommUpdateError = errors.New("更新失败")

var ErrorCommConvertError = errors.New("转换错误")

var ErrorCommArgIsBlankError = errors.New(StrCommArgIsBlank)
