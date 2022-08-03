package store

import (
	"sync"
)

type Hash struct {
	mu sync.RWMutex
	items map[string]interface{}
}

func (s *Hash) Get(key string) (interface{}, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, ok := s.items[key]
	return item, ok
}

func (s *Hash) Save(key string, value interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[key] = value
	return nil
}

func NewHash() *Hash {
	return &Hash{
		items: make(map[string]interface{}),
	}
}
