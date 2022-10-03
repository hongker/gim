package main

import (
	"gim/cmd/app"
	"github.com/ebar-go/ego/utils/runtime"
	"log"
	"os"
)

const (
	AppName = "gim"
)

func main() {
	// bootstrap with command line
	cmd := app.NewCommand(AppName)

	// run the command with os.Args.
	runtime.HandleError(cmd.Run(os.Args), func(err error) {
		log.Panic(err)
	})
}
