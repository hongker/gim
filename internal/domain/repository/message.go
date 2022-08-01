package repository

import (
	"context"
	"gim/internal/domain/dto"
	"gim/internal/domain/entity"
)

type MessageRepo interface {
	Save(ctx context.Context, message *entity.Message) error
	Query(ctx context.Context, query dto.MessageHistoryQuery) ([]entity.Message, error)
}
