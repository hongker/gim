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
		// decode request from context.
		err := DecodeRequestFromContext(ctx, req)
		// abort response if err is nil
		Abort(err)

		// process request and return response
		response, err := fn(NewValidatedContext(ctx), req)
		Abort(err)

		// output response
		Success(ctx, response)
	}
}

// RequestBodyFromContext returns the request body from the context.
func RequestBodyFromContext(ctx *gin.Context) (p []byte, err error) {
	return ctx.GetRawData()
}

// DecodeRequestFromContext
func DecodeRequestFromContext(ctx *gin.Context, container interface{}) error {
	body, err := RequestBodyFromContext(ctx)
	if err != nil {
		return err
	}
	return serializer().Decode(body, container)
}
