package event

import (
	"gim/internal/domain/dto"
	"gim/pkg/network"
)

const (
	Connect    = "Connect"
	Login      = "Login"
	Heartbeat  = "Heartbeat"
	Disconnect = "Disconnect"
	JoinGroup  = "JoinGroup"
	LeaveGroup = "LeaveGroup"
	Push       = "Push"
)

type ConnectEvent struct {
	Connection *network.Connection
}

type LoginEvent struct {
	UserId     string
	Connection *network.Connection
}

type DisconnectEvent struct {
	Connection *network.Connection
}

type JoinGroupEvent struct {
	GroupId    string
	Connection *network.Connection
}
type LeaveGroupEvent struct {
	GroupId    string
	Connection *network.Connection
}
type PushMessageEvent struct {
	SessionType  string
	TargetId     string
	BatchMessage dto.BatchMessage
}
