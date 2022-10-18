package bucket

// Channel represents a channel for many connections
type Channel struct {
	ID string
	*Room
}

func NewChannel(id string) *Channel {
	return &Channel{ID: id, Room: NewRoom()}
}
