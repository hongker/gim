package gateway

import (
	"gim/api"
	"gim/framework"
	"gim/internal/domain/stateful"
	"github.com/ebar-go/ego/component"
	"sync"
)

// Controller represents gateway module.
type Controller struct {
	name string
	once sync.Once

	config *Config

	//engine *ego.NamedEngine
}

// Run runs the controller.
func (c *Controller) Run(stopCh <-chan struct{}) {
	c.once.Do(c.initialize)

	component.Provider().Logger().Infof("controller running: [%s]", c.name)

	handler := NewHandler(c.config.HeartbeatInterval)
	app := framework.New(framework.WithConnectCallback(func(conn *framework.Connection) {
		component.Provider().Logger().Infof("[%s] Connected", conn.UUID())
		handler.InitializeConn(conn)
	}), framework.WithDisconnectCallback(func(conn *framework.Connection) {
		component.Provider().Logger().Infof("[%s] Disconnected", conn.UUID())
		handler.FinalizeConn(conn)
	}))

	app.Use(func(ctx *framework.Context) {
		if ctx.Operate() != api.LoginOperate {
			// check user login state
			if uid := stateful.GetUidFromConnection(ctx.Conn()); uid == "" {
				ctx.Abort()
			}
		}

		ctx.Next()
	})
	app.Router().Route(api.LoginOperate, framework.StandardHandler(handler.Login))
	app.Router().Route(api.HeartbeatOperate, framework.StandardHandler(handler.Heartbeat))
	app.Router().Route(api.LogoutOperate, framework.StandardHandler(handler.Logout))
	app.Router().Route(api.MessageSendOperate, framework.StandardHandler(handler.SendMessage))
	app.Router().Route(api.MessageQueryOperate, framework.StandardHandler(handler.QueryMessage))
	app.Router().Route(api.SessionListOperate, framework.StandardHandler(handler.ListSession))
	app.Router().Route(api.ChatroomJoinOperate, framework.StandardHandler(handler.JoinChatroom))

	if err := app.Listen(framework.TCP, c.config.Address).Run(stopCh); err != nil {
		panic(err)
	}

	c.shutdown()

}

// WithName set controller name.
func (c *Controller) WithName(name string) *Controller {
	c.name = name
	return c
}

// initialize init controller dependencies.
func (c *Controller) initialize() {
	//c.engine = ego.New()
	//
	//callback := NewCallback(c.config.HeartbeatInterval)
	//
	//wss := ego.NewWebsocketServer(c.config.Address).
	//	WithWorker(c.config.WorkerNumber).
	//	OnConnect(callback.OnConnect).
	//	OnDisconnect(callback.OnDisconnect).
	//	OnMessage(callback.OnMessage)
	//
	//c.engine.WithServer(wss)
}

// shutdown shuts down the controller.
func (c *Controller) shutdown() {
	component.Provider().Logger().Infof("controller shutdown: [%s]", c.name)
}
