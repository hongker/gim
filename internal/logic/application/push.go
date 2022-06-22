package application

import "gim/api/protocol"

type PushApp struct {
}

func (app *PushApp) PushUser(uid string, Proto *protocol.Proto) {

}
func (app *PushApp) PushRoom(roomId string, Proto *protocol.Proto) {

}

func (app *PushApp) Broadcast(Proto *protocol.Proto) {

}
