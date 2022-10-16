package task

import (
	"container/list"
	"sync"
	"time"
)

// Queue 批处理队列
type Queue[T any] struct {
	l     *list.List
	cap   int // 队列容量
	mu    sync.Mutex
	limit bool // 是否限制队列长度，如果为true,则标识当队列长度超过指定长度时，需要删除多余的数据
}

func (q *Queue[T]) Offer(item T) {
	q.mu.Lock()
	if q.limit && q.l.Len() == q.cap {
		q.l.Remove(q.l.Front())
	}
	q.l.PushBack(item)
	q.mu.Unlock()
}

func (q *Queue[T]) Poll(duration time.Duration, fn func(items []T)) {
	timer := time.NewTimer(duration)
	for {
		select {
		case <-timer.C: // 按时间来触发批处理
			items := make([]T, 0, q.cap)
			q.mu.Lock()

			for next := q.l.Front(); next != nil; next = next.Next() {
				items = append(items, next.Value.(T))
			}
			q.l.Init()
			if len(items) > 0 {
				fn(items)
				timer.Reset(duration)
			}
			q.mu.Unlock()
		default:
			//if q.limit || q.l.Len() < q.cap {
			//	continue
			//}
			//q.mu.Lock()
			//// 按容量来触发批处理
			//items := make([]T, 0, q.l.Len())
			//for next := q.l.Front(); next != nil; next = next.Next() {
			//	items = append(items, next.Value)
			//}
			//q.l.Init()
			//fn(items)
			//q.mu.Unlock()
		}
	}
}

func NewQueue[T any](cap int, limit bool) *Queue[T] {
	return &Queue[T]{cap: cap, l: list.New(), limit: limit}
}
