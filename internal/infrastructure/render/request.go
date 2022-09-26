package render

import (
	"github.com/gin-gonic/gin"
)

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
