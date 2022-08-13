package main

import (
	"fmt"
	"gim/internal"
	"gim/pkg/system"
	"time"
)

func main() {
	go func() {
		for {
			time.Sleep(time.Second * 5)
			fmt.Printf("memory usage: %.2fM\n", float64(system.GetMem())/1000/1000)
		}
	}()
	internal.Run()
}
