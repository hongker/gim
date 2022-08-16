package types

import (
	"gim/pkg/network"
	"log"
	"time"
)

type Channel struct {
	key       string // 用户ID
	conn      *network.Connection
	timer     *time.Timer
	onExpired func()
}

func (c *Channel) SetKey(key string) {
	c.key = key
}

func (c *Channel) Key() string {
	return c.key
}

func (c *Channel) Conn() *network.Connection {
	return c.conn
}

func (c *Channel) ResetTimer(expired time.Duration) {
	if c.timer == nil {
		c.timer = time.AfterFunc(expired, c.onExpired)
		return
	}
	c.timer.Reset(expired)
}

func NewChannel(conn *network.Connection) *Channel {
	c := &Channel{conn: conn, onExpired: func() {
		conn.Close()
		log.Println("connection was closed, because server not receive the heartbeat request")
	}}
	return c
}
