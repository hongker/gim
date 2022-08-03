package aggregate

import (
	"gim/api"
	"gim/internal/domain/dto"
	"gim/internal/domain/types"
	"gim/pkg/network"
)

type GateApp struct {
	bucket *types.Bucket
}

// RegisterConn register user connection
func (app *GateApp) RegisterConn(uid string, conn *network.Connection) {
	channel := types.NewChannel(uid, conn)
	app.bucket.AddChannel(channel)
}
// RemoveConn remove connection from bucket
func (app *GateApp) RemoveConn(conn *network.Connection)   {
	channel := app.bucket.GetChannel(conn.ID())
	if channel == nil {
		return
	}
	app.bucket.RemoveChannel(channel)
}

// CheckConnExist checks if the connection is already exist
func (app *GateApp) CheckConnExist(conn *network.Connection) bool {
	channel := app.bucket.GetChannel(conn.ID())
	return channel != nil
}

// GetUser return the user of connection
func (app *GateApp) GetUser(conn *network.Connection) *dto.User {
	channel := app.bucket.GetChannel(conn.ID())
	if channel == nil {
		return nil
	}
	return &dto.User{Id: channel.Key()}
}

// Push push message to session target
func (app *GateApp) Push(sessionType string, sessionId string, msg []byte) {
	if sessionType== api.UserSession {
		app.pushUser(sessionId, msg)
	}else {
		app.pushRoom(sessionId, msg)
	}
}

func (app *GateApp) pushUser(uid string, msg []byte) {
	channel := app.bucket.GetChannelByKey(uid)
	if channel == nil {
		return
	}
	channel.Conn().Push(msg)
}

func (app *GateApp) pushRoom(rid string, msg []byte) {
	room := app.bucket.GetRoom(rid)
	if room == nil {
		return
	}
	room.Push(msg)
}
// Broadcast push message to everyone
func (app *GateApp) Broadcast(msg []byte) {
	app.bucket.Push(msg)
}

// JoinRoom
func (app *GateApp) JoinRoom(roomId string, conn *network.Connection, ) {
	channel := app.bucket.GetChannel(conn.ID())
	if channel == nil {
		return
	}
	app.bucket.PutRoom(roomId, channel)
}
// LeaveRoom
func (app *GateApp) LeaveRoom(roomId string, conn *network.Connection, ) {
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


func NewGateApp() *GateApp {
	app := &GateApp{
		bucket: types.NewBucket(),
	}
	return app
}