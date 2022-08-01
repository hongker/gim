package applications

import (
	"context"
	"gim/internal/domain/dto"
	"gim/internal/domain/entity"
	"gim/internal/domain/repository"
	uuid "github.com/satori/go.uuid"
	"time"
)

type MessageApp struct {
	repo repository.MessageRepo
}


func (app *MessageApp) Send(ctx context.Context, fromUser *entity.User, req *dto.MessageSendRequest) error {
	item := &entity.Message{
		Id:          uuid.NewV4().String(),
		Type:        req.Type,
		Content:     req.Content,
		CreatedAt:   time.Now().UnixNano(),
		ClientMsgId: req.ClientMsgId,
		Sequence:    0,
		SessionId:   req.SessionId,
		FromUser:    fromUser,
	}
	return app.repo.Save(ctx, item)
}



func (app *MessageApp) Query(ctx context.Context,req *dto.MessageQueryRequest) (*dto.MessageQueryResponse, error) {
	items, err := app.repo.Query(ctx, dto.MessageHistoryQuery{
		SessionId: req.SessionId,
		Limit:     10,
		Last:      req.Last,
	})
	if err != nil {
		return nil, err
	}

	res := &dto.MessageQueryResponse{Items: make([]dto.Message, 0, len(items))}
	for _, item := range items {
		res.Items = append(res.Items, dto.Message{
			SessionId: item.SessionId,
			Content:   item.Content,
			CreatedAt: item.CreatedAt,
		})
	}

	return nil, nil
}

func NewMessageApp() *MessageApp {
	return &MessageApp{}
}