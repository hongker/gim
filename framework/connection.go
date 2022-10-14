package framework

type Connection struct{}

func (conn *Connection) Write(p []byte) (int, error) { return 0, nil }
func (conn *Connection) Read(p []byte) (int, error)  { return 0, nil }
func (conn *Connection) Close() error                { return nil }
func (conn *Connection) NewContext() *Context        { return &Context{} }
