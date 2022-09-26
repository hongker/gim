package api

import (
	"context"
	"gim/internal/module/gateway/render"
	"github.com/gin-gonic/gin"
)

type Handler interface {
	Install(router *gin.Engine)
}

// Action returns the formatted handler use by generics params.
func Action[Request any, Response any](fn func(context.Context, *Request) (*Response, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := new(Request)
		err := render.SerializeRequestFromContext(ctx, req)
		render.Abort(err)

		response, err := fn(NewValidatedContext(ctx), req)
		render.Abort(err)

		render.Success(ctx, response)
	}
}
