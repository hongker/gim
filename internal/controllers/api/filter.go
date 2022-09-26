package api

import (
	"context"
	"gim/internal/module/gateway/domain/types"
	"gim/internal/module/gateway/domain/types/auth"
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

const (
	TokenParam       = "token"
	CurrentUserParam = "currentUser"
)

func checkToken(auth types.Authenticator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Query(TokenParam)
		if len(token) == 0 {
			render.Abort(errors.Unauthorized("invalid token"))
		}

		uid, err := auth.Authenticate(ctx, token)
		render.Abort(errors.WithMessage(err, "authenticate"))

		ctx.Set(CurrentUserParam, uid)
		ctx.Next()
	}
}

func NewValidatedContext(ctx *gin.Context) context.Context {
	return auth.NewUserContext(ctx, ctx.GetString(CurrentUserParam))
}
