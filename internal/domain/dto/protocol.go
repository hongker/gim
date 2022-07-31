package dto

type Packet struct {

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