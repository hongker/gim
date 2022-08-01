package interfaces

import (
	"gim/api"
	"gim/internal/applications"
	"gim/pkg/errors"
	"gim/pkg/network"
	"log"
	"time"
)

type Handler func(ctx * network.Context,p *api.Packet)

type Socket struct {
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

	handler(ctx, packet)

}



func (s *Socket) Start(bind string) error {
	tcpServer := network.NewTCPServer([]string{bind}, network.WithPacketLength(api.PacketOffset))
	tcpServer.SetOnConnect(s.OnConnect)
	tcpServer.SetOnDisconnect(s.OnDisconnect)
	tcpServer.SetOnRequest(s.OnRequest)

	return tcpServer.Start()
}

func NewSocket(userApp *applications.UserApp, gateApp *applications.GateApp,
	messageApp *applications.MessageApp,) *Socket {
	s := &Socket{
		handlers: make(map[int32]Handler, 16),
		userApp: userApp,
		gateApp: gateApp,
		messageApp: messageApp,
		expired: time.Minute,
	}

	s.handlers[api.OperateAuth] = s.WrapHandler(s.login)
	s.handlers[api.OperateMessageSend] = s.WrapHandler(s.send)
	s.handlers[api.OperateMessageQuery] = s.WrapHandler(s.query)
	return s
}