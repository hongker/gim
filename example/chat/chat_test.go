package main

import (
	"gim/framework/codec"
	"github.com/ebar-go/ego/utils/runtime/signal"
	"log"
	"net"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	stop := signal.SetupSignalHandler()
	go initializeUser(stop, "test01")
	go initializeUser(stop, "test02")

	<-stop
}
func initializeUser(stop <-chan struct{}, name string) {
	conn, err := net.Dial("tcp", ":8090")
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

	login(defaultCodec, conn, name)

	time.Sleep(time.Second * 1)

	subscribeChannel(defaultCodec, conn)

	time.Sleep(time.Second * 1)
	// send heartbeat
	go func() {
		for {
			sendMessage(defaultCodec, conn)
			time.Sleep(time.Second * 5)
		}
	}()

	<-stop

}

func login(defaultCodec codec.Codec, conn net.Conn, name string) {
	buf, err := defaultCodec.Pack(&codec.Packet{
		Operate:     1,
		Seq:         1,
		ContentType: codec.ContentTypeJSON,
	}, LoginRequest{Name: name})

	if err != nil {
		panic(err)
	}

	write(conn, buf)

}

func subscribeChannel(defaultCodec codec.Codec, conn net.Conn) {
	buf, err := defaultCodec.Pack(&codec.Packet{
		Operate:     2,
		Seq:         1,
		ContentType: codec.ContentTypeJSON,
	}, SubscribeChannelRequest{ID: "10001"})

	if err != nil {
		panic(err)
	}

	write(conn, buf)

}

func sendMessage(defaultCodec codec.Codec, conn net.Conn) {
	buf, err := defaultCodec.Pack(&codec.Packet{
		Operate:     3,
		Seq:         1,
		ContentType: codec.ContentTypeJSON,
	}, SendMessageRequest{ChannelID: "10001", Content: "hello world," + time.Now().String()})

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
