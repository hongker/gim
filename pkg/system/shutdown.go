package system

import (
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

// shutdown 关闭服务
func Shutdown(callback func()) {
	// signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		sig := <-c
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			if callback != nil {
				callback()
			}

			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
func GetMem() uint64 {
	var memStat runtime.MemStats
	runtime.ReadMemStats(&memStat)
	return memStat.Sys
}
