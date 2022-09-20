package render

import (
	"encoding/json"
	"github.com/ebar-go/ego/errors"
	"github.com/ebar-go/ego/server/ws"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Success output success response.
func Success(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, Response{
		Code: 0,
		Msg:  "",
		Data: data,
	})
}

// Error output error response.
func Error(ctx *gin.Context, err error) {
	se := errors.Convert(err)
	ctx.JSON(http.StatusOK, Response{
		Code: se.Code(),
		Msg:  se.Message(),
	})
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

func SocketSuccess(ctx *ws.Context, data any) {
	response := Response{
		Code: 0,
		Msg:  "",
		Data: data,
	}
	b, _ := json.Marshal(response)
	ctx.Output(b)
}
func SocketAbort(err error) {}
