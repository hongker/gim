package render

import "github.com/gin-gonic/gin"

func Success(ctx *gin.Context, data interface{}) {

}

func Error(err error) {

}

// Abort it will abort when err is not nil,
func Abort(ctx *gin.Context, err error) {
	if err == nil {
		return
	}
	Panic(err)
}

func Panic(err error) {

}
