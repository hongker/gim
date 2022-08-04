package persistence

import (
	"context"
	"fmt"
	"gim/internal/domain/dto"
	"gim/internal/domain/entity"
	"gim/internal/domain/repository"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

type MessageRepo struct {
	redisConn redis.UniversalClient
	expired time.Duration
}

var (
	messageCachePrefix = "message"
	sequencePrefix = "sequence"
)
func (repo MessageRepo) getCacheKey(sessionId string) string {
	return fmt.Sprintf("%s:%s", messageCachePrefix, sessionId)
}
func (repo MessageRepo) getSequenceCacheKey(sessionId string) string {
	return fmt.Sprintf("%s:%s", sequencePrefix, sessionId)
}

func (repo MessageRepo) Save(ctx context.Context, message *entity.Message) error {
	err := repo.redisConn.ZAdd(ctx, repo.getCacheKey(message.SessionId), &redis.Z{
		Score: float64(message.CreatedAt),
		Member: entity.Encode(message),
	}).Err()
	return err
}

func (repo MessageRepo) Query(ctx context.Context, query dto.MessageHistoryQuery) ([]entity.Message, error) {
	items, err := repo.redisConn.ZRevRangeByScore(ctx, repo.getCacheKey(query.SessionId), &redis.ZRangeBy{
		Min:    "-1",
		Max:    strconv.FormatInt(query.Last, 10),
		Offset: 0,
		Count: int64( query.Limit),
	}).Result()
	if err != nil {
		return nil, err
	}
	res := make([]entity.Message, 0, query.Limit)
	for _, item := range items {
		message := entity.Message{}
		if err := entity.Decode([]byte(item), &message); err != nil {
			continue
		}
		res = append(res, message)
	}
	return res, nil
}

func (repo MessageRepo) Count(ctx context.Context, sessionId string) int {
	count, _ := repo.redisConn.ZCard(ctx, repo.getCacheKey(sessionId)).Result()
	return int(count)
}

func (repo MessageRepo) PopMin(ctx context.Context, sessionId string, n int) {
	repo.redisConn.ZPopMin(ctx, repo.getCacheKey(sessionId), int64(n))
}

func (repo MessageRepo) GenerateSequence(ctx context.Context, sessionId string) int64 {
	res, _ := repo.redisConn.Incr(ctx, repo.getSequenceCacheKey(sessionId)).Result()
	return res
}

func NewMessageRepo(redisConn redis.UniversalClient) repository.MessageRepo {
	return &MessageRepo{redisConn: redisConn, expired: time.Hour* 24 * 30}
}
