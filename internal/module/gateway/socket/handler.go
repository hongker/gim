package socket

import (
	"context"
	"gim/internal/module/gateway/application"
	"gim/internal/module/gateway/domain/dto"
	"gim/internal/module/gateway/render"
	"github.com/ebar-go/ego/server/ws"
	"time"
)

type Event func(ctx *ws.Context, proto *Proto)

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

		resp, err := fn(ctx, req)
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
	return
}

func (em EventManager) Heartbeat(ctx context.Context, req *dto.SocketHeartbeatRequest) (resp *dto.SocketHeartbeatResponse, err error) {
	resp = &dto.SocketHeartbeatResponse{ServerTime: time.Now().UnixMilli()}
	return
}
