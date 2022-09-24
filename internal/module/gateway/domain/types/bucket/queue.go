package bucket

import (
	"gim/pkg/queue"
	"sync"
	"time"
)

type Queue interface {
	Offer(item interface{})
}

type QueueDispatcher interface {
	GetQueue(id string) Queue
}

type DelayedQueueDispatcher struct {
	mu       sync.Mutex
	queues   map[string]*queue.Queue
	cap      int
	duration time.Duration
	callback func(items []interface{})
}

func (d *DelayedQueueDispatcher) GetQueue(id string) Queue {
	d.mu.Lock()
	if q, exist := d.queues[id]; exist {
		return q
	}
	q := queue.NewQueue(d.cap, false)
	go q.Poll(d.duration, d.callback)
	return nil
}

var queueDispatcherOnce struct {
	once     sync.Once
	instance QueueDispatcher
}

func GetQueueDispatcher() QueueDispatcher {
	queueDispatcherOnce.once.Do(func() {
		qd := &DelayedQueueDispatcher{
			queues: make(map[string]*queue.Queue),
		}

		queueDispatcherOnce.instance = qd
	})
	return queueDispatcherOnce.instance
}
