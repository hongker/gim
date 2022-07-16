package entity

type Room struct {
	id       string
	channels map[string]*Channel
}

func (room *Room) Channels() map[string]*Channel {
	return room.channels
}

func NewRoom(id string) *Room {
	return &Room{
		id:       id,
		channels: make(map[string]*Channel),
	}
}

func (room *Room) Add(channel *Channel) {
	room.channels[channel.key] = channel
}
func (room *Room) Remove(channel *Channel) {
	delete(room.channels, channel.key)
}

func (room *Room) Push(packet []byte) {

}
