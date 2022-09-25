package types

type Chatroom struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Creator   string `json:"creator"`
	CreatedAt int64  `json:"created_at"`
}
