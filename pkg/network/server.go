package network

import (
	"gim/pkg/compressor"
	"log"
	"net"
)

type Server interface {
	Start() error
	SetOnRequest(fn HandleFunc)
	SetOnConnect(func(conn *Connection))
	SetOnDisconnect(func(conn *Connection))
	Use(middlewares ...HandleFunc)
}

func NewTCPServer(bind string, opts ...Option) *TcpServer {
	conf := defaultConfig()
	conf.Bind = bind
	for _, setter := range opts {
		setter(conf)
	}
	return &TcpServer{
		engine: newEngine(conf.ContextPoolSize),
		conf:   conf,
	}
}

type TcpServer struct {
	Callback

	engine *Engine

	conf *Config
}

func (s *TcpServer) Start() error {
	s.engine.Use(s.OnRequest)

	return s.init()
}

func (s *TcpServer) Use(handlers ...HandleFunc) {
	s.engine.Use(handlers...)
}

// accept 一般使用cpu核数作为参数，提高处理能力
func (s *TcpServer) init() (err error) {
	var (
		bind     = s.conf.Bind
		listener *net.TCPListener
		addr     *net.TCPAddr
	)
	if addr, err = net.ResolveTCPAddr("tcp", bind); err != nil {
		log.Printf("net.ResolveTCPAddr(tcp, %s) errors(%v)", bind, err)
		return
	}
	if listener, err = net.ListenTCP("tcp", addr); err != nil {
		log.Printf("net.ListenTCP(tcp, %s) errors(%v)", bind, err)
		return
	}

	log.Printf("start tcp listen: %s", bind)

	// 利用多线程处理连接初始化
	for i := 0; i < s.conf.Accept; i++ {
		go s.listen(listener)
	}
	return
}

func (s *TcpServer) listen(lis *net.TCPListener) {
	var (
		conn *net.TCPConn
		err  error
	)

	for {
		if conn, err = lis.AcceptTCP(); err != nil {
			// if listener close then return
			log.Printf("listener.Accept(\"%s\") errors(%v)", lis.Addr().String(), err)
			return
		}
		if err = conn.SetKeepAlive(s.conf.KeepAlive); err != nil {
			log.Printf("conn.SetKeepAlive() errors(%v)", err)
			return
		}
		if err = conn.SetReadBuffer(s.conf.Rcvbuf); err != nil {
			log.Printf("conn.SetReadBuffer() errors(%v)", err)
			return
		}
		if err = conn.SetWriteBuffer(s.conf.Sndbuf); err != nil {
			log.Printf("conn.SetWriteBuffer() errors(%v)", err)
			return
		}

		if s.conf.Debug {
			log.Printf("client new request ,ip: %v", conn.RemoteAddr())
		}

		// 一个goroutine处理一个连接
		go s.handle(conn)

	}
}

func (s *TcpServer) handle(conn *net.TCPConn) {
	if s.conf.Debug {
		lAddr := conn.LocalAddr().String()
		rAddr := conn.RemoteAddr().String()
		log.Printf("start handle \"%s\" with \"%s\"", lAddr, rAddr)
	}

	// 初始化连接
	connection := wantsCompressConnection(conn, s.conf.ContentEncoding)
	connection.init(s.conf.QueueSize, s.conf.DataLength, s.conf.DataMaxLength)

	// 开启连接事件回调
	s.OnConnect(connection)

	// 处理接收数据
	connection.handleRequest(s.engine)

	// 关闭连接事件回调
	s.OnDisconnect(connection)
}

func wantsCompressConnection(conn net.Conn, contentEncoding string) *Connection {
	instance := conn
	c, err := compressor.NewCompressingResponseWriter(conn, contentEncoding)
	if err == nil {
		instance = &CompressorNetConn{Conn: conn, compressorWriter: c}
	}
	return &Connection{instance: instance}
}
