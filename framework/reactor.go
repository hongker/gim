package framework

import (
	"context"
	"gim/framework/poller"
	"gim/pkg/pool"
	"github.com/ebar-go/ego/utils/runtime"
	"log"
)

// Reactor represents the epoll model for processing action connections.
type Reactor struct {
	poll             poller.Poller
	thread           *Thread
	engine           *Engine
	worker           pool.Worker
	packetLengthSize int
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

// run receive active connection file descriptor and offer to thread
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

// handleActiveConnection handles active connection request
func (reactor *Reactor) handleActiveConnection(active int) {
	// receive an active connection
	conn := reactor.thread.GetConnection(active)
	if conn == nil {
		return
	}

	// read message
	msg, err := conn.readLine(reactor.packetLengthSize)
	if err != nil {
		conn.Close()
		return
	}

	// prepare Context
	ctx := reactor.engine.AcquireContext()
	ctx.reset(conn, msg)

	// process request
	reactor.worker.Schedule(func() {
		reactor.engine.HandleContext(ctx)
	})
}

// ReactorOptions represents the options for the reactor
type ReactorOptions struct {
	// EpollBufferSize is the size of the active connections in every duration
	EpollBufferSize int

	// WorkerPollSize is the size of the worker pool
	WorkerPoolSize int

	// PacketLengthSize is the size of the packet length offset
	PacketLengthSize int

	// ThreadQueueCapacity is the cap of the thread queue
	ThreadQueueCapacity int
}

func NewReactor(options ReactorOptions) (*Reactor, error) {
	poll, err := poller.NewPollerWithBuffer(options.EpollBufferSize)
	if err != nil {
		return nil, err
	}
	reactor := &Reactor{
		poll:             poll,
		engine:           NewEngine(),
		worker:           pool.NewWorkerPool(options.WorkerPoolSize),
		packetLengthSize: options.PacketLengthSize,
		thread:           NewThread(options.ThreadQueueCapacity),
	}

	return reactor, nil
}
