package interfaces

import (
	"gim/api"
	"gim/internal/applications"
	"gim/internal/interfaces/handler"
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
		if !s.gateApp.CheckConnExist(conn) {
			conn.Close()
		}

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


	h, ok := s.handlers[packet.Op]
	if !ok {
		Failure(ctx, errors.InvalidParameter("invalid operate"))
		return
	}

	h(ctx, packet)

}

func (s *Socket) recover(ctx *network.Context) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("recover: err=%v\n", err)
		}
	}()
}

func (s *Socket) Start(bind string) error {
	tcpServer := network.NewTCPServer([]string{bind}, network.WithPacketLength(api.PacketOffset))
	tcpServer.SetOnConnect(s.OnConnect)
	tcpServer.SetOnDisconnect(s.OnDisconnect)
	tcpServer.Use(s.recover)
	tcpServer.SetOnRequest(s.OnRequest)

	return tcpServer.Start()
}

func (s *Socket) wrapHandler(fn func(ctx *network.Context, p *api.Packet) error) Handler{
	return func(ctx *network.Context, p *api.Packet) {
		if err := fn(ctx, p); err != nil {
			Failure(ctx, errors.Convert(err))
		}else {
			Success(ctx, p.Encode())
		}
	}
}

func NewSocket(userHandler *handler.UserHandler, messageHandler *handler.MessageHandler) *Socket {
	s := &Socket{
		handlers: make(map[int32]Handler, 16),
		expired: time.Minute,
	}

	s.handlers[api.OperateAuth] = s.wrapHandler(userHandler.Login)
	s.handlers[api.OperateMessageSend] = s.wrapHandler(messageHandler.Send)
	s.handlers[api.OperateMessageQuery] = s.wrapHandler(messageHandler.Query)
	return s
}