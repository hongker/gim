package storage

import (
	"context"
)

type Object interface {
	ID() string
}
type Storage interface {
	Save(ctx context.Context, object Object) error
	Find(ctx context.Context, object Object) error
}

type List interface {
	Push(ctx context.Context, object Object) error
	Pop(ctx context.Context, object Object) error
}

type String interface {
	Set(ctx context.Context)
}
