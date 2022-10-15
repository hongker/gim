package framework

type Schema struct {
	Protocol Protocol
	Addr     string
}

type Protocol string

const (
	TCP       Protocol = "tcp"
	WEBSOCKET Protocol = "ws"
	HTTP      Protocol = "http"
)

func NewSchema(protocol Protocol, addr string) Schema {
	return Schema{
		Protocol: protocol,
		Addr:     addr,
	}
}
