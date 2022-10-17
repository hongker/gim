package bucket

import (
	"github.com/ebar-go/ego/utils/runtime"
	"sync"
)

type Room struct {
	rmw      sync.RWMutex
	sessions map[string]*Session
	queue    chan []byte
	once     sync.Once
	done     chan struct{}
}

func (room *Room) Broadcast(msg []byte) {
	select {
	case room.queue <- msg:
	}
}

func (room *Room) start(goroutineSize int) {
	for i := 0; i < goroutineSize; i++ {
		go func() {
			defer runtime.HandleCrash()
			room.polling(room.done)
		}()
	}
}
func (room *Room) stop() {
	room.once.Do(func() {
		close(room.done)
	})
}

func (room *Room) polling(done <-chan struct{}) {
	for {
		select {
		case <-done:
			return
		case msg, ok := <-room.queue:
			if !ok {
				return
			}

			room.broadcast(msg)
		}
	}
}

func (room *Room) broadcast(msg []byte) {
	for _, session := range room.sessions {
		session.Send(msg)
	}
}

func (room *Room) AddSession(session *Session) {
	room.rmw.Lock()
	room.sessions[session.ID] = session
	room.rmw.Unlock()
}
func (room *Room) RemoveSession(session *Session) {
	room.rmw.Lock()
	delete(room.sessions, session.ID)
	room.rmw.Unlock()
}
func (room *Room) GetSession(id string) *Session {
	room.rmw.RLock()
	session := room.sessions[id]
	room.rmw.RUnlock()
	return session
}

func NewRoom(queueSize int, goroutineSize int) *Room {
	room := &Room{
		sessions: make(map[string]*Session),
		queue:    make(chan []byte, queueSize),
		done:     make(chan struct{}),
	}
	room.start(goroutineSize)

	return room
}
