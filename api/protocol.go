package api

import "encoding/json"

type Packet struct {
	Op int
	Data []byte
}


func (p *Packet) Decode(body []byte) (err error) {
	return
}

func (p *Packet) Encode() (body []byte, err error) {
	return
}

func (p *Packet) Bind(container interface{}) (err error) {
	return json.Unmarshal(p.Data, container)
}

func (p *Packet) Marshal(container interface{}) (err error) {
	p.Data, err = json.Marshal(container)
	return
}

func NewPacket() *Packet {
	return new(Packet)
}