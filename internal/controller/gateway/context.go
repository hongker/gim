package gateway

import (
	"context"
	"gim/internal/domain/types/auth"
	"github.com/ebar-go/ego/errors"
	"github.com/ebar-go/ego/server/ws"
)

const (
	UidParam        = "uid"
	ConnectionParam = "connection"
)

func ConnectionFromContext(ctx context.Context) ws.Conn {
	return ctx.Value(ConnectionParam).(ws.Conn)
}
func NewConnectionContext(ctx context.Context, conn ws.Conn) context.Context {
	return context.WithValue(ctx, ConnectionParam, conn)
}

func NewValidatedContext(ctx *ws.Context) (context.Context, error) {
	uid := ctx.Conn().Property().GetString(UidParam)
	connCtx := NewConnectionContext(ctx, ctx.Conn())
	if uid == "" {
		return connCtx, errors.Unauthorized("login required")
	}
	return auth.NewUserContext(connCtx, uid), nil
}
