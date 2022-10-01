package gateway

import (
	"encoding/json"
	"gim/api"
)

const (
	CodecJson     = "json"
	CodecProtobuf = "protobuf"
)

type Codec interface {
	Decode([]byte, *api.Proto) error
	Encode(proto *api.Proto) []byte
}

var (
	DefaultCodec = JsonCodecInstance
)

func JsonCodecInstance() Codec {
	return &JsonCodec{}
}

type JsonCodec struct{}

func (c JsonCodec) Decode(bytes []byte, p *api.Proto) error {
	return json.Unmarshal(bytes, p)
}

func (c JsonCodec) Encode(proto *api.Proto) []byte {
	b, _ := json.Marshal(proto)
	return b
}

type ProtobufCodec struct{}

func (c ProtobufCodec) Encode(p *api.Proto) []byte {
	panic("implement me")
}
func (c ProtobufCodec) Decode(bytes []byte, p *api.Proto) error {
	panic("implement me")
}
