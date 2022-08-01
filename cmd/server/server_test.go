package main

import (
	"bufio"
	"fmt"
	"gim/api"
	"gim/internal/domain/dto"
	"gim/pkg/binary"
	"gim/pkg/system"
	"net"
	"testing"
)

func TestNewClient(t *testing.T) {
	_, err := connect("A")
	system.SecurePanic(err)


	select {}
}

func TestClientB(t *testing.T) {
	conn, err := connect("B")
	system.SecurePanic(err)

	p := api.NewPacket()
	p.Op  =api.OperateMessageSend
	p.Marshal(dto.MessageSendRequest{
		Type:        api.PrivateMessage,
		Content:     "test",
		ClientMsgId: "",
		SessionId:   "bfba28cc-18cf-4c08-8f60-3e6802b2e143",
	})

	conn.Write(p.Encode())
	select {}
}

func connect(name string) (net.Conn, error) {
	conn, err := net.Dial("tcp", "127.0.0.1:8088")
	system.SecurePanic(err)
	p := api.NewPacket()
	p.Op = api.OperateAuth
	err = p.Marshal(&dto.UserLoginRequest{Name: name})
	system.SecurePanic(err)

	go func() {
		scanner := bufio.NewScanner(conn)
		// 处理粘包问题：先读取包体长度
		scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
			if !atEOF && len(data) > 4 {
				length := int(binary.BigEndian.Int32(data[:4]))
				if length >= 0 && length <= len(data) {
					return length, data[:length], nil
				}
			}
			return
		})
		for {
			if !scanner.Scan() {
				panic(scanner.Err())
			}
			resp := api.NewPacket()
			system.SecurePanic(resp.Decode(scanner.Bytes()))
			fmt.Println(resp.Op, string(resp.Data))
		}
	}()
	_, err = conn.Write(p.Encode())
	system.SecurePanic(err)
	return conn, nil
}