package api

import (
	"gim/api"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Success output success response.
func Success(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, api.NewSuccessResponse(data))
}

// Error output error response.
func Error(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusOK, api.NewFailureResponse(err))
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
