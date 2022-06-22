package store

import "container/list"

type Queue struct {
	store *list.List
}

func NewQueue() *Queue {
	return &Queue{store: list.New()}
}

func (queue Queue) Push(v interface{}) {
	queue.store.PushBack(v)
}

func (queue Queue) Last() interface{} {
	return queue.store.Back().Value
}

func (queue Queue) Length() int {
	return queue.store.Len()
}

func (queue Queue) Query(v interface{}, n int) []interface{} {
	result := make([]interface{}, 0, n)
	item := queue.store.Front()
	for item != nil {
		if item.Value != v {
			item = item.Next()
			continue
		}

		for i := 0; i < n; i++ {
			item = item.Next()
			result = append(result, item.Value)
		}
		break
	}

	return result
}
