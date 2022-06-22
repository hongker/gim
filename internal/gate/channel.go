package gate

import "gim/pkg/network"

type Channel struct {
	key  string // 用户ID
	conn *network.Connection
}

func (c *Channel) Key() string {
	return c.key
}

func NewChannel(key string, conn *network.Connection) *Channel {
	return &Channel{key: key, conn: conn}
}
