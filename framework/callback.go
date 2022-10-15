package framework

type ConnectionHandler func(conn *Connection)
type Callback struct {
	connect    ConnectionHandler
	disconnect ConnectionHandler
	request    ConnectionHandler
}

func NewCallback() *Callback {
	return &Callback{
		connect: func(conn *Connection) {

		},
		disconnect: func(conn *Connection) {

		},
		request: func(conn *Connection) {

		},
	}
}

func (callback *Callback) OnConnect(fn ConnectionHandler) *Callback {
	if fn != nil {
		callback.connect = fn
	}
	return callback
}
func (callback *Callback) OnDisconnect(fn ConnectionHandler) *Callback {
	if fn != nil {
		callback.disconnect = fn
	}
	return callback
}
func (callback *Callback) OnRequest(fn ConnectionHandler) *Callback {
	if fn != nil {
		callback.request = fn
	}
	return callback
}
