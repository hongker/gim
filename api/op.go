package api

const (
	OperateAuth      = 101
	OperateAuthReply = 102

	OperateMessageSend      = 201
	OperateMessageSendReply = 202

	OperateMessageQuery      = 203
	OperateMessageQueryReply = 204

	OperateSessionList      = 301
	OperateSessionListReply = 302

	OperateGroupJoin      = 401
	OperateGroupJoinReply = 402
	OperateGroupLeave     = 403
	OperateGroupLeaveReply = 404

	OperateMessagePush = 501
)

const(
	PrivateMessage = "private"
	RoomMessage = "room"
)

const (
	TextMessage = "text"
)