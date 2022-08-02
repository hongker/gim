package errors

import (
	"fmt"
)

const(
	CodeFailure = 1001
	CodeInvalidParameter = 1002
	CodeDataNotFound = 1003
	CodeForbidden = 1004

)

type Error struct {
	code int
	msg string
}

func (e Error) Error() string {
	return fmt.Sprintf("errors: code=%d msg=%s", e.code, e.msg)
}

func (e Error) Code() int {
	return e.code
}
func (e Error) Message() string {
	return e.msg
}

func New(code int, msg string) *Error {
	return &Error{code: code, msg: msg}
}

func InvalidParameter(msg string) *Error {
	return New(CodeInvalidParameter, msg)
}
func DataNotFound(msg string) *Error {
	return New(CodeDataNotFound, msg)
}

func Failure(msg string) *Error {
	return New(CodeFailure, msg)
}

func WithMessage(err error, msg string) *Error {
	e := Convert(err)
	if e == nil {
		return nil
	}
	e.msg = fmt.Sprintf("%s: %s", msg, e.msg)
	return e
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