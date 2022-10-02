package api

import (
	"encoding/json"
)

const (
	CodecJson     = "json"
	CodecProtobuf = "protobuf"
)

type Codec interface {
	Decode([]byte, *Proto) error
	Encode(proto *Proto) []byte
}

var (
	DefaultCodec = JsonCodecInstance
)

func JsonCodecInstance() Codec {
	return &JsonCodec{}
}

type JsonCodec struct{}

func (c JsonCodec) Decode(bytes []byte, p *Proto) error {
	return json.Unmarshal(bytes, p)
}

func (c JsonCodec) Encode(proto *Proto) []byte {
	b, _ := json.Marshal(proto)
	return b
}

type ProtobufCodec struct{}

func (c ProtobufCodec) Encode(p *Proto) []byte {
	panic("implement me")
}
func (c ProtobufCodec) Decode(bytes []byte, p *Proto) error {
	panic("implement me")
}
