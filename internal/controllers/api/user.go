package api

import (
	"context"
	"gim/internal/application"
	"gim/internal/domain/dto"
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
	g.GET("", Action(h.find))
}

func (h *UserHandler) find(ctx context.Context, req *dto.UserFindRequest) (resp *dto.UserFindResponse, err error) {
	return h.userApp.Find(ctx, req)
}
