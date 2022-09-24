package http

import (
	"gim/internal/module/gateway/render"
	"github.com/ebar-go/ego/errors"
	"github.com/gin-gonic/gin"
)

// recoverMiddleware returns a recover middleware.
func recoverMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				var err error
				switch r.(type) {
				case error:
					err = r.(error)
				default:
					err = errors.Unknown("system error")
				}
				render.Error(ctx, err)

			}
		}()
		ctx.Next()
	}
}
