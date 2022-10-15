package framework

import (
	"context"
	"gim/framework/poller"
	"gim/pkg/pool"
	"log"
)

type Reactor struct {
	poll   poller.Poller
	thread *Thread
}

func (reactor *Reactor) Run(stopCh <-chan struct{}) {
	ctx := context.Background()
	threadCtx, threadCancel := context.WithCancel(ctx)
	defer threadCancel()
	go reactor.thread.Polling(threadCtx.Done(), reactor.buildContext)

	reactor.run(stopCh)
}

func (reactor *Reactor) run(stopCh <-chan struct{}) {

	for {
		select {
		case <-stopCh:
			return
		default:
			active, err := reactor.poll.Wait()
			if err != nil {
				log.Println("unable to get active socket connection from epoll:", err)
				continue
			}

			// 处理待读取数据的链接
			for _, fd := range active {
				reactor.thread.Offer(fd)
			}
		}

	}
}

func (reactor *Reactor) buildContext(conn *Connection) (ctx *Context, err error) { return }

func NewReactor() (*Reactor, error) {
	poll, err := poller.NewPollerWithBuffer(100)
	if err != nil {
		return nil, err
	}
	return &Reactor{
		poll: poll,
		thread: &Thread{
			queue:  make(chan int, 10),
			worker: pool.NewWorkerPool(1000),
		},
	}, nil
}

type ThreadProvider struct {
}

type Thread struct {
	queue  chan int
	worker pool.Worker
}

func (thread *Thread) Add(conn *Connection)    {}
func (thread *Thread) Remove(conn *Connection) {}
func (thread *Thread) Get(fd int) *Connection  { return nil }
func (thread *Thread) Offer(fd int)            {}
func (thread *Thread) Polling(stopCh <-chan struct{}, builder func(conn *Connection) (*Context, error)) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go thread.run(ctx.Done(), builder)
	<-stopCh
}
func (thread *Thread) run(stopCh <-chan struct{}, builder func(conn *Connection) (*Context, error)) {
	var (
		ctx *Context
		err error
	)

	for {
		select {
		case <-stopCh:
			return
		default:
			fd := <-thread.queue
			conn := thread.Get(fd)
			if conn == nil {
				continue
			}

			ctx, err = builder(conn)
			if err != nil {
				_ = conn.Close()
				continue
			}
			// 读取数据不能放在协程里执行
			thread.worker.Schedule(ctx.Run)
		}

	}
}
