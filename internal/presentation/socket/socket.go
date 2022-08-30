package socket

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
	"sync"
	"time"
)

type HandlerFunction func(ctx *network.Context) (interface{}, error)

type NotFoundHandlerFunction func(ctx *network.Context)
type RecoverHandlerFunction func(ctx *network.Context)
type SuccessHandlerFunction func(ctx *network.Context, response interface{})
type ServiceErrorHandleFunction func(ctx *network.Context, err error)

type Socket struct {
	handlers               map[int32]HandlerFunction
	addr                   string
	doNotRecover           bool
	notFoundHandlerFunc    NotFoundHandlerFunction
	successHandlerFunc     SuccessHandlerFunction
	serviceErrorHandleFunc ServiceErrorHandleFunction
	recoverHandleFunc      RecoverHandlerFunction
	contentEncodingEnabled bool
}

func (s *Socket) NotFoundHandlerFunc(handler NotFoundHandlerFunction) {
	if handler == nil {
		return
	}
	s.notFoundHandlerFunc = handler
}

func (s *Socket) RecoverHandler(handler RecoverHandlerFunction) {
	if handler == nil {
		return
	}
	s.recoverHandleFunc = handler
}

func (s *Socket) SuccessHandler(handler SuccessHandlerFunction) {
	if handler == nil {
		return
	}
	s.successHandlerFunc = handler
}

func (s *Socket) ServiceErrorHandler(handler ServiceErrorHandleFunction) {
	if handler == nil {
		return
	}
	s.serviceErrorHandleFunc = handler
}

func (s *Socket) Start() error {
	server := network.NewTCPServer(s.addr, network.WithPacketLength(api.PacketOffset))
	server.SetOnConnect(s.onConnect)
	server.SetOnDisconnect(s.onDisconnect)
	server.SetOnRequest(s.onRequest)

	if !s.doNotRecover {
		server.Use(func(ctx *network.Context) {
			s.recoverHandleFunc(ctx)
		})
	}
	server.Use(filter.Unpack, filter.Auth)

	return server.Start()
}

//--------------------private methods------------------------

func (s *Socket) registerHandler(operate int32, handler HandlerFunction) {
	s.handlers[operate] = handler
}

func (s *Socket) onConnect(conn *network.Connection) {
	log.Println("connect:", conn.IP())
	event.Trigger(event.Connect, &event.ConnectEvent{conn})
}

func (s *Socket) onDisconnect(conn *network.Connection) {
	log.Println("disconnect:", conn.IP())
	event.Trigger(event.Disconnect, &event.DisconnectEvent{conn})
}

func (s *Socket) onRequest(ctx *network.Context) {
	packet := helper.GetContextPacket(ctx)
	processor, ok := s.handlers[packet.Op]
	if !ok {
		s.notFoundHandlerFunc(ctx)
		return
	}

	response, err := processor(ctx)
	if err != nil {
		s.serviceErrorHandleFunc(ctx, err)
	} else {
		s.successHandlerFunc(ctx, response)
	}

}

func (s *Socket) registerEvents(expired time.Duration) {
	h := handler.NewEventHandler(types.GetCollection(), expired)
	event.NewEvent[*event.ConnectEvent](event.Connect).Bind(h.Connect)
	event.NewEvent[*event.HeartbeatEvent](event.Heartbeat).Bind(h.Heartbeat)
	event.NewEvent[*event.DisconnectEvent](event.Disconnect).Bind(h.Disconnect)
	event.NewEvent[*event.JoinGroupEvent](event.JoinGroup).Bind(h.JoinGroup)
	event.NewEvent[*event.LeaveGroupEvent](event.LeaveGroup).Bind(h.LeaveGroup)
	event.NewEvent[*event.PushMessageEvent](event.Push).Bind(h.Push)
	event.NewEvent[*event.LoginEvent](event.Login).Bind(h.Login)
}

func buildSocket(conf *config.Config) *Socket {
	s := &Socket{
		addr:                   conf.Addr(),
		handlers:               make(map[int32]HandlerFunction, 16),
		serviceErrorHandleFunc: helper.Failure,
		successHandlerFunc:     helper.Success,
		notFoundHandlerFunc: func(ctx *network.Context) {
			helper.Failure(ctx, errors.InvalidParameter("invalid operate"))
		},
		recoverHandleFunc:      filter.Recover,
		doNotRecover:           conf.Debug,
		contentEncodingEnabled: false,
	}

	return s
}

var socketInstance struct {
	once   sync.Once
	socket *Socket
}

func Initialize(conf *config.Config,
	userHandler *handler.UserHandler,
	messageHandler *handler.MessageHandler,
	groupHandler *handler.GroupHandler) {
	socketInstance.once.Do(func() {
		s := buildSocket(conf)

		s.registerHandler(api.OperateAuth, userHandler.Login)
		s.registerHandler(api.OperateHeartbeat, userHandler.Heartbeat)
		s.registerHandler(api.OperateMessageSend, messageHandler.Send)
		s.registerHandler(api.OperateMessageQuery, messageHandler.Query)
		s.registerHandler(api.OperateGroupJoin, groupHandler.Join)
		s.registerHandler(api.OperateGroupLeave, groupHandler.Leave)
		s.registerHandler(api.OperateGroupMember, groupHandler.QueryMember)

		s.registerEvents(conf.Server.HeartbeatInterval)
		socketInstance.socket = s
	})
}

func Get() *Socket {
	return socketInstance.socket
}
