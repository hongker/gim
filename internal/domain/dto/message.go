package dto

type MessageHistoryQuery struct {
	SessionId string
	Limit     int
	Last      int64
}

type BatchMessage struct {
	Count int       `json:"count"`
	Items []Message `json:"items"`
}

type SessionMessagePackage struct {
	Session Session   `json:"session"`
	Items   []Message `json:"items"`
}
