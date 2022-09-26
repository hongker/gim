package render

import (
	"encoding/json"
	"sync"
)

type Serializer struct{}

func (s *Serializer) Decode(data []byte, container interface{}) error {
	return json.Unmarshal(data, container)
}

func (s *Serializer) Encode(container interface{}) ([]byte, error) {
	return json.Marshal(container)
}

var serializerOnce struct {
	once       sync.Once
	serializer *Serializer
}

func serializer() *Serializer {
	serializerOnce.once.Do(func() {
		serializerOnce.serializer = &Serializer{}
	})
	return serializerOnce.serializer
}
