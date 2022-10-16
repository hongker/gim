package framework

type ConnectionHandler func(conn *Connection)
type Callback struct {
	connect    ConnectionHandler
	disconnect ConnectionHandler
}

func NewCallback() *Callback {
	return &Callback{}
}

func (callback *Callback) OnConnect(fn ConnectionHandler) *Callback {
	callback.connect = fn
	return callback
}
func (callback *Callback) OnDisconnect(fn ConnectionHandler) *Callback {
	callback.disconnect = fn
	return callback
}

func (callback *Callback) handleConnect(conn *Connection) {
	if callback.connect != nil {
		callback.connect(conn)
	}
}

func (callback *Callback) handleDisconnect(conn *Connection) {
	if callback.disconnect != nil {
		callback.disconnect(conn)
	}
}
