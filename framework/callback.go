package framework

type Callback struct{}

func NewCallback() *Callback { return &Callback{} }

func (callback *Callback) OnConnect(fn func(conn *Connection)) *Callback    { return callback }
func (callback *Callback) OnDisconnect(fn func(conn *Connection)) *Callback { return callback }
func (callback *Callback) OnRequest(fn func(conn *Connection)) *Callback    { return callback }
