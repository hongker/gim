package network

import (
	"bufio"
	"gim/pkg/binary"
	uuid "github.com/satori/go.uuid"
	"log"
	"net"
	"sync"
)

type Connection struct {
	id       string
	instance net.Conn

	sendQueue chan []byte
	scanner   *bufio.Scanner // 读取请求数据

	once *sync.Once
	done chan struct{} // 关闭标识
}

func (conn *Connection) ID() string {
	return conn.id
}

func (conn *Connection) init(sendQueueSize int, packetDataLength int) {
	conn.id = uuid.NewV4().String()
	conn.sendQueue = make(chan []byte, sendQueueSize)
	conn.scanner = conn.getScanner(packetDataLength)
	conn.once = new(sync.Once)
	conn.done = make(chan struct{})
}

func (conn *Connection) Push(msg []byte) {
	select {
	case conn.sendQueue <- msg:
	default:
	}
}
func (conn *Connection) IP() string {
	ip, _, _ := net.SplitHostPort(conn.instance.RemoteAddr().String())
	return ip
}

// Close 关闭请求
func (conn *Connection) Close() {
	conn.once.Do(func() {
		close(conn.done)
		close(conn.sendQueue)
		_ = conn.instance.Close()
	})

}

// 分发数据
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
			if _, err = conn.instance.Write(msg); err != nil {
				return
			}
		}
	}
}

func (conn *Connection) getScanner(packetDataLength int) *bufio.Scanner {
	scan := bufio.NewScanner(conn.instance)
	if packetDataLength <= 0 {
		return scan
	}

	// 处理粘包问题：先读取包体长度
	scan.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if !atEOF && len(data) > packetDataLength {
			length := int(binary.BigEndian.Int32(data[:packetDataLength]))
			if length >= 0 && length <= len(data) {
				return length, data[packetDataLength:length], nil
			}
		}
		return
	})
	return scan
}

// handleRequest 处理请求
func (conn *Connection) handleRequest(engine *Engine) {
	defer conn.Close()
	// 利用对象池实例化context,避免GC
	// 会导致内存随着连接的增加而增加
	ctxPool := engine.ContextPool()

	for {
		select {
		case <-conn.done: // 退出
			return
		default:
			if !conn.scanner.Scan() {
				log.Println("scanner failed:", conn.scanner.Err())
				return
			}
			// 通过对象池初始化时，会导致内存缓慢上涨,直到稳定
			//ctx := &Context{engine: engine}
			ctx := ctxPool.Get().(*Context)
			ctx.Reset(conn.scanner.Bytes(), conn)

			// 执行回调
			ctx.Run()

			// 回收
			ctxPool.Put(ctx)
		}

	}

}
