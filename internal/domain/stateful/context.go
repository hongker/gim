package stateful

import (
	"context"
	"gim/framework"
	"github.com/ebar-go/ego/server/socket"
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

// ConnectionFromContext returns connection from context
func ConnectionFromContext(ctx context.Context) socket.Connection {
	return ctx.Value(connectionParam).(socket.Connection)
}

// NewConnectionContext returns context.Context with socket.Connection
func NewConnectionContext(ctx context.Context, conn socket.Connection) context.Context {
	return context.WithValue(ctx, connectionParam, conn)
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
