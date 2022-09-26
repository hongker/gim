package api

import (
	"context"
	"github.com/gin-gonic/gin"
)

type Handler interface {
	Install(router *gin.Engine)
}

// generic returns the formatted handler use by generics params.
func generic[Request any, Response any](fn func(context.Context, *Request) (*Response, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := new(Request)
		err := SerializeRequestFromContext(ctx, req)
		Abort(err)

		response, err := fn(NewValidatedContext(ctx), req)
		Abort(err)

		Success(ctx, response)
	}
}

// RequestBodyFromContext returns the request body from the context.
func RequestBodyFromContext(ctx *gin.Context) (p []byte, err error) {
	return ctx.GetRawData()
}

// SerializeRequestFromContext
func SerializeRequestFromContext(ctx *gin.Context, container interface{}) error {
	body, err := RequestBodyFromContext(ctx)
	if err != nil {
		return err
	}
	return serializer().Decode(body, container)
}
