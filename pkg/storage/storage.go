package storage

import "gim/pkg/runtime"

type Interface interface {
	Save(obj runtime.Object) error
	Remove(obj runtime.Object) error
}

type MemoryStorage struct{}
type RedisStorage struct{}
type EtcdStorage struct{}
