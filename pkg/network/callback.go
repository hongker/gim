package network

// Callback 回调
type Callback struct {
	onRequest    HandleFunc
	onConnect    func(conn *Connection)
	onDisconnect func(conn *Connection)
}

func (c *Callback) SetOnConnect(onConnect func(conn *Connection)) {
	c.onConnect = onConnect
}

func (c *Callback) SetOnDisconnect(onDisconnect func(conn *Connection)) {
	c.onDisconnect = onDisconnect
}

// SetOnRequest 接收请求回调
func (c *Callback) SetOnRequest(hookFunc HandleFunc) {
	c.onRequest = hookFunc
}

func (c *Callback) OnRequest(ctx *Context) {
	if c.onRequest == nil {
		return
	}
	c.onRequest(ctx)
}

func (c *Callback) OnConnect(conn *Connection) {
	if c.onConnect == nil {
		return
	}
	c.onConnect(conn)
}

func (c *Callback) OnDisconnect(conn *Connection) {
	if c.onDisconnect == nil {
		return
	}
	c.onDisconnect(conn)
}
