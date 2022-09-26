package aggregator

import (
	"gim/internal/controllers"
	"gim/internal/controllers/api"
	"gim/internal/controllers/job"
	"gim/internal/controllers/socket"
	"gim/pkg/watcher"
	"github.com/ebar-go/ego/component"
	"github.com/ebar-go/ego/utils/runtime"
	"sync"
)

// Aggregate represents a controller aggregator
type Aggregator struct {
	once        sync.Once
	controllers []controllers.Controller
	watcher     watcher.Interface
}

// Run runs the aggregator
func (agg *Aggregator) Run() {
	// run one times.
	agg.once.Do(func() {
		agg.initialize()

		agg.run()
	})

	runtime.Shutdown(agg.shutdown)
}

// initialize init controllers.
func (agg *Aggregator) initialize() {
	agg.controllers = append(agg.controllers,
		api.NewController().WithName("api"),
		socket.NewController().WithName("gateway"),
		job.NewController().WithName("job"),
	)
}

// run start controller async.
func (agg *Aggregator) run() {
	stopChs := make([]chan struct{}, 0)
	for _, controller := range agg.controllers {
		ch := make(chan struct{})
		stopChs = append(stopChs, ch)
		go controller.Run(ch, 1)
	}

	agg.watcher = watcher.NewChanWatcher(stopChs...)
}

// shutdown stops the aggregator.
func (agg *Aggregator) shutdown() {
	agg.watcher.Stop()
	component.Provider().Logger().Info("shutdown success")
}
