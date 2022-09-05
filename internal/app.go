package internal

import (
	"gim/pkg/runtime/signal"
	"log"
)

type Server struct {
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
