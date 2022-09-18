package route

import (
	"gim/internal/module/gateway/handler"
	"github.com/gin-gonic/gin"
	"sync"
)

type container struct {
	handlers []handler.Handler
}

// RegisterHandler registers a new handler.
func (c *container) RegisterHandler(handlers ...handler.Handler) {
	c.handlers = append(c.handlers, handlers...)
}

// install prepare install handler with router
func (c *container) install(router *gin.Engine) {
	for _, h := range c.handlers {
		h.Install(router)
	}
}

// buildContainer creates a new container instance.
func buildContainer() *container {
	return &container{handlers: make([]handler.Handler, 0, 16)}
}

// containerInstance represents a container instance of singleton mode.
var containerInstance struct {
	once      sync.Once
	container *container
}

// Container returns the container instance. Use singleton mode
func Container() *container {
	containerInstance.once.Do(func() {
		containerInstance.container = buildContainer()
	})
	return containerInstance.container
}

func Loader(router *gin.Engine) {
	router.Use(recoverMiddleware())

	Container().install(router)
}
