package grpc

import (
	"gim/api/server"
	config2 "gim/internal/gate/infrastructure/config"
	"gim/pkg/grpc"
)

func NewLogicClient(conf *config2.Config) (server.LogicClient, error) {
	conn, err := grpc.NewGrpcConn(conf.LogicServer)
	if err != nil {
		return nil, err
	}

	return server.NewLogicClient(conn), nil
}
