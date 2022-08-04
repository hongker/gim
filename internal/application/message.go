package application

import (
	"context"
	"gim/api"
	"gim/internal/domain/dto"
	"gim/internal/domain/entity"
	"gim/internal/domain/event"
	"gim/internal/domain/repository"
	"gim/pkg/queue"
	"sync"
	"time"
)

type MessageApp struct {
	repo repository.MessageRepo
	rmu sync.RWMutex
	queues map[string]*queue.Queue
	queueCap int
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
	sessionId := req.SessionId(sender.Id)
	item := &entity.Message{
		SessionType: req.Type,
		Content:     req.Content,
		ContentType: req.ContentType,
		CreatedAt:   time.Now().UnixNano(),
		ClientMsgId: req.ClientMsgId,
		Sequence:    app.repo.GenerateSequence(sessionId),
		SessionId:   sessionId,
		FromUser:    &entity.User{Id: sender.Id},
	}
	if err := app.repo.Save(ctx, item); err != nil {
		return  err
	}
	res := dto.Message{
		Id: item.Id,
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

func NewMessageApp(repo repository.MessageRepo) *MessageApp {
	app := &MessageApp{repo: repo, queues: map[string]*queue.Queue{}, queueCap: 10}
	return app
}