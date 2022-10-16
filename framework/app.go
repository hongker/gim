package framework

import (
	"context"
	"errors"
	"github.com/ebar-go/ego/utils/runtime"
	"log"
	"net"
)

// Instance represents an app interface
type Instance interface {
	// Router return an router instance
	Router() *Router

	// Listen listens for different schema and address
	Listen(protocol string, addr string) *App

	// Run runs the application with the given signal handler
	Run(stopCh <-chan struct{}) error
}

// App represents im framework public access api.
type App struct {
	options  *Options
	schemas  []Schema
	callback *Callback
	router   *Router
	event    *Event
	reactor  *Reactor
}

// Listen register different protocols
func (app *App) Listen(protocol string, addr string) *App {
	app.schemas = append(app.schemas, NewSchema(protocol, addr))
	return app
}

// Router return instance of Router
func (app *App) Router() *Router {
	return app.router
}

// Run starts the app
func (app *App) Run(stopCh <-chan struct{}) error {
	ctx := context.Background()
	if len(app.schemas) == 0 {
		return errors.New("empty listen target")
	}

	// prepare servers
	schemaCtx, schemeCancel := context.WithCancel(ctx)
	// cancel schema context when app is stopped
	defer schemeCancel()
	for _, schema := range app.schemas {
		// listen with context and connection register callback function
		if err := schema.Listen(schemaCtx.Done(), app.registerConnection); err != nil {
			return err
		}
	}

	// prepare reactor
	reactor, err := NewReactor()
	if err != nil {
		return err
	}
	reactor.engine.Use(app.router.Request)
	reactorCtx, reactorCancel := context.WithCancel(ctx)
	// cancel reactor context when app is stopped
	defer reactorCancel()
	go func() {
		defer runtime.HandleCrash()
		reactor.Run(reactorCtx.Done())
	}()

	app.reactor = reactor

	log.Println("app started")
	runtime.WaitClose(stopCh)
	return nil
}

func (app *App) registerConnection(conn net.Conn) {
	connection := NewConnection(conn, app.options.MaxReadBufferSize)
	connection.fd = app.reactor.poll.SocketFD(conn)
	if err := app.reactor.thread.RegisterConnection(connection); err != nil {
		connection.Close()
		return
	}

	connection.AddBeforeCloseHook(app.callback.disconnect, app.reactor.thread.UnregisterConnection)

	app.callback.connect(connection)

}

// New returns a new app instance
func New(opts ...Option) Instance {
	options := defaultOptions()
	for _, setter := range opts {
		setter(options)
	}

	return &App{
		options:  options,
		callback: NewCallback().OnConnect(options.OnConnect).OnDisconnect(options.OnDisconnect),
		router:   NewRouter(),
	}
}
