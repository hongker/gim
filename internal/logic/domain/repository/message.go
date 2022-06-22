package repository

import "gim/internal/logic/domain/entity"

type MessageRepo interface {
	Save(item *entity.Message)
	Query(sessionId string, lastMsgId string, n int) []*entity.Message
}
