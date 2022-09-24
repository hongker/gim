package http

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

	g.POST("create", Action(h.create))
	g.POST("update", Action(h.update))
	g.POST("dismiss", Action(h.dismiss))
	g.POST("join", Action(h.join))
	g.POST("leave", Action(h.leave))

}

func (h *ChatroomHandler) create(ctx context.Context, req *dto.ChatroomCreateRequest) (resp *dto.ChatroomCreateResponse, err error) {
	return
}

func (h *ChatroomHandler) update(ctx context.Context, req *dto.ChatroomUpdateRequest) (resp *dto.ChatroomUpdateResponse, err error) {
	return
}

func (h *ChatroomHandler) dismiss(ctx context.Context, req *dto.ChatroomDismissRequest) (resp *dto.ChatroomDismissResponse, err error) {
	return
}
func (h *ChatroomHandler) join(ctx context.Context, req *dto.ChatroomJoinRequest) (res *dto.ChatroomJoinResponse, err error) {
	return
}

func (h *ChatroomHandler) leave(ctx context.Context, req *dto.ChatroomLeaveRequest) (res *dto.ChatroomLeaveResponse, err error) {
	return
}
