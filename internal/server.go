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

	gatewayController controller.Controller
	apiController     controller.Controller
	jobController     controller.Controller
}

// Run runs the server
func (srv *Server) Run() error {
	// run one times.
	srv.once.Do(srv.initialize)

	// run with signal
	return srv.run(signal.SetupSignalHandler())
}

// initialize init controllers.
func (srv *Server) initialize() {
	srv.gatewayController = srv.config.GatewayControllerConfig.New("gateway")
	srv.apiController = srv.config.ApiControllerConfig.New("api")
	srv.jobController = srv.config.JobControllerConfig.New("job")
}

// run start controller async.
func (srv *Server) run(stopCh <-chan struct{}) error {
	watch := runtime.NewWatcher(
		controller.NewDaemonController(srv.gatewayController).NonBlockingRun(),
		controller.NewDaemonController(srv.apiController).NonBlockingRun(),
		controller.NewDaemonController(srv.jobController).NonBlockingRun(),
	)

	component.Provider().Logger().Infof("server started successfully")
	runtime.WaitClose(stopCh, watch.Stop, srv.shutdown)
	return nil
}

// shutdown stops the server.
func (srv *Server) shutdown() {
	time.Sleep(time.Second)
	component.Provider().Logger().Info("server shutdown completed")
}
