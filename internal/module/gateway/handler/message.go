package handler

import (
	"context"
	"gim/internal/module/gateway/domain/dto"
	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	groupName string
}

func (m *MessageHandler) Install(router *gin.Engine) {
	g := router.Group(m.groupName)
	g.POST("query", Action(m.query))
	g.POST("session", Action(m.session))
}

// query returns message history response.
func (m *MessageHandler) query(ctx context.Context, req *dto.MessageQueryRequest) (resp *dto.MessageQueryResponse, err error) {
	return
}

// session returns user session response.
func (m *MessageHandler) session(ctx context.Context, req *dto.SessionQueryRequest) (resp *dto.SessionQueryResponse, err error) {
	return
}
