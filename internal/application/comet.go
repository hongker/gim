package application

import (
	"github.com/ebar-go/ego/errors"
	"github.com/ebar-go/ego/server/ws"
	"sync"
)

type CometApplication interface {
	SetUserConnection(uid string, conn ws.Conn)
	GetUserConnection(uid string) (ws.Conn, error)
}

type cometApplication struct {
	mu          sync.RWMutex
	connections map[string]ws.Conn
}

func (app *cometApplication) SetUserConnection(uid string, conn ws.Conn) {
	app.mu.Lock()
	app.connections[uid] = conn
	app.mu.Unlock()
}

func (app *cometApplication) GetUserConnection(uid string) (ws.Conn, error) {
	app.mu.RLock()
	conn := app.connections[uid]
	app.mu.RUnlock()
	if conn == nil {
		return nil, errors.NotFound("user not connected")
	}
	return conn, nil
}

func (app *cometApplication) JoinChatroom(roomId string, conn ws.Conn) {

}
func (app *cometApplication) LeaveChatroom(roomId string, conn ws.Conn) {}

var cometApplicationOnce struct {
	once     sync.Once
	instance *cometApplication
}

func GetCometApplication() CometApplication {
	cometApplicationOnce.once.Do(func() {
		cometApplicationOnce.instance = &cometApplication{connections: map[string]ws.Conn{}}
	})
	return cometApplicationOnce.instance
}
