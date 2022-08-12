package network

import (
	"gim/pkg/binary"
	"gim/pkg/bytes"
	uuid "github.com/satori/go.uuid"
	"log"
	"net"
	"sync"
)

type Connection struct {
	instance net.Conn
	id       string

	sendQueue chan []byte

	once *sync.Once
	done chan struct{} // 关闭标识
	packetDataLength int
}

// ID return unique id of connection
func (conn *Connection) ID() string {
	return conn.id
}

// IP return ip of connection
func (conn *Connection) IP() string {
	ip, _, _ := net.SplitHostPort(conn.instance.RemoteAddr().String())
	return ip
}

// Push send message to client, if queue is full, msg will be disposed
func (conn *Connection) Push(msg []byte) {
	select {
	case conn.sendQueue <- msg:
	default:
	}
}

// Close shutdown the connection
func (conn *Connection) Close() {
	conn.once.Do(func() {
		close(conn.done)
		//close(conn.sendQueue)
		_ = conn.instance.Close()
	})

}

// =====================private function======================================

// init initialize connection param
func (conn *Connection) init(sendQueueSize int, packetDataLength int) {
	conn.id = uuid.NewV4().String()
	conn.sendQueue = make(chan []byte, sendQueueSize)
	conn.once = new(sync.Once)
	conn.done = make(chan struct{})
	conn.packetDataLength = packetDataLength

	// 分发响应数据
	go conn.dispatchResponse()
}

// dispatchResponse
func (conn *Connection) dispatchResponse() {
	defer conn.Close()

	var err error
	for {
		select {
		case <-conn.done:
			return
		default:
			msg, ok := <-conn.sendQueue
			if !ok { // 队列已关闭
				return
			}

			// 写数据
			_, err = conn.instance.Write(msg)
			bytes.Put(msg) // 回收数组
			if err != nil {
				return
			}
		}
	}
}

// handleRequest 处理请求
func (conn *Connection) handleRequest(engine *Engine) {
	defer conn.Close()
	// 利用对象池实例化context,避免GC
	// 会导致内存随着连接的增加而增加
	ctxPool := engine.contextPool()

	for {
		select {
		case <-conn.done: // 退出
			return
		default:
			b := bytes.Get(512)
			_, err := conn.instance.Read(b[:conn.packetDataLength])
			if err != nil {
				log.Println("read error: ", err)
				return
			}
			length := int(binary.BigEndian.Int32(b[:conn.packetDataLength]))
			_, err = conn.instance.Read(b[conn.packetDataLength:length])
			if err != nil {
				log.Println("read error: ", err)
				return
			}

			// 通过对象池初始化时，会导致内存缓慢上涨,直到稳定
			ctx := ctxPool.Get().(*Context)
			ctx.Reset(b[:length], conn)

			// 执行回调
			go func() {
				defer func() {
					ctxPool.Put(ctx)
					bytes.Put(b)
				}()
				ctx.Run()
			}()
		}

	}

}

