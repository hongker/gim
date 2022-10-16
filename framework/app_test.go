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

func TestApp(t *testing.T) {
	app := New()

	app.Callback().
		OnConnect(func(conn *Connection) {
			log.Printf("[%s] connected\n", conn.UUID())
		}).OnDisconnect(func(conn *Connection) {
		log.Printf("[%s] disconnected\n", conn.UUID())
	})

	app.Router().Route(1, StandardHandler[LoginRequest, LoginResponse](LoginAction))

	app.Listen(TCP, ":8080")

	err := app.Run(signal.SetupSignalHandler())
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

	buf, err := DefaultCodec{}.Pack(&Packet{
		Operate:     1,
		ContentType: ContentTypeJSON,
	}, LoginRequest{Name: "test"})

	if err != nil {
		panic(err)
	}
	_, err = conn.Write(buf)
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
