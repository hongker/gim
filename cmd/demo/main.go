package main

import "gim/internal"

func main() {
	config := &internal.Config{}

	server := config.Complete().New()

	server.Run()
}
