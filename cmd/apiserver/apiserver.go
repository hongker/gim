package main

import (
	"gim/cmd/apiserver/app"
	"log"
	"os"
)

func main() {
	cmd := app.NewServerCommand()
	err := cmd.Run(os.Args)
	if err != nil {
		log.Fatal("failed to run server command: ", err)
	}
}
