package dto

type MessageHistoryQuery struct {
	SessionId string
	Limit int
	Last int64
}

type BatchMessage struct {
	Items []Message `json:"items"`
}