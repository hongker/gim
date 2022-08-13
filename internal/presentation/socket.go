package presentation

import (
	"gim/api"
	"gim/internal/application"
	"gim/internal/domain/event"
	"gim/internal/presentation/handler"
	"gim/internal/presentation/helper"
	"gim/pkg/errors"
	"gim/pkg/network"
	"log"
)

type Handler func(ctx * network.Context)

type Socket struct {
	handlers map[int32]Handler
	gateApp *application.GateApp
}


func (s *Socket) RegisterHandler(operate int32, handler func(ctx *network.Context) (interface{}, error))  {
	s.handlers[operate] = s.wrapHandler(handler)
}

func (s *Socket) Start(bind string) error {
	tcpServer := network.NewTCPServer(bind, network.WithPacketLength(api.PacketOffset))
	tcpServer.SetOnConnect(s.onConnect)
	tcpServer.SetOnDisconnect(s.onDisconnect)
	tcpServer.SetOnRequest(s.onRequest)
	tcpServer.Use(s.recover, s.unpack, s.validateUser)

	return tcpServer.Start()
}


func (s *Socket) onConnect(conn *network.Connection) {
	log.Println("connect:", conn.IP())
	event.Trigger(event.Connect, conn)
}


func (s *Socket) onDisconnect(conn *network.Connection) {
	log.Println("disconnect:", conn.IP())
	event.Trigger(event.Disconnect, conn)
}


func (s *Socket) onRequest(ctx *network.Context) {
	packet := helper.GetContextPacket(ctx)
	processor, ok := s.handlers[packet.Op]
	if !ok {
		helper.Failure(ctx, errors.InvalidParameter("invalid operate"))
		return
	}

	processor(ctx)

}

func (s *Socket) unpack(ctx *network.Context) {
	packet := api.NewPacket()
	if err := packet.Decode(ctx.Request().Body()); err != nil {
		helper.Failure(ctx, errors.InvalidParameter(err.Error()))
		ctx.Abort()
		return
	}

	//log.Println(packet.Op, string(packet.Data))
	helper.SetContextPacket(ctx, packet)
	ctx.Next()
}

func (s *Socket) validateUser(ctx *network.Context) {
	packet := helper.GetContextPacket(ctx)
	if packet.Op == api.OperateAuth {
		ctx.Next()
		return
	}

	user := s.gateApp.GetUser(ctx.Connection())
	if user == nil {
		helper.Failure(ctx, errors.New(errors.CodeForbidden, "auth is required"))
		ctx.Abort()
		return
	}
	helper.SetContextUser(ctx, user)
	ctx.Next()
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


func NewSocket(userHandler *handler.UserHandler, messageHandler *handler.MessageHandler,
	groupHandler *handler.GroupHandler, gateApp *application.GateApp, eventHandler *handler.EventHandler) *Socket {
	s := &Socket{
		handlers: make(map[int32]Handler, 16),
		gateApp: gateApp,
	}

	s.RegisterHandler(api.OperateAuth, userHandler.Login)
	s.RegisterHandler(api.OperateMessageSend, messageHandler.Send)
	s.RegisterHandler(api.OperateMessageQuery, messageHandler.Query)
	s.RegisterHandler(api.OperateGroupJoin, groupHandler.Join)
	s.RegisterHandler(api.OperateGroupLeave, groupHandler.Leave)

	eventHandler.RegisterEvents()
	return s
}

