package handler

import (
	"context"
	"gim/internal/module/gateway/domain/dto"
	"github.com/gin-gonic/gin"
)

type ChatroomHandler struct {
	groupName string
}

func NewChatRoomHandler() *ChatroomHandler {
	return &ChatroomHandler{
		groupName: "chatroom",
	}
}

func (h *ChatroomHandler) Install(router *gin.Engine) {
	g := router.Group(h.groupName)

	g.POST("join", Action(h.join))
	g.POST("leave", Action(h.leave))

}

func (h *ChatroomHandler) join(ctx context.Context, req *dto.ChatroomJoinRequest) (res *dto.ChatroomJoinResponse, err error) {
	return
}

func (h *ChatroomHandler) leave(ctx context.Context, req *dto.ChatroomLeaveRequest) (res *dto.ChatroomLeaveResponse, err error) {
	return
}
