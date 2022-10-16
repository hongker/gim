package framework

import (
	"context"
	"gim/framework/codec"
	"github.com/ebar-go/ego/utils/runtime/signal"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"log"
	"net"
	"testing"
)

func TestApp(t *testing.T) {
	app := New()

	app.Router().Route(2, StandardHandler[LoginRequest, LoginResponse](LoginAction))

	err := app.Listen(TCP, ":8080").
		Listen(WEBSOCKET, ":8081").
		Run(signal.SetupSignalHandler())
	assert.Nil(t, err)

}

type LoginRequest struct{ Name string }
type LoginResponse struct {
	Token string
	Name  string
}

func LoginAction(ctx context.Context, req *LoginRequest) (response *LoginResponse, err error) {
	response = &LoginResponse{Token: uuid.NewV4().String(), Name: req.Name}
	return
}

func TestClient(t *testing.T) {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	buf, err := codec.Default().Pack(&codec.Packet{
		Operate:     1,
		ContentType: codec.ContentTypeJSON,
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
