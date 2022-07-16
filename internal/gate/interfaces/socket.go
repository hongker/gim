package interfaces

import (
	"context"
	"gim/api/client"
	"gim/api/protocol"
	"gim/internal/gate/applications"
	"gim/internal/gate/infrastructure/config"
	"gim/pkg/network"
	"log"
)

type Socket struct {
	server network.Server
	messageApp applications.MessageApp
}

func NewSocket(conf *config.Config, messageApp applications.MessageApp) *Socket {
	s := &Socket{messageApp: messageApp}
	tcpServer := network.NewTCPServer([]string{conf.TcpServer}, network.WithPacketLength(protocol.PacketOffset))
	tcpServer.Use(s.HandleAuth)
	tcpServer.SetOnConnect(s.HandleConnect)
	tcpServer.SetOnDisconnect(s.HandleDisconnect)
	tcpServer.SetOnRequest(s.HandleRequest)
	s.server = tcpServer

	return s
}

func (s *Socket) Start() error {
	return s.server.Start()
}

func (s *Socket) HandleConnect(conn *network.Connection) {
	log.Println("welcome:", conn.IP())
}

func (s *Socket) HandleAuth(ctx *network.Context) {
	// 判断是否已认证过
	ch := s.messageApp.GetChannel(ctx.Connection().ID())
	if ch != nil {
		ctx.Next()
		return
	}

	// 解析数据包
	proto := new(protocol.Proto)
	if err := proto.Unpack(ctx.Request().Body()); err != nil {
		ctx.Output(proto.MustPackFromError(protocol.InvalidParameter, err))
		ctx.Abort()
		return
	}

	uid, err := s.messageApp.Auth(context.Background(), proto, ctx.Connection())
	if err != nil {
		ctx.Output(proto.MustPackFromError(protocol.AuthFailed, err))
	} else {
		ctx.Output(proto.MustPackSuccess(&client.AuthResponse{Uid: uid}))
	}
	ctx.Abort()
	return
}

func (s *Socket) HandleRequest(ctx *network.Context) {
	// 解析数据包
	proto := new(protocol.Proto)
	if err := proto.Unpack(ctx.Request().Body()); err != nil {
		ctx.Output(proto.MustPackFromError(protocol.InvalidParameter, err))
		return
	}

	// 处理其他请求
	response, err := s.messageApp.Receive(context.Background(), proto, ctx.Connection())
	if err != nil {
		ctx.Output(proto.MustPackFromError(protocol.InvalidParameter, err))
		return
	}

	ctx.Output(proto.MustPackSuccessFromBytes(response.Data))
}

func (s *Socket) HandleDisconnect(conn *network.Connection) {
	log.Println("goodbye")
}
