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


func (app *MessageApp) Send(ctx context.Context, fromUser *entity.User, req *dto.MessageSendRequest) (*dto.Message, error) {
	item := &entity.Message{
		Id:          uuid.NewV4().String(),
		Type:        req.Type,
		Content:     req.Content,
		CreatedAt:   time.Now().UnixNano(),
		ClientMsgId: req.ClientMsgId,
		Sequence:    app.repo.GenerateSequence(req.SessionId),
		SessionId:   req.SessionId,
		FromUser:    fromUser,
	}
	if err := app.repo.Save(ctx, item); err != nil {
		return nil, err
	}
	res := &dto.Message{
		SessionId: item.SessionId,
		Content:   item.Content,
		CreatedAt: item.CreatedAt,
	}
	return res, nil
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

func NewMessageApp(repo repository.MessageRepo) *MessageApp {
	return &MessageApp{repo: repo}
}