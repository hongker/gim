package bucket

import (
	"gim/framework"
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
	ID         string
	Connection *framework.Connection
}

func NewSession(id string, connection *framework.Connection) *Session {
	return &Session{ID: id, Connection: connection}
}

func (session *Session) Send(msg []byte) {
	//log.Println("session send msg: ", session.ID, string(msg))
	if session.Connection == nil {
		return
	}
	session.Connection.Push(msg)
}
