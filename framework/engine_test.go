package framework

import (
	"context"
	uuid "github.com/satori/go.uuid"
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
		WithCodec(NewJsonCodec())

	engine.Run()
	defer engine.Close()
}

type LoginRequest struct{ Name string }
type LoginResponse struct{ Token string }

func LoginAction(ctx context.Context, req *LoginRequest) (response *LoginResponse, err error) {
	response = &LoginResponse{Token: uuid.NewV4().String()}
	return
}
