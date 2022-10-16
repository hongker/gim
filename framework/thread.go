package framework

import "sync"

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

func NewThread(reactor *Reactor) *Thread {
	return &Thread{
		core:  reactor,
		queue: make(chan int, 100),

		connections: map[int]*Connection{},
	}
}
