package network

// Request 请求
type Request struct {
	body []byte // 请求内容
}

// Body 获取原始报文
func (request Request) Body() []byte {
	return request.body
}

// Context 上下文
type Context struct {
	index int8

	engine *Engine

	connection *Connection
	request    Request
}

func (c *Context) Connection() *Connection {
	return c.connection
}

// Request 获取请求
func (c *Context) Request() Request {
	return c.request
}

// Output 输出数据到客户端
func (c *Context) Output(msg []byte) {
	c.connection.Push(msg)
}

func (c *Context) Reset(body []byte, connection *Connection) {
	c.index = 0
	c.request.body = body
	c.connection = connection
}

func (ctx *Context) Run() {
	ctx.engine.handleChains[0](ctx)
}

func (ctx *Context) Next() {
	if ctx.index < maxIndex {
		ctx.index++
		ctx.engine.handleChains[ctx.index](ctx)
	}
}
func (ctx *Context) Abort() {
	ctx.index = maxIndex
}
