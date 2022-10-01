package gateway

import (
	"context"
	"gim/internal/domain/types/auth"
	"github.com/ebar-go/ego/errors"
	"github.com/ebar-go/ego/server/socket"
	"time"
)

const (
	UidParam        = "uid"
	ConnectionParam = "connection"
	TimerParam      = "timer"
)

func ConnectionFromContext(ctx context.Context) socket.Connection {
	return ctx.Value(ConnectionParam).(socket.Connection)
}
func NewConnectionContext(ctx context.Context, conn socket.Connection) context.Context {
	return context.WithValue(ctx, ConnectionParam, conn)
}

func NewValidatedContext(ctx *socket.Context) (context.Context, error) {
	uid := GetUidFromConnection(ctx.Conn())
	connCtx := NewConnectionContext(ctx, ctx.Conn())
	if uid == "" {
		return connCtx, errors.Unauthorized("login required")
	}
	return auth.NewUserContext(connCtx, uid), nil
}

func GetUidFromConnection(conn socket.Connection) string {
	return conn.Property().GetString(UidParam)
}

func SetConnectionUid(conn socket.Connection, uid string) {
	conn.Property().Set(UidParam, uid)
}

func GetTimerFromConnection(conn socket.Connection) *time.Timer {
	if conn == nil {
		return nil
	}
	t := conn.Property().Get(TimerParam)
	if t == nil {
		return nil
	}
	timer, _ := t.(*time.Timer)
	return timer
}

func SetConnectionTimer(conn socket.Connection, timer *time.Timer) {
	conn.Property().Set(TimerParam, timer)
}
