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
	OperateGroupQuit      = 403
	OperateGroupQuitReply = 404
)

const(
	PrivateMessage = "private"
	RoomMessage = "room"
)