package entity

import "sync"

// Bucket 存储channel和rooms
type Bucket struct {
	rmu      sync.RWMutex
	channels map[string]*Channel // 对应私聊,key为uid
	rooms    map[string]*Room    // 对应群聊,key为groupId
}

func NewBucket() *Bucket {
	return &Bucket{
		channels: make(map[string]*Channel, 1024),
		rooms:    make(map[string]*Room, 128),
	}
}

func (bucket *Bucket) GetChannel(connId string) *Channel {
	bucket.rmu.RLock()
	defer bucket.rmu.RUnlock()
	return bucket.channels[connId]
}

func (bucket *Bucket) Channels() map[string]*Channel {
	bucket.rmu.RLock()
	defer bucket.rmu.RUnlock()
	return bucket.channels
}

func (bucket *Bucket) GetChannelByKey(key string) *Channel {
	bucket.rmu.RLock()
	defer bucket.rmu.RUnlock()
	for _, channel := range bucket.channels {
		if channel.Key() == key {
			return channel
		}
	}

	return nil
}

func (bucket *Bucket) AddChannel(channel *Channel) {
	bucket.rmu.Lock()
	bucket.channels[channel.conn.ID()] = channel
	bucket.rmu.Unlock()
}

func (bucket *Bucket) RemoveChannel(channel *Channel) {
	bucket.rmu.Lock()
	delete(bucket.channels, channel.conn.ID())
	bucket.rmu.Unlock()
}

func (bucket *Bucket) GetRoom(roomId string) *Room {
	bucket.rmu.RLock()
	defer bucket.rmu.RUnlock()
	return bucket.rooms[roomId]
}

func (bucket *Bucket) PutRoom(roomId string, channel *Channel) {
	bucket.rmu.Lock()
	room := bucket.rooms[roomId]
	if room == nil {
		room = NewRoom(roomId)
	}
	room.Add(channel)
	bucket.rmu.Unlock()
}
