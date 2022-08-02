package api

import "gim/pkg/errors"

const (
	CodeSuccess = 0
)

type Response struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
	Data interface{} `json:"data"`
}

func FailureResponse(err error) *Response  {
	e := errors.Convert(err)
	return &Response{Code: e.Code(), Msg: e.Message()}
}

func SuccessResponse(data interface{}) *Response {
	return &Response{Code: CodeSuccess, Data: data}
}
