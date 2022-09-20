package socket

import (
	"context"
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

		resp, err := fn(ctx, req)
		if err != nil {
			return
		}

		err = proto.Marshal(render.Response{
			Code: 0,
			Msg:  "",
			Data: resp,
		})
		return

	}
}

func LoginEvent(ctx context.Context, req *dto.SocketLoginRequest) (resp *dto.SocketLoginResponse, err error) {
	resp = &dto.SocketLoginResponse{}
	return
}

func HeartbeatEvent(ctx context.Context, req *dto.SocketHeartbeatRequest) (resp *dto.SocketHeartbeatResponse, err error) {
	resp = &dto.SocketHeartbeatResponse{ServerTime: time.Now().UnixMilli()}
	return
}
