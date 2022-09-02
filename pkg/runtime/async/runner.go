package async

import "sync"

type Runner struct {
	lock          sync.Mutex
	loopFunctions []func(stop chan struct{})
	stop          *chan struct{}
}

// Start begins running.
func (r *Runner) Start() {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.stop == nil {
		c := make(chan struct{})
		r.stop = &c
		for i := range r.loopFunctions {
			go r.loopFunctions[i](*r.stop)
		}
	}
}

// Stop stops running.
func (r *Runner) Stop() {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.stop != nil {
		close(*r.stop)
		r.stop = nil
	}
}

// NewRunner makes a runner for the given function(s). The function(s) should loop until
// the channel is closed.
func NewRunner(f ...func(stop chan struct{})) *Runner {
	return &Runner{loopFunctions: f}
}
