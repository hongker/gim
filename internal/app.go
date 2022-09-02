package internal

import (
	"gim/pkg/runtime/signal"
	"gim/pkg/server"
	"log"
)

type Server struct {
	httpServer   *server.HttpServer
	grpcServer   *server.GrpcServer
	socketServer *server.SocketServer
}

func (s *Server) Run() {
	stopCh := signal.SetupSignalHandler()
	s.run(stopCh)
	s.Stop()
}

func (s *Server) Stop() {
	log.Println("server stopped")
}

func (s *Server) run(stopCh <-chan struct{}) {
	log.Println("server started")

	<-stopCh
}
