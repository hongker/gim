package socket

import "time"

type Config struct {
	Address           string
	HeartbeatInterval time.Duration
}
