package api

import (
	"encoding/json"
	"fmt"
	"gim/pkg/binary"
	"gim/pkg/bytes"
)

const(
	PacketOffset = 4
	OperateOffset = 8

)
type Packet struct {
	Op int32
	Data []byte
}

// Decode 解析包体
func (p *Packet) Decode(body []byte) (err error) {
	length := binary.BigEndian.Int32(body[:PacketOffset])
	p.Op = binary.BigEndian.Int32(body[PacketOffset:OperateOffset])
	if length > 1000 {
		fmt.Println(string(body))
	}
	p.Data = body[OperateOffset:int(length)]
	return
}

// Encode 组装包体
func (p *Packet) Encode() (buf []byte) {
	buf = bytes.Get(OperateOffset + len(p.Data))
	binary.BigEndian.PutInt32(buf[:PacketOffset], int32(len(buf)))
	binary.BigEndian.PutInt32(buf[PacketOffset:OperateOffset], p.Op)
	binary.BigEndian.PutString(buf[OperateOffset:], string(p.Data))
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
	return &Packet{}
}

func BuildPacket(operate int32, container interface{}) *Packet  {
	p := NewPacket()
	p.Op = operate
	_ = p.Marshal(container)
	return p
}
