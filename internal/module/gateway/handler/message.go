package handler

import (
	"context"
	"gim/internal/module/gateway/domain/dto"
	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	groupName string
}

func NewMessageHandler() *MessageHandler {
	return &MessageHandler{}
}

func (m *MessageHandler) Install(router *gin.Engine) {
	g := router.Group(m.groupName)
	g.POST("session", Action(m.session))
	g.POST("query", Action(m.query))
	g.POST("send", Action(m.send))
}

// session returns user session response.
func (m *MessageHandler) session(ctx context.Context, req *dto.SessionQueryRequest) (resp *dto.SessionQueryResponse, err error) {
	return
}

// query returns message history response.
func (m *MessageHandler) query(ctx context.Context, req *dto.MessageQueryRequest) (resp *dto.MessageQueryResponse, err error) {
	return
}

// send process user message send request.
func (m *MessageHandler) send(ctx context.Context, req *dto.MessageSendResponse) (resp *dto.MessageSendResponse, err error) {
	return
}
