package dto

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

func NewPacket() *Packet {
	return new(Packet)
}