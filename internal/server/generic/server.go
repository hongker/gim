package generic

import "gim/internal/application"

// Server represents a generic server with tcp service.
type Server struct {
	el *application.EventLoop
}

func (s *Server) Run(stopCh <-chan struct{}) error {
	defer s.Stop()
	if err := s.el.Start(); err != nil {
		return err
	}
	<-stopCh
	return nil
}

func (s *Server) Stop() {
	s.el.Stop()
}
