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
	conf *config.Config
	configFile string
}

func (a *App) WithDebug(debug bool) *App  {
	if debug {
		displayMemoryUsage()
	}
	return a
}
func (a *App) WithStorage(storage string) *App {
	a.conf.Server.Store = storage
	return a
}
func (a *App) WithConfigFile(configFile string) *App {
	a.configFile = configFile
	return a
}
func (a *App) WithPort(port int) *App {
	a.conf.Server.Port = port
	return a
}
func (a *App) WithLimit(limit int) *App {
	a.conf.Message.MaxStoreSize = limit
	return a
}

func (a *App) WithPushCount(count int) *App {
	a.conf.Message.PushCount = count
	return a
}

func (a *App) Run() {
	container := app.Container()
	system.SecurePanic(a.loadConfig(container))

	infrastructure.InjectStore(container, a.conf.Server.Store)
	application.Inject(container)
	presentation.Inject(container)

	err := container.Invoke(a.serve)
	system.SecurePanic(err)

	system.Shutdown(func() {
		log.Println("server shutdown")
	})
}

func (a *App) loadConfig(container *dig.Container) error {
	if a.configFile != "" {
		if err := a.conf.LoadFile(a.configFile); err != nil {
			return err
		}
	}

	return container.Provide(func() *config.Config{
		return a.conf
	})
}

func (a *App) serve(socket *presentation.Socket) error {
	return socket.Start()

}

func NewApp() *App {
	return &App{
		configFile: "",
		conf: config.New(),
	}
}