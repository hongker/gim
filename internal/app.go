package internal

import (
	"github.com/ebar-go/ego"
	"log"
)

type Server struct {
	engine *ego.NamedEngine
}

func NewServer() *Server {
	return &Server{
		engine: ego.New(),
	}
}

func (s *Server) Run(stopCh <-chan struct{}) {
	defer s.Stop()
	log.Println("server started")

	s.engine.Run()
	<-stopCh
}

func (s *Server) Stop() {
	log.Println("server stopped")
}
