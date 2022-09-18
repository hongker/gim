package render

import "github.com/gin-gonic/gin"

func ReadRequestBody(ctx *gin.Context) {

}

func RequestBodyFromContext(ctx *gin.Context) []byte {
	return nil
}

func SerializeRequestFromContext(ctx *gin.Context, container interface{}) error {
	body := RequestBodyFromContext(ctx)
	return serializer().Decode(body, container)
}
