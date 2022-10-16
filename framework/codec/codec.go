package codec

import (
	"errors"
	"gim/pkg/binary"
)

type Codec interface {
	Pack(packet *Packet, data any) ([]byte, error)
	Unpack(msg []byte) (packet *Packet, err error)
}

type DefaultCodec struct {
	options *Options
}

type Option func(options *Options)

func Default(options ...Option) Codec {
	// |-------------- header ------------- |-------- body --------|
	// |packetLength|operate|contentType|seq|-------- body --------|
	// |     4      |   2   |      2    | 2 |          n           |
	defaultOptions := &Options{
		headerSize:       10,
		packetLengthSize: 4,
		operateSize:      2,
		contentTypeSize:  2,
		seqSize:          2,
	}
	for _, setter := range options {
		setter(defaultOptions)
	}
	return DefaultCodec{options: defaultOptions}
}

type Options struct {
	headerSize       int
	packetLengthSize int
	operateSize      int
	contentTypeSize  int
	seqSize          int
}

func (codec DefaultCodec) Pack(packet *Packet, data any) ([]byte, error) {
	body, err := packet.Marshal(data)
	if err != nil {
		return nil, err
	}

	// packet header and body
	length := len(body) + codec.options.headerSize
	buf := make([]byte, length)

	packetLengthOffset := 0 + codec.options.packetLengthSize
	operateOffset := packetLengthOffset + codec.options.operateSize
	contentTypeOffset := operateOffset + codec.options.contentTypeSize
	seqOffset := contentTypeOffset + codec.options.seqSize

	binary.BigEndian.PutInt32(buf[0:packetLengthOffset], int32(length))
	binary.BigEndian.PutInt16(buf[packetLengthOffset:operateOffset], packet.Operate)
	binary.BigEndian.PutInt16(buf[operateOffset:contentTypeOffset], packet.ContentType)
	binary.BigEndian.PutInt16(buf[contentTypeOffset:seqOffset], packet.Seq)
	binary.BigEndian.PutString(buf[seqOffset:], string(body))
	return buf, nil
}

func (codec DefaultCodec) Unpack(msg []byte) (*Packet, error) {
	if len(msg) < codec.options.headerSize {
		return nil, errors.New("unexpected message")
	}

	packetLengthOffset := 0 + codec.options.packetLengthSize
	operateOffset := packetLengthOffset + codec.options.operateSize
	contentTypeOffset := operateOffset + codec.options.contentTypeSize
	seqOffset := contentTypeOffset + codec.options.seqSize

	packet := &Packet{}
	length := int(binary.BigEndian.Int32(msg[:packetLengthOffset]))
	packet.Operate = binary.BigEndian.Int16(msg[packetLengthOffset:operateOffset])
	packet.ContentType = binary.BigEndian.Int16(msg[operateOffset:contentTypeOffset])
	packet.Seq = binary.BigEndian.Int16(msg[contentTypeOffset:seqOffset])

	if length > seqOffset {
		packet.Body = msg[seqOffset:length]
	}

	return packet, nil
}
