package presentation

import (
	"gim/api"
	"gim/internal/domain/event"
	"gim/internal/domain/types"
	"gim/internal/infrastructure/config"
	"gim/internal/presentation/handler"
	"gim/internal/presentation/helper"
	"gim/pkg/errors"
	"gim/pkg/network"
	"log"
)

type Handler func(ctx * network.Context)

type Socket struct {
	conf *config.Config
	handlers map[int32]Handler
	collection *types.Collection
}


func (s *Socket) RegisterHandler(operate int32, handler func(ctx *network.Context) (interface{}, error))  {
	s.handlers[operate] = s.wrapHandler(handler)
}

func (s *Socket) Start() error {
	server := network.NewTCPServer(s.conf.Addr(), network.WithPacketLength(api.PacketOffset))
	server.SetOnConnect(s.onConnect)
	server.SetOnDisconnect(s.onDisconnect)
	server.SetOnRequest(s.onRequest)
	server.Use(s.recover, s.unpack, s.validateUser)

	return server.Start()
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

	user := s.collection.GetUser(ctx.Connection())
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


func NewSocket(conf *config.Config,
	userHandler *handler.UserHandler,
	messageHandler *handler.MessageHandler,
	groupHandler *handler.GroupHandler) *Socket {
	s := &Socket{
		conf: conf,
		handlers: make(map[int32]Handler, 16),
		collection: types.NewCollection(),
	}

	s.RegisterHandler(api.OperateAuth, userHandler.Login)
	s.RegisterHandler(api.OperateMessageSend, messageHandler.Send)
	s.RegisterHandler(api.OperateMessageQuery, messageHandler.Query)
	s.RegisterHandler(api.OperateGroupJoin, groupHandler.Join)
	s.RegisterHandler(api.OperateGroupLeave, groupHandler.Leave)

	handler.NewEventHandler(s.collection, conf.Server.HeartbeatInterval).RegisterEvents()
	return s
}

