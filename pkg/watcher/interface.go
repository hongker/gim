package watcher

import "context"

type Interface interface {
	Stop()
}
type ChanWatcher struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func (w *ChanWatcher) Stop() {
	w.cancel()
}

func (w *ChanWatcher) Watch(fn func(ch chan struct{})) {

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
