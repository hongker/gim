package interfaces

import (
	"context"
	"gim/api/server"
	"gim/internal/gate/applications"
	"gim/internal/gate/infrastructure/config"
	"gim/pkg/grpc"
)

type GRPCServer struct {
	instance *grpc.Server
	messageApp applications.MessageApp
}

func NewGRPCServer(conf *config.Config) *GRPCServer {
	grpcServer := grpc.NewServer(conf.RPC)


	return &GRPCServer{instance: grpcServer}
}

func (s *GRPCServer) Start() error {
	server.RegisterGateServer(s.instance.Instance(), s)
	return s.instance.Start()
}


func (srv *GRPCServer) PushMsg(ctx context.Context, request *server.PushMsgRequest) (res *server.PushMsgResponse, err error) {
	packet, err := request.Proto.Pack()
	if err != nil {
		return
	}
	for _, key := range request.Keys {
		srv.messageApp.Push(key, packet)
	}
	return
}

func (srv *GRPCServer) BroadcastRoom(ctx context.Context, request *server.BroadcastRoomRequest) (res *server.BroadcastRoomResponse, err error) {
	packet, err := request.Proto.Pack()
	if err != nil {
		return
	}

	err = srv.messageApp.PushRoom(request.RoomID, packet)
	return
}
