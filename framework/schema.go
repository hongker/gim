package framework

import "github.com/ebar-go/ego/utils/runtime"

type Protocol string

const (
	TCP       Protocol = "tcp"
	WEBSOCKET Protocol = "ws"
	HTTP      Protocol = "http"
)

type Schema struct {
	Protocol Protocol
	Addr     string
}

func (schema Schema) Listen(stopCh <-chan struct{}) error {
	runtime.WaitClose(stopCh, schema.Stop)
	return nil
}

func (schema Schema) Stop() {}
func NewSchema(protocol Protocol, addr string) Schema {
	return Schema{
		Protocol: protocol,
		Addr:     addr,
	}
}
