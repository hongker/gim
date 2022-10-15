package framework

import (
	"context"
	"gim/framework/poller"
	"gim/pkg/pool"
	"log"
	"sync"
)

type Reactor struct {
	poll      poller.Poller
	thread    *Thread
	container *ContextContainer
}

func (reactor *Reactor) Run(stopCh <-chan struct{}) {
	ctx := context.Background()
	threadCtx, threadCancel := context.WithCancel(ctx)
	defer threadCancel()
	go reactor.thread.Polling(threadCtx.Done(), reactor.container.BuildContext)

	log.Println("reactor started")
	reactor.run(stopCh)
	log.Println("reactor stopped")
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

func NewReactor() (*Reactor, error) {
	poll, err := poller.NewPollerWithBuffer(100)
	if err != nil {
		return nil, err
	}
	reactor := &Reactor{
		poll: poll,

		container: NewContextContainer(),
	}

	reactor.thread = &Thread{
		core:        reactor,
		queue:       make(chan int, 10),
		worker:      pool.NewWorkerPool(1000),
		connections: map[int]*Connection{},
	}
	return reactor, nil
}

type ThreadProvider struct {
}

type Thread struct {
	core        *Reactor
	queue       chan int
	worker      pool.Worker
	rmu         sync.RWMutex
	connections map[int]*Connection
}

func (thread *Thread) Add(conn *Connection) error {
	fd := conn.FD()
	if err := thread.core.poll.Add(fd); err != nil {
		return err
	}

	thread.rmu.Lock()
	thread.connections[fd] = conn
	thread.rmu.Unlock()
	return nil
}
func (thread *Thread) Remove(conn *Connection) {
	fd := conn.FD()
	if err := thread.core.poll.Remove(fd); err != nil {
		return
	}
	thread.rmu.Lock()
	delete(thread.connections, fd)
	thread.rmu.Unlock()
}
func (thread *Thread) Get(fd int) *Connection {
	thread.rmu.RLock()
	conn := thread.connections[fd]
	thread.rmu.RUnlock()
	return conn
}
func (thread *Thread) Offer(fd int) {
	select {
	case thread.queue <- fd:
	}
}
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
		case fd := <-thread.queue:
			conn := thread.Get(fd)
			if conn == nil {
				continue
			}

			ctx, err = builder(conn)
			if err != nil {
				conn.Close()
				continue
			}
			// 读取数据不能放在协程里执行
			thread.worker.Schedule(ctx.Run)
		default:
		}

	}
}
