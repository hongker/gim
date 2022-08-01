package dto

type MessageHistoryQuery struct {
	SessionId string
	Limit int
	Last int64
}
