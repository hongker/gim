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
	"time"
)


func TestQueyrMessage(t *testing.T) {
	conn, err := connect("someUserA", true)
	system.SecurePanic(err)

	for {
		p := api.NewPacket()
		p.Op  =api.OperateMessageQuery
		p.Marshal(dto.MessageQueryRequest{
			SessionId: "group:1001",
			Limit:     3,
			Last:      time.Now().UnixNano(),
		})

		conn.Write(p.Encode())
		time.Sleep(time.Second * 30)
	}

}

func TestSendUserMessage(t *testing.T) {
	conn, err := connect("someUserB", false)
	system.SecurePanic(err)

	for {
		p := api.NewPacket()
		p.Op  =api.OperateMessageSend
		p.Marshal(dto.MessageSendRequest{
			Type:        api.UserSession,
			Content:     "testSendUserMessage",
			ContentType: api.TextMessage,
			ClientMsgId: "",
			TargetId:   "8f66603c-823a-458e-93e2-647ca52fe122",
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
			ClientMsgId: "",
			TargetId:   "1001",
		})

		conn.Write(p.Encode())
		time.Sleep(time.Millisecond * 500)
	}

}

func connect(name string, joinGroup bool) (net.Conn, error) {
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

	if joinGroup {
		time.Sleep(time.Second)
		joinGroupPacket := api.BuildPacket(api.OperateGroupJoin, dto.GroupJoinRequest{GroupId: "1001"})
		_, err = conn.Write(joinGroupPacket.Encode())
	}

	return conn, nil
}