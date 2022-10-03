package gateway

import (
	"context"
	"gim/api"
	"gim/internal/application"
	"gim/internal/domain/dto"
	"gim/internal/domain/stateful"
	"github.com/ebar-go/ego/component"
	"github.com/ebar-go/ego/errors"
	"github.com/ebar-go/ego/server/socket"
	"github.com/ebar-go/ego/utils/runtime"
	"sync"
	"time"
)

type HandleFunc func(ctx *socket.Context, proto *api.Proto)

type Action[Request, Response any] func(ctx context.Context, req *Request) (*Response, error)

// newValidatedContext returns a new context
// if user is not authenticated, context only include connection param.
// if user is authenticated, context include uid param and connection param.
func newValidatedContext(ctx *socket.Context) (context.Context, error) {
	uid := stateful.GetUidFromConnection(ctx.Conn())
	connCtx := stateful.NewConnectionContext(ctx, ctx.Conn())
	if uid == "" {
		return connCtx, errors.Unauthorized("login required")
	}
	return stateful.NewUserContext(connCtx, uid), nil
}

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
				validatedCtx, err := newValidatedContext(ctx)
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

type Handler struct {
	once   sync.Once
	routes map[api.OperateType]HandleFunc

	userApp     application.UserApplication
	cometApp    application.CometApplication
	messageApp  application.MessageApplication
	chatroomApp application.ChatroomApplication

	heartbeatInterval time.Duration
}

func NewHandler(heartbeatInterval time.Duration) *Handler {
	return &Handler{
		heartbeatInterval: heartbeatInterval,
		routes:            map[api.OperateType]HandleFunc{},

		userApp:     application.NewUserApplication(),
		cometApp:    application.GetCometApplication(),
		messageApp:  application.NewMessageApplication(),
		chatroomApp: application.NewChatroomApplication(),
	}
}

// InitializeConn initializes connection
func (handler *Handler) InitializeConn(conn socket.Connection) {
	// start release timer
	handler.startReleaseTimer(conn)
}

// FinalizeConn finalizes connection
func (handler *Handler) FinalizeConn(conn socket.Connection) {
	// stop release timer
	handler.stopReleaseTimer(conn)
}

// HandleRequest handle user requests
func (handler *Handler) HandleRequest(ctx *socket.Context, proto *api.Proto) {
	handler.once.Do(handler.initialize)

	fn := handler.routes[proto.OperateType()]
	if fn == nil {
		component.Provider().Logger().Errorf("[%s] No handler registered for type: %s", ctx.Conn().ID(), proto.OperateType())
		return
	}

	fn(ctx, proto)
}

func (handler *Handler) initialize() {
	handler.registerRoutes()
}

func (handler *Handler) buildReleaseTimer(callback func()) *time.Timer {
	timer := time.NewTimer(handler.heartbeatInterval)
	go func() {
		defer runtime.HandleCrash()
		<-timer.C
		callback()
	}()
	return timer
}

// startReleaseTimer
func (handler *Handler) startReleaseTimer(conn socket.Connection) {
	// close the connection if client don't send heartbeat request.
	timer := handler.buildReleaseTimer(func() {
		runtime.HandleError(conn.Close(), func(err error) {
			component.Provider().Logger().Errorf("[%s] closed failed: %v", conn.ID(), err)
		})

		handler.cometApp.RemoveUserConnection(stateful.GetUidFromConnection(conn))

	})

	stateful.SetConnectionTimer(conn, timer)
}

// leaseReleaseTimer
func (handler *Handler) leaseReleaseTimer(conn socket.Connection, duration time.Duration) {
	runtime.HandleNil[time.Timer](stateful.GetTimerFromConnection(conn), func(timer *time.Timer) {
		timer.Reset(duration)
	})
}

// stopReleaseTimer
func (handler *Handler) stopReleaseTimer(conn socket.Connection) {
	runtime.HandleNil[time.Timer](stateful.GetTimerFromConnection(conn), func(timer *time.Timer) {
		timer.Stop()
	})
}

