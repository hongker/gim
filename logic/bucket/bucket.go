package bucket

import (
	"github.com/ebar-go/ego/utils/runtime"
	"sync"
	"sync/atomic"
)

// Bucket represents a bucket for all connections
type Bucket struct {
	rmu      sync.RWMutex
	channels map[string]*Channel
	*Room

	once sync.Once
	done chan struct{}

	workerNum  uint64
	queueSize  uint64
	queueCount uint64
	queues     []chan QueueItem
}

type QueueItem struct {
	Channel *Channel
	Msg     []byte
}

func (bucket *Bucket) AddChannel(id string) {
	channel := NewChannel(id)
	bucket.rmw.Lock()
	bucket.channels[id] = channel
	bucket.rmw.Unlock()
}
func (bucket *Bucket) RemoveChannel(channel *Channel) {
	bucket.rmw.Lock()
	delete(bucket.channels, channel.ID)
	bucket.rmw.Unlock()
}
func (bucket *Bucket) GetChannel(id string) *Channel {
	bucket.rmw.RLock()
	channel := bucket.channels[id]
	bucket.rmw.RUnlock()
	return channel
}

func (bucket *Bucket) GetOrCreate(id string) *Channel {
	bucket.rmw.Lock()
	channel, exist := bucket.channels[id]
	if !exist {
		channel = NewChannel(id)
		bucket.channels[id] = channel
	}
	bucket.rmw.Unlock()
	return channel
}

func (bucket *Bucket) SubscribeChannel(channel *Channel, sessions ...*Session) {
	for _, session := range sessions {
		channel.AddSession(session)
	}
}

func (bucket *Bucket) UnsubscribeChannel(channel *Channel, sessions ...*Session) {
	for _, session := range sessions {
		channel.RemoveSession(session)
	}
}

func (bucket *Bucket) BroadcastChannel(channel *Channel, msg []byte) {
	num := atomic.AddUint64(&bucket.workerNum, 1) % bucket.queueCount
	bucket.queues[num] <- QueueItem{Channel: channel, Msg: msg}
}

func (bucket *Bucket) Stop() {
	bucket.once.Do(func() {
		close(bucket.done)
	})
}

func (bucket *Bucket) start() {
	for i := 0; i < int(bucket.queueCount); i++ {
		bucket.queues[i] = make(chan QueueItem, bucket.queueSize)
		go func(idx int) {
			defer runtime.HandleCrash()
			bucket.polling(bucket.done, bucket.queues[idx])
		}(i)
	}
}

func (bucket *Bucket) polling(done <-chan struct{}, queue chan QueueItem) {
	for {
		select {
		case <-done:
			return
		case item, ok := <-queue:
			if !ok {
				return
			}

			item.Channel.Broadcast(item.Msg)
		}
	}

}

func NewBucket() *Bucket {
	bucket := &Bucket{
		channels:   make(map[string]*Channel),
		Room:       NewRoom(),
		queues:     make([]chan QueueItem, 32),
		queueSize:  1024,
		queueCount: 32,
	}
	bucket.start()
	return bucket
}
