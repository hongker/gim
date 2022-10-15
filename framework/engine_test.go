package framework

import (
	"context"
	"github.com/ebar-go/ego/utils/runtime/signal"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEngine(t *testing.T) {
	callback := NewCallback().
		OnRequest(func(conn *Connection) {

		}).OnConnect(func(conn *Connection) {

	}).OnDisconnect(func(conn *Connection) {})

	router := NewRouter().
		Handle(1, StandardHandler[LoginRequest, LoginResponse](LoginAction))

	engine := New().
		WithCallback(callback).
		WithRouter(router).
		WithCodec(NewJsonCodec()).Listen(TCP, ":8080")

	err := engine.Run(signal.SetupSignalHandler())
	assert.Nil(t, err)

}

type LoginRequest struct{ Name string }
type LoginResponse struct{ Token string }

func LoginAction(ctx context.Context, req *LoginRequest) (response *LoginResponse, err error) {
	response = &LoginResponse{Token: uuid.NewV4().String()}
	return
}
