package api

import (
	"encoding/json"
	"gim/pkg/binary"
	"gim/pkg/bytes"
)

type Packet struct {
	Op int32
	Data []byte
}

// Decode 解析包体
func (p *Packet) Decode(body []byte) (err error) {
	p.Op = binary.BigEndian.Int32(body[:4])
	p.Data = body[4:]
	return
}

// Encode 组装包体
func (p *Packet) Encode() (buf []byte) {
	buf = bytes.Get(8 + len(p.Data))
	binary.BigEndian.PutInt32(buf[:4], int32(len(buf)))
	binary.BigEndian.PutInt32(buf[4:8], p.Op)
	copy(buf[8:], p.Data)
	return
}

// Bind 解析数据体
func (p *Packet) Bind(container interface{}) (err error) {
	return json.Unmarshal(p.Data, container)
}

// Marshal 格式化数据
func (p *Packet) Marshal(container interface{}) (err error) {
	p.Data, err = json.Marshal(container)
	return
}

func NewPacket() *Packet {
	return new(Packet)
}