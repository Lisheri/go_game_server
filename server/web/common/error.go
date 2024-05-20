package common

type MyError struct {
	code int
	msg  string
}

// 这里直接实现一个Error方法
func (m *MyError) Error() string {
	return m.msg
}

func (m *MyError) Code() int {
	return m.code
}

func New(code int, msg string) error {
	return &MyError{
		code: code,
		msg:  msg,
	}
}
