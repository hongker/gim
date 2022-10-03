package api

import (
	"context"
	"gim/internal/domain/stateful"
	"gim/internal/domain/types"
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
				Error(ctx, err)

			}
		}()
		ctx.Next()
	}
}

const (
	tokenParam = "token"
	userParam  = "user"
)

func checkToken(auth types.Authenticator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Query(tokenParam)
		if len(token) == 0 {
			Abort(errors.Unauthorized("invalid token"))
		}

		uid, err := auth.Authenticate(ctx, token)
		Abort(errors.WithMessage(err, "authenticate"))

		ctx.Set(userParam, uid)
		ctx.Next()
	}
}

func NewValidatedContext(ctx *gin.Context) context.Context {
	return stateful.NewUserContext(ctx, ctx.GetString(userParam))
}
