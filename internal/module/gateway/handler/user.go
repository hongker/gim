package handler

import (
	"context"
	"gim/internal/module/gateway/application"
	"gim/internal/module/gateway/domain/dto"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	groupName string
	userApp   application.UserApplication
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		groupName: "user",
		userApp:   application.NewUserApplication(),
	}
}

func (h *UserHandler) Install(route *gin.Engine) {
	g := route.Group(h.groupName)
	g.POST("login", Action(h.login))
}

// login represents user login request.
func (h *UserHandler) login(ctx context.Context, req *dto.UserLoginRequest) (resp *dto.UserLoginResponse, err error) {
	return h.userApp.Login(ctx, req)
}
