package framework

import (
	"errors"
	"gim/pkg/bytes"
	uuid "github.com/satori/go.uuid"
	"net"
	"sync"
)

// Connection represents client connection
type Connection struct {
	uuid             string
	conn             net.Conn
	fd               int
	once             sync.Once
	beforeCloseHooks []func(connection *Connection)
	maxReadBytes     int
}

// UIID returns the uuid associated with the connection
func (conn *Connection) UUID() string { return conn.uuid }

// Push send message to the connection
func (conn *Connection) Push(p []byte) {
	conn.Write(p)
}

// Write writes message to the connection
func (conn *Connection) Write(p []byte) (int, error) {
	return conn.conn.Write(p)
}

// Read reads message from the connection
func (conn *Connection) Read(p []byte) (int, error) {
	return conn.conn.Read(p)
}

// Close closes the connection
func (conn *Connection) Close() {
	conn.once.Do(func() {
		for _, hook := range conn.beforeCloseHooks {
			hook(conn)
		}
		_ = conn.conn.Close()
	})
}

// FD returns the file descriptor of the connection
func (conn *Connection) FD() int { return conn.fd }

// AddBeforeCloseHook adds a hook to the connection before closed
func (conn *Connection) AddBeforeCloseHook(hooks ...func(conn *Connection)) {
	conn.beforeCloseHooks = append(conn.beforeCloseHooks, hooks...)
}

// readLine reads a line message from the connection
func (conn *Connection) readLine(maxReadBytes int) ([]byte, error) {
	buf := bytes.Get(maxReadBytes)
	n, err := conn.Read(buf)
	if err != nil {
		bytes.Put(buf)
		return nil, err
	}

	if n == 0 {
		bytes.Put(buf)
		return nil, errors.New("empty packet")
	}

	return buf[:n], nil

}
func NewConnection(conn net.Conn) *Connection {
	return &Connection{conn: conn, uuid: uuid.NewV4().String()}
}
