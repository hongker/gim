package network

import (
	"encoding/json"
	"github.com/golang/protobuf/proto"
)

type Codec interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(p []byte, v interface{}) (error)
}

type JsonCodec struct {}

func (JsonCodec) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func ( JsonCodec) Unmarshal(p []byte, v interface{}) error {
	return json.Unmarshal(p, v)
}

func NewJsonCodec() Codec {
	return &JsonCodec{}
}
type ProtobufCodec struct {}

func ( ProtobufCodec) Marshal(v interface{}) ([]byte, error) {
	return proto.Marshal(v.(proto.Message))
}

func ( ProtobufCodec) Unmarshal(p []byte, v interface{}) error {
	return proto.Unmarshal(p, v.(proto.Message))
}

func NewProtobufCodec() Codec {
	return &ProtobufCodec{}
}
