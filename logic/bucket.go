package logic

// Bucket represents a bucket for all connections
type Bucket struct {
	sessions map[string]*Session
	channels map[string]*Channel
}

func (bucket *Bucket) Broadcast(msg []byte) {
	for _, session := range bucket.sessions {
		session.Send(msg)
	}
}
func (bucket *Bucket) AddChannel(channel *Channel) {
	bucket.channels[channel.ID] = channel
}
func (bucket *Bucket) RemoveChannel(channel *Channel) {
	delete(bucket.channels, channel.ID)
}
func (bucket *Bucket) GetChannel(id string) *Channel {
	return bucket.channels[id]
}

func (bucket *Bucket) AddSession(session *Session) {
	bucket.sessions[session.ID] = session
}
func (bucket *Bucket) RemoveSession(session *Session) {
	delete(bucket.sessions, session.ID)
}
func (bucket *Bucket) GetSession(id string) *Session {
	return bucket.sessions[id]
}

func (bucket *Bucket) SubscribeChannel(channel *Channel, sessions ...*Session) {
	for _, session := range sessions {
		channel.AddSession(session)
	}
}

func (bucket *Bucket) UnsubscribeChannel(channel *Channel, sessions ...*Session) {
	for _, session := range sessions {
		channel.RemoveSession(session)
	}
}

// Channel represents a channel for many connections
type Channel struct {
	ID       string
	sessions map[string]*Session
}

func (channel *Channel) AddSession(session *Session) {
	channel.sessions[session.ID] = session
}
func (channel *Channel) RemoveSession(session *Session) {
	delete(channel.sessions, session.ID)
}
func (channel *Channel) GetSession(id string) *Session {
	return channel.sessions[id]
}

func (channel *Channel) Broadcast(msg []byte) {
	for _, session := range channel.sessions {
		session.Send(msg)
	}
}
func NewChannel(id string) *Channel {
	return &Channel{}
}

// Session represents a session for one connection
type Session struct {
	ID string
}

func (session *Session) Send(msg []byte) {}
