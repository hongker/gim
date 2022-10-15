package pool

import (
	"bytes"
	"sync"
)

type BufferPool sync.Pool

func NewBufferPool() *BufferPool {
	return &BufferPool{
		New: func() interface{} {
			return &bytes.Buffer{}
		},
	}
}
func (pool *BufferPool) Get() *bytes.Buffer {
	return (*sync.Pool)(pool).Get().(*bytes.Buffer)
}

func (pool *BufferPool) Put(buffer *bytes.Buffer) {
	(*sync.Pool)(pool).Put(buffer)
}
