package framework

import (
	"context"
	"github.com/ebar-go/ego/utils/runtime/signal"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"log"
	"net"
	"testing"
)

func TestEngine(t *testing.T) {
	router := NewRouter().
		Handle(1, StandardHandler[LoginRequest, LoginResponse](LoginAction))

	callback := NewCallback().
		OnConnect(func(conn *Connection) {
			log.Printf("[%s] connected\n", conn.UUID())
		}).OnDisconnect(func(conn *Connection) {
		log.Printf("[%s] disconnected\n", conn.UUID())
	})

	engine := New().
		WithCallback(callback).
		WithRouter(router).Listen(TCP, ":8080")

	err := engine.Run(signal.SetupSignalHandler())
	assert.Nil(t, err)

}

type LoginRequest struct{ Name string }
type LoginResponse struct{ Token string }

func LoginAction(ctx context.Context, req *LoginRequest) (response *LoginResponse, err error) {
	response = &LoginResponse{Token: uuid.NewV4().String()}
	return
}

func TestClient(t *testing.T) {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	_, err = conn.Write([]byte("hello"))
	if err != nil {
		panic(err)
	}

	log.Println("send success")

	receive := make([]byte, 512)
	n, err := conn.Read(receive)
	if err != nil {
		panic(err)
	}

	log.Println(string(receive[:n]))
}
