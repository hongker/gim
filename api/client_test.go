package api

import (
	"gim/framework/codec"
	"gim/internal/domain/dto"
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
	buf, err := defaultCodec.Pack(&codec.Packet{
		Operate:     LoginOperate,
		Seq:         1,
		ContentType: codec.ContentTypeJSON,
	}, dto.UserLoginRequest{ID: "1001", Name: "test01"})

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

	time.Sleep(time.Second * 30)

	// send heartbeat

}
