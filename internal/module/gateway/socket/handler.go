package socket

import (
	"context"
	"gim/internal/module/gateway/domain/dto"
	"github.com/ebar-go/ego/server/ws"
)

type Event func(ctx *ws.Context, proto *Proto) error

func Action[Request any, Response any](fn func(context.Context, *Request) (*Response, error)) Event {
	return func(ctx *ws.Context, proto *Proto) error {
		req := new(Request)
		if err := proto.Bind(req); err != nil {
			return err
		}

		resp, err := fn(ctx, req)
		if err != nil {
			return err
		}

		return proto.Marshal(resp)

	}
}

func ConnectEvent(ctx context.Context, req *dto.ConnectRequest) (resp *dto.ConnectResponse, err error) {
	return
}

func DisconnectEvent(ctx context.Context, req *dto.DisconnectRequest) (resp *dto.DisconnectResponse, err error) {
	return
}

func HeartbeatEvent(ctx context.Context, req *dto.HeartbeatRequest) (resp *dto.HeartbeatResponse, err error) {
	return
}
