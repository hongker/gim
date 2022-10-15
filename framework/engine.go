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
	router   *Router
	event    *Event
	reactor  *Reactor
}

// WithProtocol set different protocol
func (engine *Engine) Listen(protocol Protocol, addr string) *Engine {
	engine.schemas = append(engine.schemas, NewSchema(protocol, addr))
	return engine
}

// Callback return instance of Callback
func (engine *Engine) Callback() *Callback {
	return engine.callback
}

// Router return instance of Router
func (engine *Engine) Router() *Router {
	return engine.router
}

// WithEvent set event
func (engine *Engine) WithEvent(event *Event) *Engine {
	engine.event = event
	return engine
}

// Start starts the engine
func (engine *Engine) Run(stopCh <-chan struct{}) error {

	ctx := context.Background()

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

	reactor, err := NewReactor()
	if err != nil {
		return err
	}
	reactor.container.Use(engine.router.Request())
	reactorContext, reactorCancel := context.WithCancel(ctx)
	defer reactorCancel()
	go func() {
		defer runtime.HandleCrash()
		reactor.Run(reactorContext.Done())
	}()

	engine.reactor = reactor

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
	reactor.container.Use(engine.router.Request())
	refactorContext, refactorCancel := context.WithCancel(ctx)
	defer refactorCancel()
	go func() {
		defer runtime.HandleCrash()
		reactor.Run(refactorContext.Done())
	}()

	engine.reactor = reactor
	return nil

}

func (engine *Engine) handle(conn net.Conn) {
	connection := NewConnection(conn)
	connection.fd = engine.reactor.poll.SocketFD(conn)
	if err := engine.reactor.thread.Add(connection); err != nil {
		connection.Close()
		return
	}
	connection.AddBeforeCloseHook(engine.callback.disconnect, engine.reactor.thread.Remove)

	engine.callback.connect(connection)

}

// New returns a new engine instance
func New() *Engine {
	return &Engine{callback: NewCallback(), router: NewRouter()}
}
