package system

import (
	"log"
	"runtime"
	"time"
)

func ShowMemoryUsage() {
	for {
		time.Sleep(time.Second * 5)
		log.Printf("memory usage: %.2fM\n", float64(GetMem())/1000/1000)
	}
}

func GetMem() uint64 {
	var memStat runtime.MemStats
	runtime.ReadMemStats(&memStat)
	return memStat.Sys
}
