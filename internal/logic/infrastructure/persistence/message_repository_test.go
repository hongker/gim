package persistence

import (
	"fmt"
	"gim/internal/logic/domain/entity"
	"testing"
	"time"
)

func TestMessageRepo(t *testing.T) {
	repo := NewMessageRepo()
	sessionId := "1:1001"
	for i := 0; i < 10; i++ {
		repo.Save(sessionId, entity.Message{
			Id:      fmt.Sprintf("id:%d", i+1000),
			Content: fmt.Sprintf("content:%d", i+1000),
			Time:    time.Now().UnixNano(),
		})
	}

	items := repo.Query(sessionId, "", 5)
	for _, item := range items {
		fmt.Println("message:", item)
	}
}

func BenchmarkSaveMessage(b *testing.B) {
	repo := NewMessageRepo()
	sessionId := "1:1001"

	for i := 0; i < b.N; i++ {
		repo.Save(sessionId, entity.Message{
			Id:      fmt.Sprintf("id:%d", i),
			Content: fmt.Sprintf("content:%d", i),
			Time:    time.Now().UnixNano(),
		})
	}
}
