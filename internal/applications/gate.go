package applications

import (
	"gim/api"
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


func (app *GateApp) PushUser(uid string, msg []byte) {
	channel := app.bucket.GetChannelByKey(uid)
	if channel == nil {
		return
	}
	channel.Conn().Push(msg)
}

func (app *GateApp) Push(sessionType string, sessionId string, msg []byte) {
	if sessionType== api.PrivateMessage {
		app.PushUser(sessionId, msg)
	}else {
		app.PushRoom(sessionId, msg)
	}
}

func (app *GateApp) PushRoom(rid string, msg []byte) {
	room := app.bucket.GetRoom(rid)
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

func NewGateApp() *GateApp {
	return &GateApp{
		bucket: entity.NewBucket(),
	}
}