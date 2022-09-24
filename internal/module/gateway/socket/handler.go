package socket

import (
	"context"
	"gim/internal/module/gateway/application"
	"gim/internal/module/gateway/domain/dto"
	"gim/internal/module/gateway/domain/types/auth"
	"gim/internal/module/gateway/render"
	"github.com/ebar-go/ego/errors"
	"github.com/ebar-go/ego/server/ws"
	"time"
)

const (
	LoginUserParam = "currentUser"
)

type Event func(ctx *ws.Context, proto *Proto)

func SetContextProperty(ctx context.Context, key string, value any) {
	wc, ok := ctx.(*ws.Context)
	if !ok {
		return
	}
	wc.Conn().Property().Set(key, value)
}
func NewValidatedContext(ctx *ws.Context) (context.Context, error) {
	uid := ctx.Conn().Property().GetString(LoginUserParam)
	if uid == "" {
		return ctx, errors.Unauthorized("login required")
	}
	return auth.NewUserContext(ctx, uid), nil
}
func Action[Request any, Response any](fn func(context.Context, *Request) (*Response, error)) Event {
	return func(ctx *ws.Context, proto *Proto) {
		var err error
		defer func() {
			if err != nil {
				_ = proto.Marshal(render.ErrorResponse(err))
			}
		}()
		req := new(Request)
		if err = proto.Bind(req); err != nil {
			return
		}

		if err = dto.Validate(req); err != nil {
			return
		}

		validatedCtx, err := NewValidatedContext(ctx)
		if proto.Operate != LoginOperate && err != nil {
			return
		}

		resp, err := fn(validatedCtx, req)
		if err != nil {
			return
		}
		err = proto.Marshal(render.SuccessResponse(resp))

		return

	}
}

type EventManager struct {
	userApp application.UserApplication
}

func (em EventManager) Login(ctx context.Context, req *dto.UserLoginRequest) (resp *dto.UserLoginResponse, err error) {
	resp, err = em.userApp.Login(ctx, req)
	SetContextProperty(ctx, LoginUserParam, req.ID)
	return
}

func (em EventManager) Logout(ctx context.Context, req *dto.UserLogoutRequest) (resp *dto.UserLogoutResponse, err error) {
	resp, err = em.userApp.Logout(ctx, req)
	return
}

func (em EventManager) Heartbeat(ctx context.Context, req *dto.SocketHeartbeatRequest) (resp *dto.SocketHeartbeatResponse, err error) {
	resp = &dto.SocketHeartbeatResponse{ServerTime: time.Now().UnixMilli()}
	return
}

func (em EventManager) FindUser(ctx context.Context, req *dto.UserFindRequest) (resp *dto.UserFindResponse, err error) {
	return em.userApp.Find(ctx, req)
}
