package queue

import (
	"container/list"
	"time"
)
// Queue 批处理队列
type Queue struct {
	l *list.List
	cap int // 队列容量
	limit bool // 是否限制队列长度，如果为true,则标识当队列长度超过指定长度时，需要删除多余的数据
}

func (q *Queue) Offer(item interface{}) {
	if q.limit && q.l.Len() == q.cap {
		q.l.Remove(q.l.Front())
	}
	q.l.PushBack(item)
}

func (q *Queue) Poll(duration time.Duration, fn func(items []interface{}) )   {
	timer := time.NewTimer(duration)
	for {
		select {
		case <- timer.C: // 按时间来触发批处理
			items := make([]interface{},0, q.cap)
			for next := q.l.Front(); next != nil; next = next.Next() {
				items = append(items, next.Value)
			}
			q.l.Init()
			if len(items) == 0 {
				continue
			}
			fn(items)
			timer.Reset(duration)
		default:
			if q.limit || q.l.Len() < q.cap {
				continue
			}
			// 按容量来触发批处理
			items := make([]interface{},0, q.l.Len())
			for next := q.l.Front(); next != nil; next = next.Next() {
				items = append(items, next.Value)
			}
			q.l.Init()
			fn(items)
		}
	}
}


func NewQueue(cap int, limit bool) *Queue {
	return &Queue{cap: cap, l : list.New(), limit: limit}
}


