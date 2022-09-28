package entity

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Message struct {
	Id        string `json:"id"`
	SenderId  string `json:"sender_id"`
	Content   string `json:"content"`
	Category  string `json:"category"`
	Status    int    `json:"status"`
	CreatedAt int64  `json:"created_at"`
}

type Chatroom struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Creator   string `json:"creator"`
	CreatedAt int64  `json:"created_at"`
}

type Session struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}
