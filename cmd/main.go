package main

import (
	"gim/cmd/app"
	"log"
	"os"
)

func main() {
	// bootstrap with command line
	cmd := app.NewCommand("gim")

	// run the command with os.Args.
	err := cmd.Run(os.Args)

	if err != nil {
		log.Fatal("run failed:", err)
	}
}
