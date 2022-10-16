package framework

import (
	"context"
	"errors"
	"github.com/ebar-go/ego/utils/runtime"
	"net"
)

// App represents im framework public access api.
type App struct {
	schemas  []Schema
	callback *Callback
	router   *Router
	event    *Event
	reactor  *Reactor
}

// Listen register different protocols
func (app *App) Listen(protocol Protocol, addr string) *App {
	app.schemas = append(app.schemas, NewSchema(protocol, addr))
	return app
}

// Callback return instance of Callback
func (app *App) Callback() *Callback {
	return app.callback
}

// Router return instance of Router
func (app *App) Router() *Router {
	return app.router
}

// WithEvent set event
func (app *App) WithEvent(event *Event) *App {
	app.event = event
	return app
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
	reactor.engine.Use(app.router.Request())
	reactorCtx, reactorCancel := context.WithCancel(ctx)
	// cancel reactor context when app is stopped
	defer reactorCancel()
	go func() {
		defer runtime.HandleCrash()
		reactor.Run(reactorCtx.Done())
	}()

	app.reactor = reactor

	runtime.WaitClose(stopCh)
	return nil
}

func (app *App) registerConnection(conn net.Conn) {
	connection := NewConnection(conn)
	connection.fd = app.reactor.poll.SocketFD(conn)
	if err := app.reactor.thread.RegisterConnection(connection); err != nil {
		connection.Close()
		return
	}

	connection.AddBeforeCloseHook(app.callback.disconnect, app.reactor.thread.UnregisterConnection)

	app.callback.connect(connection)

}

// New returns a new app instance
func New() *App {
	return &App{callback: NewCallback(), router: NewRouter()}
}
