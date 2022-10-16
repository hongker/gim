package application

import (
	"context"
	"gim/framework"
	"gim/internal/domain/repository"
	"github.com/ebar-go/ego/errors"
	"sync"
)

type CometApplication interface {
	SetUserConnection(uid string, conn *framework.Connection)
	GetUserConnection(uid string) (*framework.Connection, error)
	RemoveUserConnection(uid string)
	PushUserMessage(uid string, msg []byte) error
	PushChatroomMessage(roomId string, msg []byte) error
}

type cometApplication struct {
	mu           sync.RWMutex
	connections  map[string]*framework.Connection
	chatroomRepo repository.ChatroomRepository
}

func (app *cometApplication) SetUserConnection(uid string, conn *framework.Connection) {
	app.mu.Lock()
	app.connections[uid] = conn
	app.mu.Unlock()
}

func (app *cometApplication) GetUserConnection(uid string) (*framework.Connection, error) {
	app.mu.RLock()
	conn := app.connections[uid]
	app.mu.RUnlock()
	if conn == nil {
		return nil, errors.NotFound("user not connected")
	}
	return conn, nil
}

func (app *cometApplication) RemoveUserConnection(uid string) {
	if uid == "" {
		return
	}
	app.mu.Lock()
	delete(app.connections, uid)
	app.mu.Unlock()
}

func (app *cometApplication) PushUserMessage(uid string, msg []byte) error {
	conn, err := app.GetUserConnection(uid)
	if err != nil {
		return err
	}
	conn.Push(msg)
	return nil
}

func (app *cometApplication) PushChatroomMessage(roomId string, msg []byte) error {
	ctx := context.Background()
	chatroom, err := app.chatroomRepo.Find(ctx, roomId)
	if err != nil {
		return errors.WithMessage(err, "find chatroom")
	}
	members, err := app.chatroomRepo.GetMember(ctx, chatroom)
	if err != nil {
		return errors.WithMessage(err, "get chatroom members")
	}

	for _, uid := range members {
		// only push online members.
		conn, err := app.GetUserConnection(uid)
		if err != nil {
			continue
		}
		conn.Push(msg)
	}
	return nil
}

var cometApplicationOnce struct {
	once     sync.Once
	instance *cometApplication
}

func GetCometApplication() CometApplication {
	cometApplicationOnce.once.Do(func() {
		cometApplicationOnce.instance = &cometApplication{
			connections:  map[string]*framework.Connection{},
			chatroomRepo: repository.NewChatroomRepository(),
		}
	})
	return cometApplicationOnce.instance
}
