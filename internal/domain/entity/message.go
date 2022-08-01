package entity


type Message struct {
	Id          string
	Type        string
	Content     string
	CreatedAt        int64
	ClientMsgId string
	Sequence    int64
	SessionId   string
	FromUser    *User
}

