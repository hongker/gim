package framework

import "log"

type Codec interface {
	Pack(operate int, data any) ([]byte, error)
	Unpack(msg []byte) (operate int, err error)
	Serializer() Serializer
}

type EmptyCodec struct{}

func (codec *EmptyCodec) Unpack(msg []byte) (operate int, err error) {
	log.Println("receive:", string(msg))
	return 0, nil
}
func (codec *EmptyCodec) Pack(operate int, data any) ([]byte, error) {
	return nil, nil
}
func (codec *EmptyCodec) Serializer() Serializer {
	return nil
}
func NewJsonCodec() Codec {
	return nil
}

func NewProtobufCodec() Codec {
	return nil
}
