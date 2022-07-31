package interfaces

import (
	"gim/api/protocol"
	"gim/internal/domain/dto"
	"gim/pkg/network"
	"log"
)

type Socket struct {
	server network.Server
	
}


func (s *Socket) OnConnect(conn *network.Connection) {
	log.Println("connect:", conn.IP())
}


func (s *Socket) OnDisconnect(conn *network.Connection) {
	log.Println("disconnect:", conn.IP())
}


func (s *Socket) OnRequest(ctx *network.Context) {
	packet := dto.NewPacket()
	if err := packet.Decode(ctx.Request().Body()); err != nil {
		ctx.Output(proto.MustPackFromError(protocol.InvalidParameter, err))
		return
	}
}


func (s *Socket) Start() error {
	return s.server.Start()
}

func NewSocket(bind string) *Socket {
	s := &Socket{}

	tcpServer := network.NewTCPServer([]string{bind}, network.WithPacketLength(protocol.PacketOffset))
	tcpServer.SetOnConnect(s.OnConnect)
	tcpServer.SetOnDisconnect(s.OnDisconnect)
	tcpServer.SetOnRequest(s.OnRequest)
	s.server = tcpServer
	return s
}