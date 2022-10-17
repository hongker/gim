package bucket

import (
	"log"
	"sync"
)

// Bucket represents a bucket for all connections
type Bucket struct {
	rmu      sync.RWMutex
	channels map[string]*Channel
	*Room
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
	channel.stop()
}
func (bucket *Bucket) GetChannel(id string) *Channel {
	bucket.rmw.RLock()
	channel := bucket.channels[id]
	bucket.rmw.RUnlock()
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

func (bucket *Bucket) Stop() {
	bucket.Room.stop()
	for _, channel := range bucket.channels {
		channel.Room.stop()
	}
}

func NewBucket() *Bucket {
	return &Bucket{
		channels: make(map[string]*Channel),
		Room:     NewRoom(1024, 32),
	}
}

// Session represents a session for one connection
type Session struct {
	ID string
}

func NewSession(id string) *Session {
	return &Session{ID: id}
}

func (session *Session) Send(msg []byte) {
	log.Println("session send msg: ", session.ID, string(msg))
}
