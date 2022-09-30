package gateway

import (
	"context"
	"gim/internal/domain/types/auth"
	"github.com/ebar-go/ego/errors"
	"github.com/ebar-go/ego/server/socket"
)

const (
	UidParam        = "uid"
	ConnectionParam = "connection"
)

func ConnectionFromContext(ctx context.Context) socket.Connection {
	return ctx.Value(ConnectionParam).(socket.Connection)
}
func NewConnectionContext(ctx context.Context, conn socket.Connection) context.Context {
	return context.WithValue(ctx, ConnectionParam, conn)
}

func NewValidatedContext(ctx *socket.Context) (context.Context, error) {
	uid := ctx.Conn().Property().GetString(UidParam)
	connCtx := NewConnectionContext(ctx, ctx.Conn())
	if uid == "" {
		return connCtx, errors.Unauthorized("login required")
	}
	return auth.NewUserContext(connCtx, uid), nil
}
