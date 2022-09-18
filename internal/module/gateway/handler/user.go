package handler

import (
	"gim/internal/module/gateway/application"
	"gim/internal/module/gateway/domain/dto"
	"gim/internal/module/gateway/render"
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
	g.POST("login", h.login)
}

// login represents user login request.
func (h *UserHandler) login(ctx *gin.Context) {
	req := &dto.UserLoginRequest{}
	err := render.SerializeRequestFromContext(ctx, req)
	render.Abort(err)

	response, err := h.userApp.Login(ctx, req)
	render.Abort(err)

	render.Success(ctx, response)
}
