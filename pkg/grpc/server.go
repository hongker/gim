package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"log"
	"net"
	"time"
)

type Server struct {
	instance *grpc.Server
	conf     ServerConfig
}

func (server *Server) Instance() *grpc.Server {
	return server.instance
}

func NewServer(conf ServerConfig) *Server {
	srv := &Server{conf: conf}
	srv.init()
	return srv
}

func (server *Server) init() {
	keepParams := grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle:     time.Duration(server.conf.IdleTimeout),
		MaxConnectionAgeGrace: time.Duration(server.conf.ForceCloseWait),
		Time:                  time.Duration(server.conf.KeepAliveInterval),
		Timeout:               time.Duration(server.conf.KeepAliveTimeout),
		MaxConnectionAge:      time.Duration(server.conf.MaxLifeTime),
	})
	//opt := grpc.ChainUnaryInterceptor(grpcProto.ServerRecoverInterceptor, grpcProto.ServerErrorConvertInterceptor, grpcProto.ServerTracerInterceptor())
	server.instance = grpc.NewServer(keepParams, grpc.ChainUnaryInterceptor())
}

func (server *Server) Start() error {
	lis, err := net.Listen(server.conf.Network, server.conf.Addr)
	if err != nil {
		return err
	}

	log.Println("start grpc listen:", server.conf.Addr)
	go func() {
		if err := server.instance.Serve(lis); err != nil {
			panic(err)
		}
	}()
	return nil
}
