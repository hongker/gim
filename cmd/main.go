package main

import (
	"gim/cmd/app"
	"github.com/ebar-go/ego/utils/runtime"
	"log"
	"os"
)

const (
	name  = "gim"
	usage = "simple and fast im service"
)

func main() {
	// bootstrap with command line
	cmd := app.NewCommand(name, usage)

	// run the command with os.Args.
	runtime.HandleError(cmd.Run(os.Args), func(err error) {
		log.Panic(err)
	})
}
