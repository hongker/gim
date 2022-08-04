package memory

import (
	"context"
	"gim/internal/domain/dto"
	"gim/internal/domain/entity"
	"gim/internal/domain/repository"
	"gim/pkg/store"
	uuid "github.com/satori/go.uuid"
	"sync"
)

var (
	minScore = store.SCORE(-1)
)

type MessageRepo struct {
	colLock sync.Mutex
	collections map[string]*store.SortedSet
	seqLock sync.Mutex
	sequences map[string]int64
}

func (repo *MessageRepo) Count(ctx context.Context, sessionId string) int {
	return repo.getCollection(sessionId).GetCount()
}

func (repo *MessageRepo) PopMin(ctx context.Context,sessionId string, n int) {
	collection := repo.getCollection(sessionId)
	for i := 0; i < n; i++ {
		collection.PopMin()
	}
}

func (repo *MessageRepo) Save(ctx context.Context, message *entity.Message) error {
	collection := repo.getCollection(message.SessionId)
	message.Id = uuid.NewV4().String()
	collection.AddOrUpdate(message.Id, store.SCORE(message.CreatedAt), message)
	return nil
}

func (repo *MessageRepo) Query(ctx context.Context, query dto.MessageHistoryQuery) ([]entity.Message, error) {
	collection := repo.getCollection(query.SessionId)
	nodes := collection.GetByScoreRange( store.SCORE(query.Last), minScore, &store.GetByScoreRangeOptions{
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

func (repo *MessageRepo) GenerateSequence(sessionId string) int64 {
	repo.seqLock.Lock()
	sequence := repo.sequences[sessionId]
	sequence++
	repo.sequences[sessionId] = sequence
	repo.seqLock.Unlock()
	return sequence
}

func (repo *MessageRepo) getCollection(sessionId string) (*store.SortedSet) {
	repo.colLock.Lock()
	defer repo.colLock.Unlock()
	collection, ok := repo.collections[sessionId]
	if !ok {
		collection = store.New()
		repo.collections[sessionId] = collection
	}
	return collection
}

func NewMessageRepo() repository.MessageRepo {
	return &MessageRepo{
		collections: make(map[string]*store.SortedSet),
		sequences: make(map[string]int64),
	}
}
