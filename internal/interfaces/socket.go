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
	handlers map[int]Handler
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

	response, err := packet.Encode()
	if err != nil {
		Failure(ctx, errors.Convert(err))
	}
	Success(ctx, response)

}



func (s *Socket) Start() error {
	return s.server.Start()
}

func NewSocket(bind string) *Socket {
	s := &Socket{handlers: make(map[int]Handler, 16)}

	tcpServer := network.NewTCPServer([]string{bind}, network.WithPacketLength(protocol.PacketOffset))
	tcpServer.SetOnConnect(s.OnConnect)
	tcpServer.SetOnDisconnect(s.OnDisconnect)
	tcpServer.SetOnRequest(s.OnRequest)
	s.server = tcpServer


	s.handlers[api.OperateAuth] = s.login
	return s
}