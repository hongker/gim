package http

import (
	"gim/internal/module/gateway/domain/types"
	"github.com/gin-gonic/gin"
	"sync"
)

type container struct {
	handlers []Handler
}

// RegisterHandler registers a new handler.
func (c *container) RegisterHandler(handlers ...Handler) {
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
	return &container{handlers: make([]Handler, 0, 16)}
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

func RouteLoader(router *gin.Engine) {
	// register handlers
	Container().RegisterHandler(
		NewUserHandler(),
		NewMessageHandler(),
		NewChatRoomHandler(),
	)

	router.Use(recoverMiddleware(), checkToken(types.DefaultAuthenticator()))

	Container().install(router)
}
