package interfaces

import (
	"gim/api"
	"gim/api/protocol"
	"gim/internal/applications"
	"gim/pkg/errors"
	"gim/pkg/network"
	"log"
	"time"
)

type Handler func(ctx * network.Context,p *api.Packet) error

type Socket struct {
	server network.Server
	handlers map[int32]Handler
	userApp *applications.UserApp
	gateApp *applications.GateApp
	messageApp *applications.MessageApp
	expired time.Duration
}


func (s *Socket) OnConnect(conn *network.Connection) {
	log.Println("connect:", conn.IP())

	// 如果用户未按时登录，通过定时任务关闭连接，释放资源
	time.AfterFunc(s.expired, func() {
		uid := s.gateApp.GetUser(conn)
		if uid != "" {
			return
		}
		conn.Close()
	})
}


func (s *Socket) OnDisconnect(conn *network.Connection) {
	log.Println("disconnect:", conn.IP())
	s.gateApp.RemoveConn(conn)
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

	if err := handler(ctx, packet); err != nil {
		Failure(ctx, errors.Convert(err))
		return
	}

	packet.Op += 1
	Success(ctx, packet.Encode())

}



func (s *Socket) Start(bind string) error {
	tcpServer := network.NewTCPServer([]string{bind}, network.WithPacketLength(protocol.PacketOffset))
	tcpServer.SetOnConnect(s.OnConnect)
	tcpServer.SetOnDisconnect(s.OnDisconnect)
	tcpServer.SetOnRequest(s.OnRequest)
	s.server = tcpServer

	return s.server.Start()
}

func NewSocket(userApp *applications.UserApp,
	gateApp *applications.GateApp,
	messageApp *applications.MessageApp) *Socket {
	s := &Socket{handlers: make(map[int32]Handler, 16)}
	s.expired = time.Minute

	s.userApp = userApp
	s.gateApp = gateApp
	s.messageApp = messageApp

	s.handlers[api.OperateAuth] = s.login
	s.handlers[api.OperateMessageSend] = s.send
	s.handlers[api.OperateMessageQuery] = s.query
	return s
}