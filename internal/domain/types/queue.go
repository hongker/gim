package types

import (
	"container/list"
	"time"
)

type Queue struct {
	size int
	l *list.List
}

func (q *Queue) Offer(item interface{}) {
	if q.l.Len() == q.size {
		q.l.Remove(q.l.Front())
	}
	q.l.PushBack(item)
}

func (q *Queue) Poll(duration time.Duration, fn func(items []interface{}) )   {
	timer := time.NewTimer(duration)


	for {
		select {
		case <- timer.C:
			items := make([]interface{},0, q.size)
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

		}
	}
}


func NewQueue(size int) *Queue {
	return &Queue{size: size, l : list.New()}
}
