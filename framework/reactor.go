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

type ReactorOptions struct {
	EpollBufferSize  int
	WorkerPoolSize   int
	PacketLengthSize int
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
	}

	reactor.thread = NewThread(reactor)
	return reactor, nil
}
