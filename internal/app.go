package internal

import (
	"gim/internal/extension"
	"gim/internal/generic"
	"log"
)

type Server struct {
	genericServer   *generic.Server
	extensionServer *extension.Server
}

func (s *Server) Run(stopCh <-chan struct{}) {
	log.Println("server started")

	defer s.Stop()
	s.run(stopCh)

}

func (s *Server) Stop() {
	log.Println("server stopped")
}

func (s *Server) run(stopCh <-chan struct{}) {
	go s.genericServer.Run(stopCh)
	go s.extensionServer.Run(stopCh)
}
