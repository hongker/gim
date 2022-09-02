package internal

import (
	"gim/pkg/runtime"
	"log"
)

type Server struct {
}

func (s *Server) Run() {
	stopCh := runtime.SetupSignalHandler()
	s.run(stopCh)
	s.Stop()
}

func (s *Server) Stop() {
	log.Println("server stopped")
}

func (s *Server) run(stopCh <-chan struct{}) {
	<-stopCh
}
