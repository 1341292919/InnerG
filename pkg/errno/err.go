package errno

import (
	"errors"
	"fmt"
)

type Err struct {
	ErrorCode int64
	ErrorMsg  string
}

// Error 是实现原生err的必要函数
func (e Err) Error() string {
	return fmt.Sprintf("[%d]:%s", e.ErrorCode, e.ErrorMsg)
}

func NewErr(code int64, msg string) Err {
	return Err{
		ErrorMsg:  msg,
		ErrorCode: code,
	}
}

func ConvertErr(err error) Err {
	if err == nil {
		return Success
	}
	errno := Err{}
	if errors.As(err, &errno) {
		return errno
	}
	s := InternalServiceError
	s.ErrorMsg = err.Error()
	return s
}

func (e Err) WithMessage(msg string) Err {
	e.ErrorMsg = msg
	return e
}
