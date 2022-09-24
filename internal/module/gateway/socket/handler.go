package socket

import (
	"context"
	"gim/internal/module/gateway/application"
	"gim/internal/module/gateway/domain/dto"
	"gim/internal/module/gateway/domain/types/auth"
	"gim/internal/module/gateway/render"
	"github.com/ebar-go/ego/errors"
	"github.com/ebar-go/ego/server/ws"
	"github.com/ebar-go/ego/utils/runtime"
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
		req := new(Request)
		err := runtime.Call(func() error {
			return proto.Bind(req)
		}, func() error {
			return dto.Validate(req)
		}, func() error {
			validatedCtx, err := NewValidatedContext(ctx)
			if proto.Operate != LoginOperate && err != nil {
				return err
			}

			resp, err := fn(validatedCtx, req)
			if err != nil {
				return err
			}
			return proto.Marshal(render.SuccessResponse(resp))
		})

		runtime.HandlerError(err, func(err error) {
			_ = proto.Marshal(render.ErrorResponse(err))
		})
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
