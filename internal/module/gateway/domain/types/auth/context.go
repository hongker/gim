package auth

import "context"

const (
	CurrentUserParam = "currentUser"
)

func NewUserContext(ctx context.Context, uid string) context.Context {
	return context.WithValue(ctx, CurrentUserParam, uid)
}

func UserFromContext(ctx context.Context) string {
	return ctx.Value(CurrentUserParam).(string)
}
