package framework

import (
	"context"
	"errors"
	"github.com/ebar-go/ego/utils/runtime"
	"log"
	"net"
)

// Engine represents im framework public access api.
type Engine struct {
	schemas  []Schema
	callback *Callback
	codec    Codec
	router   *Router
	event    *Event
	reactor  *Reactor
}

// WithProtocol set different protocol
func (engine *Engine) Listen(protocol Protocol, addr string) *Engine {
	engine.schemas = append(engine.schemas, NewSchema(protocol, addr))
	return engine
}

// WithCallback use callback
func (engine *Engine) WithCallback(callback *Callback) *Engine {
	engine.callback = callback
	return engine
}

// WithCodec use codec to pack/unpack message.
func (engine *Engine) WithCodec(codec Codec) *Engine {
	engine.codec = codec
	return engine
}

// WithRouter set router
func (engine *Engine) WithRouter(router *Router) *Engine {
	engine.router = router
	return engine
}

// WithEvent set event
func (engine *Engine) WithEvent(event *Event) *Engine {
	engine.event = event
	return engine
}

// Start starts the engine
func (engine *Engine) Run(stopCh <-chan struct{}) error {
	ctx := context.Background()

	err := runtime.Call(func() error {
		return engine.runServer(ctx)
	}, func() error {
		return engine.runReactor(ctx)
	})
	if err != nil {
		return err
	}

	log.Println("engine started")
	runtime.WaitClose(stopCh)

	return nil
}

func (engine *Engine) runServer(ctx context.Context) error {
	if len(engine.schemas) == 0 {
		return errors.New("empty listen target")
	}
	// start listen protocol
	schemaContext, schemeCancel := context.WithCancel(ctx)
	defer schemeCancel()
	for _, schema := range engine.schemas {
		err := schema.Listen(schemaContext.Done(), engine.handle)
		runtime.HandleError(err, func(err error) {
			log.Println("listen error:", err)
		})
	}

	return nil
}

func (engine *Engine) runReactor(ctx context.Context) error {
	reactor, err := NewReactor()
	if err != nil {
		return err
	}
	refactorContext, refactorCancel := context.WithCancel(ctx)
	defer refactorCancel()
	go func() {
		defer runtime.HandleCrash()
		reactor.Run(refactorContext.Done())
	}()
	return nil

}

func (engine *Engine) handle(conn net.Conn) {
}

// New returns a new engine instance
func New() *Engine {
	return &Engine{}
}
