package gateway

import "encoding/json"

type Codec interface {
	Decode([]byte) (*Proto, error)
	Encode(proto *Proto) []byte
}

var (
	DefaultCodec = JsonCodecInstance
)

func JsonCodecInstance() Codec {
	return &JsonCodec{}
}

type JsonCodec struct{}

func (c JsonCodec) Decode(bytes []byte) (*Proto, error) {
	p := new(Proto)
	err := json.Unmarshal(bytes, p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (c JsonCodec) Encode(proto *Proto) []byte {
	b, _ := json.Marshal(proto)
	return b
}
