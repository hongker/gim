package persistence

import (
	"gim/internal/logic/domain/entity"
	"gim/internal/logic/domain/repository"
	"gim/pkg/store"
	"sync"
)

type MessageRepo struct {
	rmu       sync.RWMutex
	items     map[string]*store.SortedSet
	sequences map[string]int64
}

func NewMessageRepo() repository.MessageRepo {
	repo := &MessageRepo{items: make(map[string]*store.SortedSet, 128), sequences: make(map[string]int64, 128)}
	return repo
}

func (repo *MessageRepo) Save(item *entity.Message) {
	repo.rmu.Lock()
	if _, ok := repo.items[item.SessionId]; !ok {
		repo.items[item.SessionId] = store.New()
	}
	sequence, _ := repo.sequences[item.SessionId]
	item.Sequence = sequence + 1
	repo.sequences[item.SessionId] = item.Sequence

	_, ok := repo.items[item.SessionId]
	if !ok {
		repo.items[item.SessionId] = store.New()
	}
	repo.items[item.SessionId].AddOrUpdate(item.Id, store.SCORE(item.Time), item)
	repo.rmu.Unlock()
}

func (repo *MessageRepo) Query(sessionId string, lastMsgId string, n int) []*entity.Message {
	repo.rmu.RLock()
	defer repo.rmu.RUnlock()

	var start, end store.SCORE
	if lastMsgId == "" { // 查询最近n条消息
		start = store.MaxScore
	} else {
		node := repo.items[sessionId].GetByKey(lastMsgId)
		if node != nil {
			start = node.Score()
		}
	}

	set, ok := repo.items[sessionId]
	if !ok {
		return nil
	}

	nodes := set.GetByScoreRange(start, end, &store.GetByScoreRangeOptions{Limit: n})
	result := make([]*entity.Message, 0, n)
	for _, node := range nodes {
		result = append(result, node.Value.(*entity.Message))
	}

	return result

}
