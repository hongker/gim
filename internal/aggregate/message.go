package aggregate

import (
	"context"
	"gim/internal/domain/dto"
	"gim/internal/domain/entity"
	"gim/internal/domain/repository"
	"time"
)

type MessageApp struct {
	repo repository.MessageRepo
}


func (app *MessageApp) Send(ctx context.Context, fromUser *entity.User, req *dto.MessageSendRequest) (*dto.Message, error) {
	item := &entity.Message{
		Type:        req.Type,
		Content:     req.Content,
		ContentType: req.ContentType,
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
		Id: item.Id,
		SessionId: item.SessionId,
		Content:   item.Content,
		ContentType: item.ContentType,
		CreatedAt: item.CreatedAt,
		Sequence: item.Sequence,
	}
	return res, nil
}


func (app *MessageApp) Query(ctx context.Context,req *dto.MessageQueryRequest) (*dto.MessageQueryResponse, error) {
	items, err := app.repo.Query(ctx, dto.MessageHistoryQuery{
		SessionId: req.SessionId,
		Limit:     req.Limit,
		Last:      req.Last,
	})
	if err != nil {
		return nil, err
	}

	res := &dto.MessageQueryResponse{Items: make([]dto.Message, 0, len(items))}
	for _, item := range items {
		res.Items = append(res.Items, dto.Message{
			Id: item.Id,
			SessionId: item.SessionId,
			Content:   item.Content,
			ContentType: item.ContentType,
			CreatedAt: item.CreatedAt,
			Sequence: item.Sequence,
		})
	}

	return res, nil
}

func NewMessageApp(repo repository.MessageRepo) *MessageApp {
	return &MessageApp{repo: repo}
}