package bucket

import (
	"log"
	"sync"
)

type Room struct {
	rmw      sync.RWMutex
	sessions map[string]*Session
}

func (room *Room) Broadcast(msg []byte) {
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

func NewRoom() *Room {
	room := &Room{
		sessions: make(map[string]*Session),
	}
	return room
}

// Session represents a session for one connection
type Session struct {
	ID string
}

func NewSession(id string) *Session {
	return &Session{ID: id}
}

func (session *Session) Send(msg []byte) {
	log.Println("session send msg: ", session.ID, string(msg))
}
