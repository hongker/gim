package framework

import (
	"errors"
	"gim/pkg/binary"
)

type Codec interface {
	Pack(packet *Packet, data any) ([]byte, error)
	Unpack(msg []byte) (packet *Packet, err error)
}

type DefaultCodec struct {
	headerSize int
}

func (codec DefaultCodec) Pack(packet *Packet, data any) ([]byte, error) {
	body, err := packet.Marshal(data)
	if err != nil {
		return nil, err
	}

	// packet header and body
	length := len(body) + 10
	buf := make([]byte, length)
	binary.BigEndian.PutInt32(buf[:4], int32(length))
	binary.BigEndian.PutInt32(buf[4:8], packet.Operate)
	binary.BigEndian.PutInt16(buf[8:10], packet.ContentType)
	binary.BigEndian.PutString(buf[10:], string(body))
	return buf, nil
}

func (codec DefaultCodec) Unpack(msg []byte) (*Packet, error) {
	if len(msg) < 10 {
		return nil, errors.New("unexpected message")
	}

	packet := &Packet{}
	length := int(binary.BigEndian.Int32(msg[:4]))
	packet.Operate = binary.BigEndian.Int32(msg[4:8])
	packet.ContentType = binary.BigEndian.Int16(msg[8:10])

	if length > 10 {
		packet.Body = msg[10:length]
	}

	return packet, nil
}
