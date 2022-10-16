package framework

import (
	"gim/framework/codec"
	"github.com/ebar-go/ego/utils/runtime/signal"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"log"
	"net"
	"testing"
)

func TestApp(t *testing.T) {
	app := New(WithConnectCallback(func(conn *Connection) {
		log.Printf("[%s] connected\n", conn.UUID())
	}), WithDisconnectCallback(func(conn *Connection) {
		log.Printf("[%s] disconnected\n", conn.UUID())
	}))

	app.Router().Route(1, StandardHandler[LoginRequest, LoginResponse](LoginAction))

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

func LoginAction(ctx *Context, req *LoginRequest) (response *LoginResponse, err error) {
	response = &LoginResponse{Token: uuid.NewV4().String(), Name: req.Name}
	return
}

func TestClient(t *testing.T) {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	defaultCodec := codec.Default()
	buf, err := defaultCodec.Pack(&codec.Packet{
		Operate:     1,
		Seq:         1,
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

	log.Println("receive: ", string(receive[:n]))
	packet, err := defaultCodec.Unpack(receive[:n])
	if err != nil {
		panic(err)
	}

	log.Println("packet:", packet.Operate, packet.Seq, packet.ContentType, string(packet.Body))

}
