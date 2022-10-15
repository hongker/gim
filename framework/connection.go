package framework

import (
	"net"
	"sync"
)

type Connection struct {
	conn             net.Conn
	fd               int
	once             sync.Once
	beforeCloseHooks []func(connection *Connection)
}

func (conn *Connection) Push(p []byte)               {}
func (conn *Connection) Write(p []byte) (int, error) { return 0, nil }
func (conn *Connection) Read(p []byte) (int, error)  { return 0, nil }
func (conn *Connection) Close() {
	conn.once.Do(func() {
		for _, hook := range conn.beforeCloseHooks {
			hook(conn)
		}
		conn.conn.Close()
	})
}
func (conn *Connection) NewContext() *Context { return &Context{} }
func (conn *Connection) FD() int              { return conn.fd }
func (conn *Connection) AddBeforeCloseHook(hooks ...func(conn *Connection)) {
	conn.beforeCloseHooks = append(conn.beforeCloseHooks, hooks...)
}
func NewConnection(conn net.Conn) *Connection {
	return &Connection{conn: conn}
}
