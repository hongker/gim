package applications

import (
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

func (app *GateApp) GetUser(conn *network.Connection) string {
	channel := app.bucket.GetChannel(conn.ID())
	if channel == nil {
		return ""
	}
	return channel.Key()
}

func (app *GateApp) PushUser() {

}

func (app *GateApp) PushGroup() {}

func (app *GateApp) Broadcast() {}

func NewGateApp() *GateApp {
	return &GateApp{
		bucket: entity.NewBucket(),
	}
}