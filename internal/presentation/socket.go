package presentation

import (
	"gim/api"
	"gim/internal/domain/event"
	"gim/internal/domain/types"
	"gim/internal/infrastructure/config"
	"gim/internal/presentation/filter"
	"gim/internal/presentation/handler"
	"gim/internal/presentation/helper"
	"gim/pkg/errors"
	"gim/pkg/network"
	"log"
	"time"
)

type Handler func(ctx *network.Context)

type Socket struct {
	addr     string
	handlers map[int32]Handler
}

func (s *Socket) RegisterHandler(operate int32, handler func(ctx *network.Context) (interface{}, error)) {
	s.handlers[operate] = s.wrapHandler(handler)
}

func (s *Socket) Start() error {
	server := network.NewTCPServer(s.addr, network.WithPacketLength(api.PacketOffset))
	server.SetOnConnect(s.onConnect)
	server.SetOnDisconnect(s.onDisconnect)
	server.SetOnRequest(s.onRequest)

	server.Use(filter.Recover, filter.Unpack, filter.Auth)

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

func (s *Socket) wrapHandler(fn func(ctx *network.Context) (interface{}, error)) Handler {
	return func(ctx *network.Context) {
		if response, err := fn(ctx); err != nil {
			helper.Failure(ctx, err)
		} else {
			helper.Success(ctx, response)
		}
	}
}

func NewSocket(conf *config.Config,
	userHandler *handler.UserHandler,
	messageHandler *handler.MessageHandler,
	groupHandler *handler.GroupHandler) *Socket {
	s := &Socket{
		addr:     conf.Addr(),
		handlers: make(map[int32]Handler, 16),
	}

	s.RegisterHandler(api.OperateAuth, userHandler.Login)
	s.RegisterHandler(api.OperateHeartbeat, userHandler.Heartbeat)
	s.RegisterHandler(api.OperateMessageSend, messageHandler.Send)
	s.RegisterHandler(api.OperateMessageQuery, messageHandler.Query)
	s.RegisterHandler(api.OperateGroupJoin, groupHandler.Join)
	s.RegisterHandler(api.OperateGroupLeave, groupHandler.Leave)
	s.RegisterHandler(api.OperateGroupMember, groupHandler.QueryMember)

	s.registerEvents(conf.Server.HeartbeatInterval)
	return s
}

func (s *Socket) registerEvents(expired time.Duration) {
	h := handler.NewEventHandler(types.GetCollection(), expired)
	event.Listen(event.Connect, h.Connect)
	event.Listen(event.Heartbeat, h.Heartbeat)
	event.Listen(event.Login, h.Login)
	event.Listen(event.Disconnect, h.Disconnect)
	event.Listen(event.JoinGroup, h.JoinGroup)
	event.Listen(event.LeaveGroup, h.LeaveGroup)
	event.Listen(event.Push, h.Push)
}
