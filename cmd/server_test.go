package main

import (
	"bufio"
	"fmt"
	"gim/api"
	"gim/internal/domain/dto"
	"gim/pkg/binary"
	"gim/pkg/system"
	uuid "github.com/satori/go.uuid"
	"net"
	"testing"
	"time"
)


func TestQueryMessage(t *testing.T) {
	conn, err := connect("someUserA", true)
	system.SecurePanic(err)

	for {
		p := api.NewPacket()
		p.Op  =api.OperateMessageQuery
		p.Marshal(dto.MessageQueryRequest{
			SessionId: "group:1",
			Limit:     3,
			Last:      time.Now().UnixNano(),
		})

		conn.Write(p.Encode())
		time.Sleep(time.Second * 30)
	}

}

func TestSendUserMessage(t *testing.T) {
	conn, err := connect("someUser", false)
	system.SecurePanic(err)

	for {
		p := api.NewPacket()
		p.Op  =api.OperateMessageSend
		p.Marshal(dto.MessageSendRequest{
			Type:        api.UserSession,
			Content:     "testSendUserMessage",
			ContentType: api.TextMessage,
			RequestId: "",
			TargetId:   "8",
		})

		conn.Write(p.Encode())
		time.Sleep(time.Millisecond * 500)
	}
}


func TestSendGroupMessage(t *testing.T) {
	conn, err := connect("someUserB", true)
	system.SecurePanic(err)

	for {
		p := api.NewPacket()
		p.Op  =api.OperateMessageSend
		p.Marshal(dto.MessageSendRequest{
			Type:        api.GroupSession,
			Content:     "testRoom",
			ContentType: api.TextMessage,
			RequestId: uuid.NewV4().String(),
			TargetId:   "1",
		})

		conn.Write(p.Encode())
		time.Sleep(time.Millisecond * 500)
	}

}

func connect(name string, joinGroup bool) (net.Conn, error) {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
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

	if joinGroup {
		time.Sleep(time.Second)
		joinGroupPacket := api.BuildPacket(api.OperateGroupJoin, dto.GroupJoinRequest{GroupId: "1"})
		_, err = conn.Write(joinGroupPacket.Encode())
	}

	return conn, nil
}