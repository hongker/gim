package api

import "gim/pkg/errors"

const (
	CodeSuccess = 0
)

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func NewFailureResponse(err error) *Response {
	e := errors.Convert(err)
	return &Response{Code: e.Code(), Msg: e.Message()}
}

func NewSuccessResponse(data any) *Response {
	return &Response{Code: CodeSuccess, Data: data}
}
