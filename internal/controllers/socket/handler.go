package socket

import (
	"context"
	"gim/internal/application"
	"gim/internal/domain/dto"
	"gim/internal/domain/types/auth"
	"gim/internal/infrastructure/render"
	"github.com/ebar-go/ego/errors"
	"github.com/ebar-go/ego/server/ws"
	"github.com/ebar-go/ego/utils/runtime"
	"time"
)

const (
	UidParam        = "uid"
	ConnectionParam = "connection"
)

type Event func(ctx *ws.Context, proto *Proto)

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
	userApp     application.UserApplication
	cometApp    application.CometApplication
	messageApp  application.MessageApplication
	chatroomApp application.ChatroomApplication
}

func NewEventManager() *EventManager {
	return &EventManager{
		userApp:     application.NewUserApplication(),
		cometApp:    application.GetCometApplication(),
		messageApp:  application.NewMessageApplication(),
		chatroomApp: application.NewChatroomApplication(),
	}
}

func (em EventManager) Login(ctx context.Context, req *dto.UserLoginRequest) (resp *dto.UserLoginResponse, err error) {
	resp, err = em.userApp.Login(ctx, req)
	if err != nil {
		return
	}

	conn := ConnectionFromContext(ctx)
	conn.Property().Set(UidParam, req.ID)
	em.cometApp.SetUserConnection(req.ID, conn)
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

func (em EventManager) SendMessage(ctx context.Context, req *dto.MessageSendRequest) (resp *dto.MessageSendResponse, err error) {
	uid := auth.UserFromContext(ctx)
	return em.messageApp.Send(ctx, uid, req)
}

func (em EventManager) JoinChatroom(ctx context.Context, req *dto.ChatroomJoinRequest) (resp *dto.ChatroomJoinResponse, err error) {
	uid := auth.UserFromContext(ctx)
	return em.chatroomApp.Join(ctx, uid, req)
}
