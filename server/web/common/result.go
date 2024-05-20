package common

type Result struct {
	Code   int         `json:"code"`
	ErrMsg string      `json:"errmsg"`
	Data   interface{} `json:"data"`
}

func Error(code int, msg string) *Result {
	return &Result{
		Code:   code,
		ErrMsg: msg,
	}
}

func Success(code int, data interface{}, msg string) *Result {
	return &Result{
		Code:   code,
		Data:   data,
		ErrMsg: msg,
	}
}
