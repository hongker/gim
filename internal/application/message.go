package application

import (
	"context"
	"gim/api"
	"gim/internal/domain/dto"
	"gim/internal/domain/entity"
	"gim/internal/domain/event"
	"gim/internal/domain/repository"
	"gim/internal/infrastructure/config"
	"gim/pkg/errors"
	"gim/pkg/queue"
	"sync"
	"time"
)

type MessageApp struct {
	repo repository.MessageRepo
	groupRepo repository.GroupRepo
	rmu sync.RWMutex
	queues map[string]*queue.Queue
	queueCap int
	messageCap int

}

func (app *MessageApp) getQueue(sessionType string, targetId string) *queue.Queue {
	app.rmu.Lock()
	defer app.rmu.Unlock()
	if q, ok := app.queues[targetId]; ok {
		return q
	}
	limit := true
	if sessionType == api.UserSession {
		limit = false
	}
	q := queue.NewQueue(app.queueCap, limit)
	app.queues[targetId] = q
	go q.Poll(time.Second , func(items []interface{}) {
		batchMessages := &dto.BatchMessage{Items: make([]dto.Message, len(items))}
		for i, item := range items {
			batchMessages.Items[i] = item.(dto.Message)
		}

		event.Trigger(event.Push, sessionType, targetId, batchMessages)
	})
	return q
}

func (app *MessageApp) Send(ctx context.Context, sender *dto.User, req *dto.MessageSendRequest) ( error) {
	if err := app.validate(ctx, sender, req); err != nil {
		return err
	}
	sessionId := req.SessionId(sender.Id)
	item := &entity.Message{
		SessionType: req.Type,
		Content:     req.Content,
		ContentType: req.ContentType,
		CreatedAt:   time.Now().UnixNano(),
		RequestId: req.RequestId,
		Sequence:    app.repo.GenerateSequence(ctx, sessionId),
		SessionId:   sessionId,
		FromUser:    &entity.User{Id: sender.Id},
	}
	if err := app.repo.Save(ctx, item); err != nil {
		return  err
	}

	// 超过容量后删除早期数据
	count := app.repo.Count(ctx, item.SessionId)
	if diff := app.messageCap - count; diff > 0 {
		app.repo.PopMin(ctx, item.SessionId, diff)
	}
	res := dto.Message{
		Id: item.Id,
		RequestId: item.RequestId,
		Session: dto.Session{Id: item.SessionId, Type: item.SessionType},
		Content:   item.Content,
		ContentType: item.ContentType,
		CreatedAt: item.CreatedAt,
		Sequence: item.Sequence,
		FromUser: dto.User{
			Id: item.FromUser.Id,
		},
	}

	app.getQueue(item.SessionType, req.TargetId).Offer(res)

	return  nil
}

func (app *MessageApp) validate(ctx context.Context, sender *dto.User, req *dto.MessageSendRequest) error{
	if req.Type == api.UserSession {
		if sender.Id ==req.TargetId {
			return errors.InvalidParameter("sender same as target")
		}
	} else if req.Type == api.GroupSession {
		_, err := app.groupRepo.Find(ctx, req.TargetId)
		if err != nil {
			return err
		}

	}
	return nil
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
			RequestId: item.RequestId,
			Session: dto.Session{Id: item.SessionId, Type: item.SessionType},
			Content:   item.Content,
			ContentType: item.ContentType,
			CreatedAt: item.CreatedAt,
			Sequence: item.Sequence,
			FromUser: dto.User{Id: item.FromUser.Id},
		})
	}

	return res, nil
}

func NewMessageApp(repo repository.MessageRepo, groupRepo repository.GroupRepo, config *config.Config) *MessageApp {
	app := &MessageApp{
		repo: repo,
		groupRepo: groupRepo,
		queues: map[string]*queue.Queue{},
		queueCap: config.Message.PushCount,
		messageCap: config.Message.MaxStoreSize,
	}
	return app
}