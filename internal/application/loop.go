package application

import "gim/pkg/network"

type EventLoop struct {
	socketInstance network.Server
}

func (e *EventLoop) OnConnect(conn *network.Connection) {

}

func (e *EventLoop) OnDisconnect(conn *network.Connection) {}
func (e *EventLoop) OnRequest(ctx *network.Context)        {}

func (e *EventLoop) Start() error {
	e.socketInstance.SetOnConnect(e.OnConnect)
	e.socketInstance.SetOnDisconnect(e.OnDisconnect)
	e.socketInstance.SetOnRequest(e.OnRequest)

	return e.socketInstance.Start()
}

func (e *EventLoop) Stop() {
}
func BuildEventLoop(socketInstance network.Server) *EventLoop {
	return &EventLoop{socketInstance: socketInstance}
}
