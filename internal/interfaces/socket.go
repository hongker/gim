package interfaces

import (
	"gim/api"
	"gim/api/protocol"
	"gim/internal/applications"
	"gim/pkg/errors"
	"gim/pkg/network"
	"log"
)

type Handler func(p *api.Packet) error

type Socket struct {
	server network.Server
	handlers map[int32]Handler
	userApp *applications.UserApp
}


func (s *Socket) OnConnect(conn *network.Connection) {
	log.Println("connect:", conn.IP())
}


func (s *Socket) OnDisconnect(conn *network.Connection) {
	log.Println("disconnect:", conn.IP())
}


func (s *Socket) OnRequest(ctx *network.Context) {
	packet := api.NewPacket()
	if err := packet.Decode(ctx.Request().Body()); err != nil {
		Failure(ctx, errors.InvalidParameter(err.Error()))
		return
	}

	log.Println(packet.Op, packet.Data)


	handler, ok := s.handlers[packet.Op]
	if !ok {
		Failure(ctx, errors.InvalidParameter("invalid operate"))
		return
	}

	if err := handler(packet); err != nil {
		Failure(ctx, errors.Convert(err))
		return
	}

	packet.Op += 1
	Success(ctx, packet.Encode())

}



func (s *Socket) Start() error {
	return s.server.Start()
}

func NewSocket(bind string) *Socket {
	s := &Socket{handlers: make(map[int32]Handler, 16)}

	tcpServer := network.NewTCPServer([]string{bind}, network.WithPacketLength(protocol.PacketOffset))
	tcpServer.SetOnConnect(s.OnConnect)
	tcpServer.SetOnDisconnect(s.OnDisconnect)
	tcpServer.SetOnRequest(s.OnRequest)
	s.server = tcpServer

	s.userApp = applications.NewUserApp()

	s.handlers[api.OperateAuth] = s.login
	return s
}