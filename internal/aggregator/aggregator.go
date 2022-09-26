package aggregator

import (
	"gim/internal/controllers"
	"gim/pkg/runtime/signal"
	"gim/pkg/watcher"
	"github.com/ebar-go/ego/component"
	"github.com/ebar-go/ego/utils/runtime"
	"sync"
	"time"
)

// Aggregator represents a controller aggregator
type Aggregator struct {
	once   sync.Once
	config *Config

	gatewayController controllers.Controller
	apiController     controllers.Controller
	jobController     controllers.Controller
}

// Run runs the aggregator
func (agg *Aggregator) Run() {
	// run one times.
	agg.once.Do(agg.initialize)

	agg.run(signal.SetupSignalHandler())
}

// initialize init controllers.
func (agg *Aggregator) initialize() {
	agg.gatewayController = agg.config.GatewayControllerConfig.New()
	agg.apiController = agg.config.ApiControllerConfig.New()
	agg.jobController = agg.config.JobControllerConfig.New()
}

// run start controller async.
func (agg *Aggregator) run(stopCh <-chan struct{}) {
	stopChs := make([]chan struct{}, 0)

	stopChs = append(stopChs,
		controllers.NewDaemonController(agg.gatewayController).NonBlockingRun(),
		controllers.NewDaemonController(agg.apiController).NonBlockingRun(),
		controllers.NewDaemonController(agg.jobController).NonBlockingRun(),
	)

	watch := watcher.NewChanWatcher(stopChs...)

	runtime.WaitClose(stopCh, watch.Stop, agg.shutdown)
}

// shutdown stops the aggregator.
func (agg *Aggregator) shutdown() {
	time.Sleep(time.Second)
	component.Provider().Logger().Info("aggregator shutdown completed")
}
