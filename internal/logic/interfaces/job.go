package interfaces

import (
	"gim/api/protocol"
	"gim/internal/logic/application"
	"log"
)

type Job struct {
	pushApp   *application.PushApp
	sendQueue chan Packet
}

const (
	TypePush      = 1
	TypeRoom      = 2
	TypeBroadcast = 3
)

type Packet struct {
	Type   int
	Uid    string
	RoomId string
	Proto  *protocol.Proto
}

func NewJob() *Job {
	return &Job{}
}

func (job *Job) Start() {

}

func (job *Job) consumer() {

}

func (job *Job) dispatch() {
	for {
		select {
		default:
			packet, ok := <-job.sendQueue
			if !ok {
				return
			}
			if err := job.push(packet); err != nil {
				log.Println("push failed:", err)
			}
		}
	}
}

func (job *Job) push(packet Packet) (err error) {
	switch packet.Type {
	case TypePush:
		job.pushApp.PushUser(packet.Uid, packet.Proto)
	case TypeRoom:
		job.pushApp.PushRoom(packet.RoomId, packet.Proto)
	case TypeBroadcast:
		job.pushApp.Broadcast(packet.Proto)
	}
	return
}
