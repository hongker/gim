//go:build windows && cgo

package poller

import (
	"gim/framework/poller/wepoll"
)

type epoll = wepoll.Epoll

func NewPollerWithBuffer(size int) (Poller, error) {
	return wepoll.NewPollerWithBuffer(size)
}
