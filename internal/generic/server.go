package generic

// Server represents a generic server with tcp service.
type Server struct {
}

func (s *Server) Run(stopCh <-chan struct{}) {
	<-stopCh
}
