package application

import (
	"github.com/ebar-go/ego/errors"
	"github.com/ebar-go/ego/server/socket"
	"sync"
)

type CometApplication interface {
	SetUserConnection(uid string, conn socket.Connection)
	GetUserConnection(uid string) (socket.Connection, error)
	RemoveUserConnection(uid string)
	PushUserMessage(uid string, msg []byte) error
	PushChatroomMessage(roomId string, msg []byte) error
}

type cometApplication struct {
	mu          sync.RWMutex
	connections map[string]socket.Connection
}

func (app *cometApplication) SetUserConnection(uid string, conn socket.Connection) {
	app.mu.Lock()
	app.connections[uid] = conn
	app.mu.Unlock()
}

func (app *cometApplication) GetUserConnection(uid string) (socket.Connection, error) {
	app.mu.RLock()
	conn := app.connections[uid]
	app.mu.RUnlock()
	if conn == nil {
		return nil, errors.NotFound("user not connected")
	}
	return conn, nil
}

func (app *cometApplication) RemoveUserConnection(uid string) {
	app.mu.Lock()
	delete(app.connections, uid)
	app.mu.Unlock()
}

func (app *cometApplication) PushUserMessage(uid string, msg []byte) error {
	conn, err := app.GetUserConnection(uid)
	if err != nil {
		return err
	}
	return conn.Push(msg)
}

func (app *cometApplication) PushChatroomMessage(roomId string, msg []byte) error {
	return nil
}

func (app *cometApplication) JoinChatroom(roomId string, conn socket.Connection) {

}
func (app *cometApplication) LeaveChatroom(roomId string, conn socket.Connection) {}

var cometApplicationOnce struct {
	once     sync.Once
	instance *cometApplication
}

func GetCometApplication() CometApplication {
	cometApplicationOnce.once.Do(func() {
		cometApplicationOnce.instance = &cometApplication{connections: map[string]socket.Connection{}}
	})
	return cometApplicationOnce.instance
}
