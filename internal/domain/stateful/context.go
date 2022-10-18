package stateful

import (
	"context"
	"gim/framework"
	"time"
)

const (
	uidParam        = "uid"
	connectionParam = "connection"
	timerParam      = "timer"
)

// NewUserContext return context.Context with uid value
func NewUserContext(ctx context.Context, uid string) context.Context {
	return context.WithValue(ctx, uidParam, uid)
}

// UserFromContext returns user from context
func UserFromContext(ctx context.Context) string {
	return ctx.Value(uidParam).(string)
}

func GetUidFromConnection(conn *framework.Connection) string {
	return conn.Property().GetString(uidParam)
}

func SetConnectionUid(conn *framework.Connection, uid string) {
	conn.Property().Set(uidParam, uid)
}

func GetTimerFromConnection(conn *framework.Connection) *time.Timer {
	if conn == nil {
		return nil
	}
	t := conn.Property().Get(timerParam)
	if t == nil {
		return nil
	}
	timer, _ := t.(*time.Timer)
	return timer
}

func SetConnectionTimer(conn *framework.Connection, timer *time.Timer) {
	conn.Property().Set(timerParam, timer)
}
