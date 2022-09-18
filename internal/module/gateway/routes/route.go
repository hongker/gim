package routes

import (
	"gim/internal/module/gateway/handler"
	"github.com/gin-gonic/gin"
	"sync"
)

type container struct {
	handlers []handler.Handler
}

func (c *container) UseHandler(handlers ...handler.Handler) {
	c.handlers = append(c.handlers, handlers...)
}

func (c *container) install(router *gin.Engine) {
	for _, h := range c.handlers {
		h.Install(router)
	}
}

var containerInstance struct {
	once      sync.Once
	container *container
}

func Container() *container {
	containerInstance.once.Do(func() {
		containerInstance.container = &container{handlers: make([]handler.Handler, 0, 16)}
	})
	return containerInstance.container
}

func Loader(router *gin.Engine) {
	Container().install(router)
}