// registerRoutes registers routes
func (handler *Handler) registerRoutes() {
	handler.routes[api.LoginOperate] = generic[dto.UserLoginRequest, dto.UserLoginResponse](handler.Login)
	handler.routes[api.LogoutOperate] = generic[dto.UserLogoutRequest, dto.UserLogoutResponse](handler.Logout)
	handler.routes[api.HeartbeatOperate] = generic[dto.SocketHeartbeatRequest, dto.SocketHeartbeatResponse](handler.Heartbeat)
	handler.routes[api.MessageSendOperate] = generic[dto.MessageSendRequest, dto.MessageSendResponse](handler.SendMessage)
	handler.routes[api.MessageQueryOperate] = generic[dto.MessageQueryRequest, dto.MessageQueryResponse](handler.QueryMessage)
	handler.routes[api.SessionListOperate] = generic[dto.SessionQueryRequest, dto.SessionQueryResponse](handler.ListSession)
	handler.routes[api.ChatroomJoinOperate] = generic[dto.ChatroomJoinRequest, dto.ChatroomJoinResponse](handler.JoinChatroom)
}

// Login handle user login request
func (handler *Handler) Login(ctx context.Context, req *dto.UserLoginRequest) (resp *dto.UserLoginResponse, err error) {
	resp, err = handler.userApp.Login(ctx, req)
	if err != nil {
		return
	}

	conn := stateful.ConnectionFromContext(ctx)
	stateful.SetConnectionUid(conn, req.ID)
	handler.cometApp.SetUserConnection(req.ID, conn)
	return
}

// Logout handle user logout request
func (handler *Handler) Logout(ctx context.Context, req *dto.UserLogoutRequest) (resp *dto.UserLogoutResponse, err error) {
	resp, err = handler.userApp.Logout(ctx, req)

	if err == nil {
		handler.cometApp.RemoveUserConnection(stateful.UserFromContext(ctx))
	}
	return
}

// Heartbeat handle user heartbeat request
func (handler *Handler) Heartbeat(ctx context.Context, req *dto.SocketHeartbeatRequest) (resp *dto.SocketHeartbeatResponse, err error) {
	resp = &dto.SocketHeartbeatResponse{ServerTime: time.Now().UnixMilli()}

	// lease timer for close connection
	conn := stateful.ConnectionFromContext(ctx)
	handler.leaseReleaseTimer(conn, handler.heartbeatInterval)
	return
}

// FindUser handle user query request
func (handler *Handler) FindUser(ctx context.Context, req *dto.UserFindRequest) (resp *dto.UserFindResponse, err error) {
	return handler.userApp.Find(ctx, req)
}

// SendMessage handle user send message request
func (handler *Handler) SendMessage(ctx context.Context, req *dto.MessageSendRequest) (resp *dto.MessageSendResponse, err error) {
	uid := stateful.UserFromContext(ctx)
	return handler.messageApp.Send(ctx, uid, req)
}

// JoinChatroom handle user join chatroom request
func (handler *Handler) JoinChatroom(ctx context.Context, req *dto.ChatroomJoinRequest) (resp *dto.ChatroomJoinResponse, err error) {
	uid := stateful.UserFromContext(ctx)
	return handler.chatroomApp.Join(ctx, uid, req)
}

// QueryMessage handle user query session message request
func (handler *Handler) QueryMessage(ctx context.Context, req *dto.MessageQueryRequest) (resp *dto.MessageQueryResponse, err error) {
	return handler.messageApp.Query(ctx, req)
}

// ListSession handle use list session request
func (handler *Handler) ListSession(ctx context.Context, req *dto.SessionQueryRequest) (resp *dto.SessionQueryResponse, err error) {
	uid := stateful.UserFromContext(ctx)
	return handler.messageApp.ListSession(ctx, uid, req)
}
