package entity


type Message struct {
	Id          string
	SessionId   string
	SessionType string
	Content     string
	ContentType string
	CreatedAt   int64
	ClientMsgId string
	Sequence    int64
	FromUser    *User
}

