package internal

import (
	"gim/internal/application"
	"gim/internal/infrastructure"
	"gim/internal/infrastructure/config"
	"gim/internal/presentation"
	"gim/pkg/app"
	"gim/pkg/system"
	"go.uber.org/dig"
	"log"
	"time"
)



var (
	Version = "1.0.0"
)

func displayMemoryUsage()  {
	go func() {
		for {
			time.Sleep(time.Second * 5)
			log.Printf("memory usage: %.2fM\n", float64(system.GetMem())/1000/1000)
		}
	}()
}


type App struct {
	configFile, storage string
	port, limit, pushCount int
	debug bool
}

func (a *App) WithDebug(debug bool) *App  {
	a.debug = debug
	return a
}
func (a *App) WithStorage(storage string) *App {
	a.storage = storage
	return a
}
func (a *App) WithConfigFile(configFile string) *App {
	a.configFile = configFile
	return a
}
func (a *App) WithPort(port int) *App {
	a.port = port
	return a
}
func (a *App) WithLimit(limit int) *App {
	a.limit = limit
	return a
}

func (a *App) WithPushCount(count int) *App {
	a.pushCount = count
	return a
}

func (a *App) Run() {
	if a.debug {
		displayMemoryUsage()
	}
	container := app.Container()
	system.SecurePanic(a.loadConfig(container))

	infrastructure.InjectStore(container, a.storage)
	application.Inject(container)
	presentation.Inject(container)

	err := container.Invoke(a.serve)
	system.SecurePanic(err)

	system.Shutdown(func() {
		log.Println("server shutdown")
	})
}

func (a *App) loadConfig(container *dig.Container) error {
	conf := config.New()
	if a.configFile != "" {

	}
	if err := conf.LoadFile(a.configFile); err != nil {
		return err
	}

	conf.Server.Port = a.port
	conf.Message.MaxStoreSize = a.limit
	conf.Message.PushCount = a.pushCount
	return container.Provide(func() *config.Config{
		return conf
	})
}

func (a *App) serve(socket *presentation.Socket) error {
	return socket.Start()

}

func NewApp() *App {
	return &App{
		configFile: "",
		storage:    infrastructure.MemoryStore,
		port:       8080,
		limit:      10000,
		debug:      false,
		pushCount: 10,
	}
}