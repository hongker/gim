package main

import (
	"bufio"
	"context"
	"fmt"
	"gim/api"
	"gim/internal/domain/dto"
	"gim/pkg/binary"
	"gim/pkg/system"
	"github.com/rcrowley/go-metrics"
	uuid "github.com/satori/go.uuid"
	"log"
	"math/rand"
	"net"
	"os"
	"testing"
	"time"
)

func connect(uid, name string, joinGroup bool, received bool) (net.Conn, error) {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	system.SecurePanic(err)
	p := api.NewPacket()
	p.Op = api.OperateAuth
	err = p.Marshal(&dto.UserLoginRequest{UID: uid, Name: name})
	system.SecurePanic(err)

	if received {
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
				log.Println(resp.Op, string(resp.Data))
			}
		}()
	}

	go func() {
		for {
			time.Sleep(time.Second * time.Duration(rand.Intn(30)+30))
			heartbeatPacket := api.BuildPacket(api.OperateHeartbeat, dto.UserHeartbeatRequest{})
			_, err = conn.Write(heartbeatPacket.Encode())
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

func TestQueryMessage(t *testing.T) {
	conn, err := connect("10001", "someUserA", true, true)
	system.SecurePanic(err)

	for {
		p := api.NewPacket()
		p.Op = api.OperateMessageQuery
		p.Marshal(dto.MessageQueryRequest{
			SessionId: "group:1",
			Limit:     3,
			Last:      time.Now().UnixNano(),
		})

		conn.Write(p.Encode())
		time.Sleep(time.Second * 30)
	}

}

func TestGroupQueryMember(t *testing.T) {
	conn, err := connect("10004", "someUserC", true, true)
	system.SecurePanic(err)

	for {
		p := api.NewPacket()
		p.Op = api.OperateGroupMember
		p.Marshal(dto.GroupMemberQuery{
			GroupId: "1",
		})

		conn.Write(p.Encode())
		time.Sleep(time.Second * 10)
	}

}

func TestSendUserMessage(t *testing.T) {
	conn, err := connect("10002", "someUser", false, true)
	system.SecurePanic(err)

	for {
		p := newGroupMessagePacket()

		conn.Write(p.Encode())
		time.Sleep(time.Millisecond * 500)
	}
}

func newGroupMessagePacket() *api.Packet {
	p := api.NewPacket()
	p.Op = api.OperateMessageSend
	p.Marshal(dto.MessageSendRequest{
		Type:        api.GroupSession,
		Content:     "testRoom",
		ContentType: api.TextMessage,
		RequestId:   uuid.NewV4().String(),
		TargetId:    "1",
	})
	return p
}

func TestSendGroupMessage(t *testing.T) {
	conn, err := connect("10003", "someUserB", true, true)
	system.SecurePanic(err)

	for {
		p := newGroupMessagePacket()

		conn.Write(p.Encode())
		time.Sleep(time.Millisecond * 100)
	}

}

func BenchmarkSendMessage(b *testing.B) {
	opsRate := metrics.NewRegisteredTimer("ops", nil)

	ch := make(chan net.Conn, 50)
	n := 1000
	ctx, cancel := context.WithCancel(context.Background())
	for i := 0; i < 50; i++ {
		go func(ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					c, err := connect(fmt.Sprintf("%d", time.Now().UnixNano()), uuid.NewV4().String(), true, false)
					if err == nil {
						ch <- c
					}
				}
			}

		}(ctx)
	}
	connections := make([]net.Conn, 0, 1024)
	for len(connections) < n {
		connections = append(connections, <-ch)
	}
	cancel()

	go func() {
		metrics.Log(metrics.DefaultRegistry, 5*time.Second, log.New(os.Stderr, "metrics: ", log.Lmicroseconds))
	}()

	for i := 0; i < 50; i++ {
		go func() {
			for {
				n := rand.Intn(len(connections) - 1)
				c := connections[n]
				p := newGroupMessagePacket()

				before := time.Now()
				if _, err := c.Write(p.Encode()); err != nil {
					_ = c.Close()
					log.Println(err)

				} else {
					opsRate.Update(time.Now().Sub(before))
				}
			}

		}()
	}
	select {}
}
