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

type HandleFunc func(ctx *socket.Context, proto *api.Proto)

type Action[Request, Response any] func(ctx context.Context, req *Request) (*Response, error)

func generic[Request any, Response any](action Action[Request, Response]) HandleFunc {
	return func(ctx *socket.Context, proto *api.Proto) {
		req := new(Request)
		err := runtime.Call(
			// bind with request.
			proto.BindFunc(req),
			// validate request.
			dto.ValidateFunc(req),
			// invoke action.
			func() error {
				validatedCtx, err := NewValidatedContext(ctx)
				if proto.Operate != api.LoginOperate && err != nil {
					return err
				}

				resp, err := action(validatedCtx, req)
				if err != nil {
					return err
				}
				return proto.Marshal(api.NewSuccessResponse(resp))
			})

		runtime.HandleError(err, func(err error) {
			_ = proto.Marshal(api.NewFailureResponse(err))
		})
		return

	}
}

type EventManager struct {
	once     sync.Once
	handlers map[api.OperateType]HandleFunc

	userApp     application.UserApplication
	cometApp    application.CometApplication
	messageApp  application.MessageApplication
	chatroomApp application.ChatroomApplication

	heartbeatInterval time.Duration
}

func NewEventManager(heartbeatInterval time.Duration) *EventManager {
	return &EventManager{
		heartbeatInterval: heartbeatInterval,
		handlers:          map[api.OperateType]HandleFunc{},

		userApp:     application.NewUserApplication(),
		cometApp:    application.GetCometApplication(),
		messageApp:  application.NewMessageApplication(),
		chatroomApp: application.NewChatroomApplication(),
	}
}

func (em *EventManager) buildReleaseTimer(callback func()) *time.Timer {
	timer := time.NewTimer(em.heartbeatInterval)
	go func() {
		defer runtime.HandleCrash()
		<-timer.C
		callback()
	}()
	return timer
}

func (em *EventManager) RegisterReleaseTimer(conn socket.Connection) {
	// close the connection if client don't send heartbeat request.
	timer := em.buildReleaseTimer(func() {
		runtime.HandleError(conn.Close(), func(err error) {
			component.Provider().Logger().Errorf("[%s] closed failed: %v", conn.ID(), err)
		})

		em.cometApp.RemoveUserConnection(GetUidFromConnection(conn))

	})

	SetConnectionTimer(conn, timer)
}

func (em *EventManager) leaseReleaseTimer(conn socket.Connection, duration time.Duration) {
	timer := GetTimerFromConnection(conn)
	if timer == nil {
		return
	}

	timer.Reset(duration)
}

func (em *EventManager) Handle(ctx *socket.Context, proto *api.Proto) {
	em.once.Do(em.initialize)

	handler := em.handlers[proto.OperateType()]
	if handler == nil {
		component.Provider().Logger().Errorf("[%s] No handler registered for type %s", ctx.Conn().ID(), proto.OperateType())
		return
	}

	handler(ctx, proto)
}

func (em *EventManager) initialize() {
	em.handlers[api.LoginOperate] = generic[dto.UserLoginRequest, dto.UserLoginResponse](em.Login)
	em.handlers[api.LogoutOperate] = generic[dto.UserLogoutRequest, dto.UserLogoutResponse](em.Logout)
	em.handlers[api.HeartbeatOperate] = generic[dto.SocketHeartbeatRequest, dto.SocketHeartbeatResponse](em.Heartbeat)
	em.handlers[api.MessageSendOperate] = generic[dto.MessageSendRequest, dto.MessageSendResponse](em.SendMessage)
	em.handlers[api.MessageQueryOperate] = generic[dto.MessageQueryRequest, dto.MessageQueryResponse](em.QueryMessage)
	em.handlers[api.SessionListOperate] = generic[dto.SessionQueryRequest, dto.SessionQueryResponse](em.ListSession)
	em.handlers[api.ChatroomJoinOperate] = generic[dto.ChatroomJoinRequest, dto.ChatroomJoinResponse](em.JoinChatroom)
}

func (em *EventManager) Login(ctx context.Context, req *dto.UserLoginRequest) (resp *dto.UserLoginResponse, err error) {
	resp, err = em.userApp.Login(ctx, req)
	if err != nil {
		return
	}

	conn := ConnectionFromContext(ctx)
	SetConnectionUid(conn, req.ID)
	em.cometApp.SetUserConnection(req.ID, conn)
	return
}

func (em *EventManager) Logout(ctx context.Context, req *dto.UserLogoutRequest) (resp *dto.UserLogoutResponse, err error) {
	resp, err = em.userApp.Logout(ctx, req)
	if err == nil {
		em.cometApp.RemoveUserConnection(auth.UserFromContext(ctx))
	}
	return
}

func (em *EventManager) Heartbeat(ctx context.Context, req *dto.SocketHeartbeatRequest) (resp *dto.SocketHeartbeatResponse, err error) {
	resp = &dto.SocketHeartbeatResponse{ServerTime: time.Now().UnixMilli()}

	// lease timer for close connection
	conn := ConnectionFromContext(ctx)
	em.leaseReleaseTimer(conn, em.heartbeatInterval)
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
