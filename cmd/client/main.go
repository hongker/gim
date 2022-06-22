package main

import (
	"bufio"
	"fmt"
	"gim/api"
	"gim/api/client"
	"gim/api/protocol"
	"gim/pkg/binary"
	uuid "github.com/satori/go.uuid"
	"log"
	"net"
	"time"
)

const (
	serverUrl = "127.0.0.1:8001"
)

func getScanner(conn net.Conn, packetDataLength int) *bufio.Scanner {
	scan := bufio.NewScanner(conn)
	if packetDataLength <= 0 {
		return scan
	}

	// 处理粘包问题：先读取包体长度
	scan.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if !atEOF && len(data) > packetDataLength {
			length := int(binary.BigEndian.Int32(data[:packetDataLength]))
			if length <= len(data) {
				return length, data[:length], nil
			}
		}
		return
	})
	return scan
}

func main() {
	conn, err := net.Dial("tcp", serverUrl)
	if err != nil {
		panic(err)
	}

	scanner := getScanner(conn, 4)
	go func() {
		for {
			if !scanner.Scan() {
				log.Println("read failed:", scanner.Err())
				return
			}
			response := protocol.Proto{}
			if err = response.Unpack(scanner.Bytes()); err != nil {
				panic(err)
			}
			fmt.Println("receive:", response.Op, string(response.Body))
		}
	}()
	seq := int32(0)
	go func() {
		time.Sleep(time.Second * 10)
		for {
			time.Sleep(time.Second * 5)
			proto := &protocol.Proto{Ver: 1}
			proto.Op = api.OperateMessageSend
			body, err := proto.Marshal(&client.MessageSendRequest{
				SessionType: api.SessionTypeGroup,
				TargetId:    "1001",
				Type:        "text",
				Content:     "今天天气不错",
				ClientMsgId: "client:" + uuid.NewV4().String(),
			})
			if err != nil {
				panic(err)
			}
			_, err = conn.Write(body)
			if err != nil {
				panic(err)
			}
		}
	}()
	for {
		proto := &protocol.Proto{Ver: 1}
		var body []byte
		if seq == 0 { // 发送验证包
			proto.Op = api.OperateAuth
			body, err = proto.Marshal(&client.AuthRequest{AppId: "10001", Name: "test01"})
		} else {
			proto.Op = api.OperateMessageQuery
			body, err = proto.Marshal(&client.MessageQueryRequest{
				SessionId: "1:1001",
				LastMsgId: "",
				Count:     1,
			})
		}
		if err != nil {
			panic(err)
		}
		seq = seq + 1

		_, err = conn.Write(body)
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second * 5)

	}
}
