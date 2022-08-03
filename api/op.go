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
	UserSession = "user" // session of user to user
	GroupSession = "group" // session of group
)

const (
	TextMessage = "text"
)