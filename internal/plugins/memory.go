package plugins

import "context"

type StorePlugin interface {
	Save(ctx context.Context)
}

type MemoryStore struct {

}
