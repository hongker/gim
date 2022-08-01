package cache

import (
	"context"
	"gim/internal/domain/dto"
	"gim/internal/domain/entity"
	"gim/internal/domain/repository"
	"gim/pkg/store"
)

var (
	minScore = store.SCORE(-1)
)

type MessageRepo struct {
	collections map[string]*store.SortedSet
}


func (repo *MessageRepo) Save(ctx context.Context, message *entity.Message) error {

	repo.getCollection(message.SessionId).AddOrUpdate(message.Id, store.SCORE(message.CreatedAt), message)
	return nil
}

func (repo *MessageRepo) Query(ctx context.Context, query dto.MessageHistoryQuery) ([]entity.Message, error) {
	collection := repo.getCollection(query.SessionId)
	nodes := collection.GetByScoreRange(minScore, store.SCORE(query.Last), &store.GetByScoreRangeOptions{
		Limit:        query.Limit,
		ExcludeStart: false,
		ExcludeEnd:   false,
	})

	res := make([]entity.Message, 0, query.Limit)
	for _, node := range nodes {
		item := node.Value.(*entity.Message)
		res = append(res, *item)
	}
	return res, nil
}

func (repo *MessageRepo) getCollection(sessionId string) (*store.SortedSet) {
	collection, ok := repo.collections[sessionId]
	if !ok {
		collection = store.New()
		repo.collections[sessionId] = collection
	}
	return collection
}

func NewMessageRepo() repository.MessageRepo {
	return &MessageRepo{collections: make(map[string]*store.SortedSet)}
}
