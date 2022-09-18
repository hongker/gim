package handler

import (
	"github.com/gin-gonic/gin"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) Install(route *gin.Engine) {
	g := route.Group("user")
	g.POST("login", h.login)
}

func (h *UserHandler) login(ctx *gin.Context) {

}
