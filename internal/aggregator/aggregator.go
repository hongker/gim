package aggregator

import (
	"context"
	"gim/internal/controllers"
	"gim/internal/controllers/api"
	"gim/internal/controllers/job"
	"gim/internal/controllers/socket"
	"github.com/ebar-go/ego/component"
	"github.com/ebar-go/ego/utils/runtime"
	"sync"
)

type Aggregator struct {
	once        sync.Once
	controllers []controllers.Controller
	watcher     Watcher
}

func (agg *Aggregator) Run() {
	agg.once.Do(agg.initialize)

	agg.run()

	runtime.Shutdown(agg.shutdown)
}

func (agg *Aggregator) initialize() {
	agg.controllers = append(agg.controllers,
		api.NewController().WithName("api"),
		socket.NewController().WithName("gateway"),
		job.NewController().WithName("job"),
	)
}
func (agg *Aggregator) run() {
	stopChs := make([]chan struct{}, 0)
	for _, controller := range agg.controllers {
		ch := make(chan struct{})
		stopChs = append(stopChs, ch)
		go controller.Run(ch, 1)
	}

	agg.watcher = NewChanWatcher(stopChs...)
}
func (agg *Aggregator) shutdown() {
	agg.watcher.Stop()
	component.Provider().Logger().Info("shutdown success")
}

type Watcher interface {
	Stop()
}
type ChanWatcher struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func (w ChanWatcher) Stop() {
	w.cancel()
}

func NewChanWatcher(chs ...chan struct{}) *ChanWatcher {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-ctx.Done()
		for _, ch := range chs {
			close(ch)
		}
	}()
	return &ChanWatcher{
		ctx:    ctx,
		cancel: cancel,
	}
}
