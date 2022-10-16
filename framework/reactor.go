package framework

import (
	"context"
	"gim/framework/poller"
	"gim/pkg/pool"
	"github.com/ebar-go/ego/utils/runtime"
	"log"
	"sync"
)

// Reactor represents the epoll model for processing action connections.
type Reactor struct {
	poll   poller.Poller
	thread *Thread
	engine *Engine
	worker pool.Worker
}

// Run runs the Reactor with the given signal.
func (reactor *Reactor) Run(stopCh <-chan struct{}) {
	ctx := context.Background()

	threadCtx, threadCancel := context.WithCancel(ctx)
	// cancel context when the given signal is closed
	defer threadCancel()
	go func() {
		runtime.HandleCrash()
		// start thead polling task with active connection handler
		reactor.thread.Polling(threadCtx.Done(), reactor.handleActiveConnection)
	}()

	reactor.run(stopCh)
}

func (reactor *Reactor) handleActiveConnection(active int) {
	// receive an active connection
	conn := reactor.thread.GetConnection(active)
	if conn == nil {
		return
	}

	// read message
	msg, err := conn.readLine(512)
	if err != nil {
		conn.Close()
		return
	}

	// prepare Context
	ctx := reactor.engine.AcquireContext()
	ctx.reset(conn, msg)

	// process request
	reactor.worker.Schedule(ctx.Run)
}

func (reactor *Reactor) run(stopCh <-chan struct{}) {
	for {
		select {
		case <-stopCh:
			return
		default:
			// get the active connections
			active, err := reactor.poll.Wait()
			if err != nil {
				log.Println("unable to get active socket connection from epoll:", err)
				continue
			}

			// push the active connections to queue
			reactor.thread.Offer(active...)
		}

	}
}

func NewReactor() (*Reactor, error) {
	poll, err := poller.NewPollerWithBuffer(100)
	if err != nil {
		return nil, err
	}
	reactor := &Reactor{
		poll:   poll,
		engine: NewEngine(),
		worker: pool.NewWorkerPool(1000),
	}

	reactor.thread = &Thread{
		core:  reactor,
		queue: make(chan int, 100),

		connections: map[int]*Connection{},
	}
	return reactor, nil
}

type ThreadProvider struct {
}

// Thread represents sub reactor
type Thread struct {
	core  *Reactor
	queue chan int

	rmu         sync.RWMutex
	connections map[int]*Connection
}

// RegisterConnection registers a new connection to the epoll listener
func (thread *Thread) RegisterConnection(conn *Connection) error {
	fd := conn.FD()
	if err := thread.core.poll.Add(fd); err != nil {
		return err
	}

	thread.rmu.Lock()
	thread.connections[fd] = conn
	thread.rmu.Unlock()
	return nil
}

// UnregisterConnection removes the connection from the epoll listener
func (thread *Thread) UnregisterConnection(conn *Connection) {
	fd := conn.FD()
	if err := thread.core.poll.Remove(fd); err != nil {
		return
	}
	thread.rmu.Lock()
	delete(thread.connections, fd)
	thread.rmu.Unlock()
}

// GetConnection returns a connection by fd
func (thread *Thread) GetConnection(fd int) *Connection {
	thread.rmu.RLock()
	conn := thread.connections[fd]
	thread.rmu.RUnlock()
	return conn
}

// Offer push the active connections fd to the queue
func (thread *Thread) Offer(fds ...int) {
	for _, fd := range fds {
		// depose fd when queue is full
		select {
		case thread.queue <- fd:
		}
	}

}

// Polling poll with callback function
func (thread *Thread) Polling(stopCh <-chan struct{}, handler func(active int)) {
	for {
		select {
		// stop when signal is closed
		case <-stopCh:
			return
		case active := <-thread.queue:
			handler(active)
		default:
		}

	}
}
