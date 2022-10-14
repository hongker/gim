package framework

type Codec interface {
	Pack()
	Unpack()
}

func NewJsonCodec() Codec {
	return nil
}

func NewProtobufCodec() Codec {
	return nil
}
