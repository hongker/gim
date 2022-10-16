package internal

import (
	"context"
	"gim/internal/controller"
	"github.com/ebar-go/ego/component"
	"github.com/ebar-go/ego/utils/runtime"
	"github.com/ebar-go/ego/utils/runtime/signal"
	"sync"
)

// Server represents a controller server
type Server struct {
	config *Config

	once sync.Once

	controllers []controller.Controller
}

// Run runs the server
func (srv *Server) Run() error {
	// run one times.
	srv.once.Do(srv.initialize)

	controllerCtx, controllerCancel := context.WithCancel(context.Background())
	defer controllerCancel()

	for _, c := range srv.controllers {
		target := c
		go func() {
			defer runtime.HandleCrash()
			target.Run(controllerCtx.Done())
		}()
	}

	component.Provider().Logger().Infof("server started successfully")
	runtime.WaitClose(signal.SetupSignalHandler())
	return nil

}

// initialize init controllers.
func (srv *Server) initialize() {
	srv.controllers = append(srv.controllers,
		srv.config.GatewayControllerConfig.New("gateway"),
		srv.config.ApiControllerConfig.New("api"),
		srv.config.JobControllerConfig.New("job"),
	)
}
