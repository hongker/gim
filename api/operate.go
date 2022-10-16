package api

type OperateType int

const (
	LoginOperate      = 101
	LoginOperateReply = 102

	LogoutOperate                  = 103
	LogoutOperateReply OperateType = 104

	HeartbeatOperate      = 105
	HeartbeatOperateReply = 106

	MessageSendOperate                  = 201
	MessageSendOperateReply OperateType = 202

	MessageQueryOperate                  = 203
	MessageQueryOperateReply OperateType = 204

	MessagePushOperate                  = 205
	MessagePushOperateReply OperateType = 206

	SessionListOperate                  = 205
	SessionListOperateReply OperateType = 206

	ChatroomJoinOperate      = 301
	ChatroomJoinOperateReply = 302
)
