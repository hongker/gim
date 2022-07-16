package protocol

import (
	"gim/pkg/binary"
	"google.golang.org/protobuf/proto"
)

const (
	PacketOffset  = 4
	VersionOffset = PacketOffset + 2
	OperateOffset = VersionOffset + 4
)

// Unpack 解包
func (p *Proto) Unpack(buf []byte) error {
	p.Ver = int32(binary.BigEndian.Int16(buf[PacketOffset:VersionOffset]))
	p.Op = binary.BigEndian.Int32(buf[VersionOffset:OperateOffset])
	p.Body = buf[OperateOffset:]
	return nil
}

// Bind 解body
func (p *Proto) Bind(message proto.Message) error {
	return proto.Unmarshal(p.Body, message)
}

func (p *Proto) Marshal(message proto.Message) ([]byte, error) {
	body, err := proto.Marshal(message)
	if err != nil {
		return nil, err
	}
	p.Body = body
	return p.Pack()
}

// Pack 打包
func (p *Proto) Pack() ([]byte, error) {
	packetLen := len(p.Body) + OperateOffset
	buf := make([]byte, packetLen)
	binary.BigEndian.PutInt32(buf[0:], int32(packetLen))
	binary.BigEndian.PutInt16(buf[PacketOffset:], int16(p.Ver))
	binary.BigEndian.PutInt32(buf[VersionOffset:], p.Op)
	binary.BigEndian.PutString(buf[OperateOffset:], string(p.Body))
	return buf, nil
}

func (p *Proto) PackFromError(code int32, err error) ([]byte, error) {
	response := &Response{
		Code: code,
		Msg:  err.Error(),
		Data: nil,
	}
	return p.packResponse(response)
}
func (p *Proto) MustPackFromError(code int32, err error) []byte {
	response := &Response{
		Code: code,
		Msg:  err.Error(),
		Data: nil,
	}
	packet, _ := p.packResponse(response)
	return packet
}
func (p *Proto) MustPackSuccess(message proto.Message) []byte {
	response := &Response{
		Code: 0,
		Msg:  "",
	}

	if message != nil {
		data, err := proto.Marshal(message)
		if err != nil {
			return nil
		}

		response.Data = data
	}
	packet, _ := p.packResponse(response)
	return packet
}

func (p *Proto) MustPackSuccessFromBytes(data []byte) []byte {
	response := &Response{
		Code: 0,
		Msg:  "",
		Data: data,
	}
	packet, _ := p.packResponse(response)
	return packet
}

func (p *Proto) packResponse(response *Response) ([]byte, error) {
	body, err := proto.Marshal(response)
	if err != nil {
		return nil, err
	}
	p.Body = body
	p.Op += 1
	return p.Pack()
}
