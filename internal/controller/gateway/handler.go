package gateway

import (
	"context"
	"gim/api"
	"gim/internal/application"
	"gim/internal/domain/dto"
	"gim/internal/domain/types/auth"
	"github.com/ebar-go/ego/component"
	"github.com/ebar-go/ego/server/socket"
	"github.com/ebar-go/ego/utils/runtime"
	"sync"
	"time"
)

type HandleFunc func(ctx *socket.Context, proto *Proto)

type Action[Request, Response any] func(ctx context.Context, req *Request) (*Response, error)

func generic[Request any, Response any](action Action[Request, Response]) HandleFunc {
	return func(ctx *socket.Context, proto *Proto) {
		req := new(Request)
		err := runtime.Call(
			// bind with request.
			proto.BindFunc(req),
			// validate request.
			dto.ValidateFunc(req),
			// invoke action.
			func() error {
				validatedCtx, err := NewValidatedContext(ctx)
				if proto.Operate != LoginOperate && err != nil {
					return err
				}

				resp, err := action(validatedCtx, req)
				if err != nil {
					return err
				}
				return proto.Marshal(api.NewSuccessResponse(resp))
			})

		runtime.HandlerError(err, func(err error) {
			_ = proto.Marshal(api.NewFailureResponse(err))
		})
		return

	}
}

type EventManager struct {
	once     sync.Once
	handlers map[OperateType]HandleFunc

	userApp     application.UserApplication
	cometApp    application.CometApplication
	messageApp  application.MessageApplication
	chatroomApp application.ChatroomApplication
}

func NewEventManager() *EventManager {
	return &EventManager{
		handlers: map[OperateType]HandleFunc{},

		userApp:     application.NewUserApplication(),
		cometApp:    application.GetCometApplication(),
		messageApp:  application.NewMessageApplication(),
		chatroomApp: application.NewChatroomApplication(),
	}
}

func (em *EventManager) Handle(ctx *socket.Context, proto *Proto) {
	em.once.Do(em.initialize)

	handler := em.handlers[proto.OperateType()]
	if handler == nil {
		component.Provider().Logger().Errorf("[%s] No handler registered for type %s", ctx.Conn().ID(), proto.OperateType())
		return
	}

	handler(ctx, proto)
}

func (em *EventManager) initialize() {
	em.handlers[LoginOperate] = generic[dto.UserLoginRequest, dto.UserLoginResponse](em.Login)
	em.handlers[LogoutOperate] = generic[dto.UserLogoutRequest, dto.UserLogoutResponse](em.Logout)
	em.handlers[HeartbeatOperate] = generic[dto.SocketHeartbeatRequest, dto.SocketHeartbeatResponse](em.Heartbeat)
	em.handlers[MessageSendOperate] = generic[dto.MessageSendRequest, dto.MessageSendResponse](em.SendMessage)
	em.handlers[MessageQueryOperate] = generic[dto.MessageQueryRequest, dto.MessageQueryResponse](em.QueryMessage)
	em.handlers[SessionListOperate] = generic[dto.SessionQueryRequest, dto.SessionQueryResponse](em.ListSession)
	em.handlers[ChatroomJoinOperate] = generic[dto.ChatroomJoinRequest, dto.ChatroomJoinResponse](em.JoinChatroom)
}

func (em *EventManager) Login(ctx context.Context, req *dto.UserLoginRequest) (resp *dto.UserLoginResponse, err error) {
	resp, err = em.userApp.Login(ctx, req)
	if err != nil {
		return
	}

	conn := ConnectionFromContext(ctx)
	conn.Property().Set(UidParam, req.ID)
	em.cometApp.SetUserConnection(req.ID, conn)
	return
}

func (em *EventManager) Logout(ctx context.Context, req *dto.UserLogoutRequest) (resp *dto.UserLogoutResponse, err error) {
	resp, err = em.userApp.Logout(ctx, req)
	return
}

func (em *EventManager) Heartbeat(ctx context.Context, req *dto.SocketHeartbeatRequest) (resp *dto.SocketHeartbeatResponse, err error) {
	resp = &dto.SocketHeartbeatResponse{ServerTime: time.Now().UnixMilli()}

	conn := ConnectionFromContext(ctx)
	property := conn.Property().Get("timer")
	if property == nil {
		return
	}
	timer, ok := property.(*time.Timer)
	if !ok {
		return
	}
	timer.Reset(time.Minute)
	return
}

func (em *EventManager) FindUser(ctx context.Context, req *dto.UserFindRequest) (resp *dto.UserFindResponse, err error) {
	return em.userApp.Find(ctx, req)
}

func (em *EventManager) SendMessage(ctx context.Context, req *dto.MessageSendRequest) (resp *dto.MessageSendResponse, err error) {
	uid := auth.UserFromContext(ctx)
	return em.messageApp.Send(ctx, uid, req)
}

func (em *EventManager) JoinChatroom(ctx context.Context, req *dto.ChatroomJoinRequest) (resp *dto.ChatroomJoinResponse, err error) {
	uid := auth.UserFromContext(ctx)
	return em.chatroomApp.Join(ctx, uid, req)
}

func (em *EventManager) QueryMessage(ctx context.Context, req *dto.MessageQueryRequest) (resp *dto.MessageQueryResponse, err error) {
	return em.messageApp.Query(ctx, req)
}
func (em *EventManager) ListSession(ctx context.Context, req *dto.SessionQueryRequest) (resp *dto.SessionQueryResponse, err error) {
	uid := auth.UserFromContext(ctx)
	return em.messageApp.ListSession(ctx, uid, req)
}
