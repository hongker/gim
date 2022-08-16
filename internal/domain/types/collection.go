package types

import (
	"gim/api"
	"gim/internal/domain/dto"
	"gim/pkg/network"
	"sync"
	"time"
)

type Collection struct {
	bucket *Bucket
}

func (app *Collection) Add(conn *network.Connection) {
	app.bucket.AddChannel(NewChannel(conn))
}

// RegisterConn register user connection
func (app *Collection) RegisterConn(uid string, conn *network.Connection) {
	channel := app.bucket.GetChannel(conn.ID())
	if channel == nil {
		return
	}
	channel.SetKey(uid)
}

// RemoveConn remove connection from bucket
func (app *Collection) RemoveConn(conn *network.Connection) {
	channel := app.bucket.GetChannel(conn.ID())
	if channel == nil {
		return
	}
	app.bucket.RemoveChannel(channel)
}

// CheckConnExist checks if the connection is already exist
func (app *Collection) IsRegistered(conn *network.Connection) bool {
	channel := app.bucket.GetChannel(conn.ID())
	if channel == nil {
		return false
	}
	return channel.Key() == ""
}

func (app *Collection) Refresh(conn *network.Connection, expired time.Duration) {
	channel := app.bucket.GetChannel(conn.ID())
	if channel == nil {
		return
	}
	channel.ResetTimer(expired)

}

// GetUser return the user of connection
func (app *Collection) GetUser(conn *network.Connection) *dto.User {
	channel := app.bucket.GetChannel(conn.ID())
	if channel == nil {
		return nil
	}
	return &dto.User{Id: channel.Key()}
}

// Push push message to session target
func (app *Collection) Push(sessionType string, targetId string, msg []byte) {
	if sessionType == api.UserSession {
		app.pushUser(targetId, msg)
	} else {
		app.pushRoom(targetId, msg)
	}
}

func (app *Collection) pushUser(uid string, msg []byte) {
	channel := app.bucket.GetChannelByKey(uid)
	if channel == nil {
		return
	}
	channel.Conn().Push(msg)
}

func (app *Collection) pushRoom(rid string, msg []byte) {
	room := app.bucket.GetRoom(rid)
	if room == nil {
		return
	}
	room.Push(msg)
}

// Broadcast push message to everyone
func (app *Collection) Broadcast(msg []byte) {
	app.bucket.Push(msg)
}

// JoinRoom
func (app *Collection) JoinRoom(roomId string, conn *network.Connection) {
	channel := app.bucket.GetChannel(conn.ID())
	if channel == nil {
		return
	}
	app.bucket.PutRoom(roomId, channel)
}

// LeaveRoom
func (app *Collection) LeaveRoom(roomId string, conn *network.Connection) {
	channel := app.bucket.GetChannel(conn.ID())
	if channel == nil {
		return
	}

	room := app.bucket.GetRoom(roomId)
	if room == nil {
		return
	}
	room.Remove(channel)
}

// collectionInstance the singleton instance of collection
var collectionInstance struct {
	once       sync.Once
	lock       sync.Mutex
	collection *Collection
}

// Initialize the bucket set.  This can only be done once per binary, subsequent calls are ignored.
func Initialize(bucket *Bucket) {
	collectionInstance.once.Do(func() {
		collectionInstance.collection = &Collection{bucket: bucket}
	})
}

func GetCollection() *Collection {
	collectionInstance.lock.Lock()
	defer collectionInstance.lock.Unlock()
	if collectionInstance.collection == nil {
		Initialize(NewBucket())
	}
	return collectionInstance.collection
}
