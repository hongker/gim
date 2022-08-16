package internal

import (
	"gim/internal/application"
	"gim/internal/infrastructure"
	"gim/internal/infrastructure/config"
	"gim/internal/presentation"
	"gim/internal/presentation/socket"
	"gim/pkg/app"
	"gim/pkg/system"
)

type Server struct {
	conf  *config.Config
	debug bool
}

func (s *Server) WithDebug(debug bool) *Server {
	s.debug = debug
	return s
}

func (s *Server) Run() error {
	if s.debug {
		go system.ShowMemoryUsage()
	}

	container := app.Container()
	if err := container.Provide(func() *config.Config {
		return s.conf
	}); err != nil {
		return err
	}

	infrastructure.InjectStore(container, s.conf.Server.Store)
	application.Inject(container)
	presentation.Inject(container)

	if err := container.Invoke(socket.Initialize); err != nil {
		return err
	}

	return socket.Start()
}

func NewServer(conf *config.Config) *Server {
	return &Server{conf: conf}
}
