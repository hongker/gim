package route

import (
	"gim/internal/module/gateway/handler"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestContainer(t *testing.T) {
	c := Container()
	assert.NotNil(t, c)

	// test singleton
	c1 := Container()
	assert.Equal(t, c1, c)
}

func TestContainer_RegisterHandler(t *testing.T) {
	type mockHandler struct {
		mock.Mock
		handler.Handler
	}

	h := new(mockHandler)
	h.On("Install", mock.Anything)
	Container().RegisterHandler(h)
}

func TestLoader(t *testing.T) {
	router := gin.Default()
	Loader(router)
}
