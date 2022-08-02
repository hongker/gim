package interfaces

import (
	"gim/api"
	"gim/internal/applications"
	"gim/internal/interfaces/handler"
	"gim/internal/interfaces/helper"
	"gim/pkg/errors"
	"gim/pkg/network"
	"log"
	"time"
)

type Handler func(ctx * network.Context)

type Socket struct {
	handlers map[int32]Handler
	gateApp *applications.GateApp
	expired time.Duration
}


func (s *Socket) Start(bind string) error {
	tcpServer := network.NewTCPServer([]string{bind}, network.WithPacketLength(api.PacketOffset))
	tcpServer.SetOnConnect(s.onConnect)
	tcpServer.SetOnDisconnect(s.onDisconnect)
	tcpServer.Use(s.recover,)
	tcpServer.SetOnRequest(s.onRequest)

	return tcpServer.Start()
}


func (s *Socket) onConnect(conn *network.Connection) {
	log.Println("connect:", conn.IP())

	// 如果用户未按时登录，通过定时任务关闭连接，释放资源
	time.AfterFunc(s.expired, func() {
		if !s.gateApp.CheckConnExist(conn) {
			conn.Close()
		}

	})
}


func (s *Socket) onDisconnect(conn *network.Connection) {
	log.Println("disconnect:", conn.IP())
	s.gateApp.RemoveConn(conn)
}


func (s *Socket) onRequest(ctx *network.Context) {
	packet := api.NewPacket()
	if err := packet.Decode(ctx.Request().Body()); err != nil {
		helper.Failure(ctx, errors.InvalidParameter(err.Error()))
		return
	}

	log.Println(packet.Op, packet.Data)
	helper.SetContextPacket(ctx, packet)

	processor, ok := s.handlers[packet.Op]
	if !ok {
		helper.Failure(ctx, errors.InvalidParameter("invalid operate"))
		return
	}

	processor(ctx)

}



func (s *Socket) recover(ctx *network.Context) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("recover: err=%v\n", err)
		}
	}()

	ctx.Next()
}



func (s *Socket) wrapHandler(fn func(ctx *network.Context) (interface{}, error)) Handler{
	return func(ctx *network.Context) {
		if response, err := fn(ctx); err != nil {
			helper.Failure(ctx, err)
		}else {
			helper.Success(ctx, response)
		}
	}
}

func (s *Socket) registerHandler(operate int32, handler func(ctx *network.Context) (interface{}, error))  {
	s.handlers[operate] = s.wrapHandler(handler)
}

func NewSocket(userHandler *handler.UserHandler, messageHandler *handler.MessageHandler,
	groupHandler *handler.GroupHandler,
	gateApp *applications.GateApp) *Socket {
	s := &Socket{
		handlers: make(map[int32]Handler, 16),
		expired: time.Minute,
	}

	s.gateApp = gateApp
	s.registerHandler(api.OperateAuth, userHandler.Login)
	s.registerHandler(api.OperateMessageSend, messageHandler.Send)
	s.registerHandler(api.OperateMessageQuery, messageHandler.Query)
	s.registerHandler(api.OperateGroupJoin, groupHandler.Join)
	s.registerHandler(api.OperateGroupLeave, groupHandler.Leave)
	return s
}

