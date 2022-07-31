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

func (app *GateApp) RemoveChannel(connId string)   {
	channel := app.bucket.GetChannel(connId)
	if channel == nil {
		return
	}
	app.bucket.RemoveChannel(channel)
}

func (app *GateApp) GetChannel() {}

func (app *GateApp) PushUser() {

}

func (app *GateApp) PushGroup() {}

func (app *GateApp) Broadcast() {}