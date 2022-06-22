package entity

import "encoding/json"

type Message struct {
	Id          string
	Type        string
	Content     string
	Time        int64
	ClientMsgId string
	Sequence    int64
	SessionId   string
	FromUser    *User
}

func (item *Message) Marshal() []byte {
	b, _ := json.Marshal(item)
	return b
}

func (item *Message) Unmarshal(source []byte) error {
	return json.Unmarshal(source, item)
}
