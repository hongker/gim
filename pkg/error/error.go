package error

import (
	"fmt"
)

const(
	CodeFailure = 1001
	CodeInvalidParameter = 1002

)

type Error struct {
	code int
	msg string
}

func (e Error) Error() string {
	return fmt.Sprintf("error: code=%d msg=%s", e.code, e.msg)
}

func New(code int, msg string) *Error {
	return &Error{code: code, msg: msg}
}

func InvalidParameter(msg string) *Error {
	return New(CodeInvalidParameter, msg)
}

func Failure(msg string) *Error {
	return New(CodeFailure, msg)
}


func Convert(err error) *Error {
	if err == nil {
		return nil
	}

	if e, ok := err.(*Error); ok {
		return e
	}
	return Failure(err.Error())
}