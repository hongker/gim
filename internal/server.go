package internal

import (
	"gim/internal/module/gateway"
	"gim/internal/module/job"
	"gim/internal/module/message"
	"log"
)

// Server
type Server struct {
	gatewayInstance *gateway.Instance
	messageInstance *message.Instance
	jobInstance     *job.Instance
}

func NewServer(gatewayInstance *gateway.Instance) *Server {
	return &Server{gatewayInstance: gatewayInstance}
}

func (s *Server) Run(stopCh <-chan struct{}) {
	defer s.Stop()
	log.Println("server started")

	s.gatewayInstance.Start()
	<-stopCh
}

func (s *Server) Stop() {
	log.Println("server stopped")
}
