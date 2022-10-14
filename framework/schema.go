package framework

type Schema struct {
	Protocol string
	Addr     string
}

const (
	tcp       = "tcp"
	websocket = "websocket"
	http      = "http"
)

func NewTcpSchema(address string) *Schema {
	return &Schema{Protocol: tcp, Addr: address}
}

func NewWebsocketSchema(address string) *Schema {
	return &Schema{Protocol: websocket, Addr: address}
}
