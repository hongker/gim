package framework

import (
	"context"
	"errors"
	"github.com/ebar-go/ego/utils/runtime"
	"log"
	"net"
	"sync"
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
	once     sync.Once
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
		if err := schema.Listen(schemaCtx.Done(), app.handleNewConnection); err != nil {
			return err
		}

		log.Printf("start listener: %v\n", schema)
	}

	// prepare reactor
	app.reactor.engine.Use(app.router.unpack)
	app.reactor.engine.Use(app.options.Middlewares...)
	app.reactor.engine.Use(app.router.onRequest)

	reactorCtx, reactorCancel := context.WithCancel(ctx)
	// cancel reactor context when app is stopped
	defer reactorCancel()
	go func() {
		defer runtime.HandleCrash()
		app.reactor.Run(reactorCtx.Done())
	}()

	runtime.WaitClose(stopCh, app.shutdown)
	return nil
}

func (app *App) shutdown() {
	log.Println("application shutdown complete")
}

func (app *App) handleNewConnection(conn net.Conn) {
	connection := NewConnection(conn, app.reactor.poll.SocketFD(conn))
	if err := app.reactor.poll.Add(connection.fd); err != nil {
		connection.Close()
		return
	}
	app.reactor.thread.RegisterConnection(connection)

	connection.AddBeforeCloseHook(
		app.callback.handleDisconnect,
		func(conn *Connection) {
			_ = app.reactor.poll.Remove(conn.fd)
		},
		app.reactor.thread.UnregisterConnection,
	)

	app.callback.handleConnect(connection)

}

// New returns a new app instance
func New(opts ...Option) Instance {
	options := defaultOptions()
	for _, setter := range opts {
		setter(options)
	}

	return &App{
		options:  options,
		reactor:  options.NewReactor(),
		callback: NewCallback().OnConnect(options.OnConnect).OnDisconnect(options.OnDisconnect),
		router:   NewRouter(),
	}
}
