package framework

import (
	"errors"
	"gim/pkg/binary"
	"gim/pkg/bytes"
	"github.com/ebar-go/ego/utils/runtime"
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
func (conn *Connection) readLine(packetLengthSize int) ([]byte, error) {
	var (
		n int
	)

	// get bytes from pool
	buf := bytes.Get(conn.maxReadBytes)

	lastErr := runtime.Call(func() error {
		var err error
		// if not set packetLengthSize, read buf directly
		if packetLengthSize == 0 {
			n, err = conn.Read(buf)
			return err
		}

		// process tcp sticky package, read packet length first
		_, err = conn.Read(buf[:packetLengthSize])
		if err != nil {
			return err
		}

		packetLength := int(binary.BigEndian.Int32(buf[:packetLengthSize]))
		if packetLength > conn.maxReadBytes {
			return errors.New("packet exceeded")
		}
		_, err = conn.Read(buf[packetLengthSize:packetLength])
		n = packetLength
		return err
	}, func() error {
		if n == 0 {
			return errors.New("empty packet")
		}
		return nil
	})

	if lastErr != nil {
		// release bytes immediately
		bytes.Put(buf)
		return nil, lastErr
	}

	return buf[:n], nil

}
func NewConnection(conn net.Conn, maxReadBytes int) *Connection {
	return &Connection{conn: conn, uuid: uuid.NewV4().String(), maxReadBytes: maxReadBytes}
}
