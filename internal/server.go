package internal

import (
	"gim/internal/controller"
	"github.com/ebar-go/ego/component"
	"github.com/ebar-go/ego/utils/runtime"
	"github.com/ebar-go/ego/utils/runtime/signal"
	"sync"
	"time"
)

// Server represents a controller server
type Server struct {
	once   sync.Once
	config *Config

	controllers       []controller.Controller
	gatewayController controller.Controller
	apiController     controller.Controller
	jobController     controller.Controller
}

// Run runs the server
func (srv *Server) Run() {
	// run one times.
	srv.once.Do(srv.initialize)

	srv.run(signal.SetupSignalHandler())
}

// initialize init controllers.
func (srv *Server) initialize() {
	srv.gatewayController = srv.config.GatewayControllerConfig.New()
	srv.apiController = srv.config.ApiControllerConfig.New()
	srv.jobController = srv.config.JobControllerConfig.New()
}

// run start controller async.
func (srv *Server) run(stopCh <-chan struct{}) {
	watch := runtime.NewWatcher(
		controller.NewDaemonController(srv.gatewayController).NonBlockingRun(),
		controller.NewDaemonController(srv.apiController).NonBlockingRun(),
		controller.NewDaemonController(srv.jobController).NonBlockingRun(),
	)

	runtime.WaitClose(stopCh, watch.Stop, srv.shutdown)
}

// shutdown stops the server.
func (srv *Server) shutdown() {
	time.Sleep(time.Second)
	component.Provider().Logger().Info("server shutdown completed")
}
