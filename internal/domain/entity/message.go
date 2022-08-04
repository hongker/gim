package entity


type Message struct {
	Id          string
	SessionId   string
	SessionType string
	Content     string
	ContentType string
	CreatedAt   int64
	RequestId string
	Sequence    int64
	FromUser    *User
}

