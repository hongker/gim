package render

import (
	"github.com/ebar-go/ego/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Success output success response.
func Success(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, SuccessResponse(data))
}

// Error output error response.
func Error(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusOK, ErrorResponse(err))
	ctx.Abort()
}

// Abort it will abort when err is not nil.
func Abort(err error) {
	if err == nil {
		return
	}
	abortPanic(err)
}

func abortPanic(err error) {
	panic(err)
}

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func ErrorResponse(err error) Response {
	se := errors.Convert(err)
	return Response{
		Code: se.Code(),
		Msg:  se.Message(),
		Data: nil,
	}
}

func SuccessResponse(data any) Response {
	return Response{
		Code: 0,
		Msg:  "",
		Data: data,
	}
}
