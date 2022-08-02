package applications

import (
	"gim/api"
	"gim/internal/domain/dto"
	"gim/internal/domain/entity"
	"gim/pkg/network"
)

type GateApp struct {
	bucket *entity.Bucket
}

func (app *GateApp) RegisterConn(uid string, conn *network.Connection) {
	channel := entity.NewChannel(uid, conn)
	app.bucket.AddChannel(channel)
}

func (app *GateApp) RemoveConn(conn *network.Connection)   {
	channel := app.bucket.GetChannel(conn.ID())
	if channel == nil {
		return
	}
	app.bucket.RemoveChannel(channel)
}

func (app *GateApp) CheckConnExist(conn *network.Connection) bool {
	channel := app.bucket.GetChannel(conn.ID())
	return channel != nil
}

func (app *GateApp) GetUser(conn *network.Connection) *dto.User {
	channel := app.bucket.GetChannel(conn.ID())
	if channel == nil {
		return nil
	}
	return &dto.User{Id: channel.Key()}
}


func (app *GateApp) pushUser(uid string, msg []byte) {
	channel := app.bucket.GetChannelByKey(uid)
	if channel == nil {
		return
	}
	channel.Conn().Push(msg)
}

func (app *GateApp) Push(sessionType string, sessionId string, msg []byte) {
	if sessionType== api.PrivateMessage {
		app.pushUser(sessionId, msg)
	}else {
		app.pushRoom(sessionId, msg)
	}
}

func (app *GateApp) pushRoom(rid string, msg []byte) {
	room := app.bucket.GetRoom(rid)
	if room == nil {
		return
	}
	channels := room.Channels()
	for _, channel := range channels {
		channel.Conn().Push(msg)
	}
}

func (app *GateApp) Broadcast(msg []byte) {
	for _, channel := range app.bucket.Channels() {
		channel.Conn().Push(msg)
	}
}

func (app *GateApp) JoinRoom(roomId string, conn *network.Connection, ) {
	channel := app.bucket.GetChannel(conn.ID())
	if channel == nil {
		return
	}
	app.bucket.PutRoom(roomId, channel)
}


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
	return &GateApp{
		bucket: entity.NewBucket(),
	}
}