package repository

import (
	"context"
	"gim/internal/module/gateway/domain/types"
)

type MessageRepository interface {
	Save(ctx context.Context, msg *types.Message) error
}

type SessionRepository interface {
	SaveMessage(ctx context.Context, session *types.Session, msg *types.Message) error
	Query(ctx context.Context, session *types.Session)
}
