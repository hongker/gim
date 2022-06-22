package gate

// 长链接网关，负责维护客户端的连接、数据收发
// 收到客户端请求后，通过grpc调用，将请求数据发送给logic服务，并将响应数据返回给客户端

import (
	"gim/api/protocol"
	"gim/api/server"
	"gim/pkg/grpc"
	"gim/pkg/network"
	"gim/pkg/system"
	"log"
)

func Run() {
	srv := &Server{bucket: NewBucket()}

	srv.conf = InitConfig()
	srv.initLogicClient()
	srv.initTCP()
	srv.initGRPC()

	system.Shutdown(func() {
		log.Println("server shutdown")
	})
}

type Server struct {
	conf   *Config
	bucket *Bucket
	logic  server.LogicClient
}

// initTCP 初始化tcp服务
func (srv *Server) initTCP() {
	tcpServer := network.NewTCPServer([]string{srv.conf.TcpServer}, network.WithPacketLength(protocol.PacketOffset))

	tcpServer.Use(srv.HandleAuth)
	tcpServer.SetOnConnect(srv.HandleConnect)
	tcpServer.SetOnDisconnect(srv.HandleDisconnect)
	tcpServer.SetOnRequest(srv.HandleRequest)

	system.SecurePanic(tcpServer.Start())
}

// initGRPC 初始化grpc服务
func (srv *Server) initGRPC() {
	grpcServer := grpc.NewServer(srv.conf.RPC)
	server.RegisterGateServer(grpcServer.Instance(), srv)
	system.SecurePanic(grpcServer.Start())
}

// initLogicClient 初始化logicClient
func (srv *Server) initLogicClient() {
	conn, err := grpc.NewGrpcConn(srv.conf.LogicServer)
	if err != nil {
		panic(err)
	}

	srv.logic = server.NewLogicClient(conn)
}
