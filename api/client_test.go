package api

import (
	"gim/framework/codec"
	"gim/internal/domain/dto"
	"github.com/ebar-go/ego/utils/runtime/signal"
	"log"
	"net"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defaultCodec := codec.Default()
	go func() {
		for {
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
	}()

	login(defaultCodec, conn)

	time.Sleep(time.Second * 1)

	// send heartbeat
	go func() {
		for {
			heartbeat(defaultCodec, conn)
			time.Sleep(time.Second * 5)
		}
	}()

	<-signal.SetupSignalHandler()

}

func login(defaultCodec codec.Codec, conn net.Conn) {
	buf, err := defaultCodec.Pack(&codec.Packet{
		Operate:     LoginOperate,
		Seq:         1,
		ContentType: codec.ContentTypeJSON,
	}, dto.UserLoginRequest{ID: "1001", Name: "test01"})

	if err != nil {
		panic(err)
	}

	write(conn, buf)

}

func heartbeat(defaultCodec codec.Codec, conn net.Conn) {
	buf, err := defaultCodec.Pack(&codec.Packet{
		Operate:     HeartbeatOperate,
		Seq:         1,
		ContentType: codec.ContentTypeJSON,
	}, dto.SocketHeartbeatRequest{})

	if err != nil {
		panic(err)
	}

	write(conn, buf)

}

func write(conn net.Conn, buf []byte) {
	_, err := conn.Write(buf)
	if err != nil {
		panic(err)
	}

	log.Println("send  success")
}
