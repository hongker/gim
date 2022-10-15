package framework

type Codec interface {
	Pack(operate int, data any) ([]byte, error)
	Unpack(msg []byte) (operate int)
	Serializer() Serializer
}

func NewJsonCodec() Codec {
	return nil
}

func NewProtobufCodec() Codec {
	return nil
}
