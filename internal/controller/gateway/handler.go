package gateway

import (
	"gim/api"
	"gim/framework"
	"gim/internal/application"
	"gim/internal/domain/dto"
	"gim/internal/domain/stateful"
	"github.com/ebar-go/ego/component"
	"github.com/ebar-go/ego/utils/runtime"
	"sync"
	"time"
)

type Handler struct {
	once   sync.Once
	routes map[api.OperateType]framework.HandleFunc

	userApp     application.UserApplication
	cometApp    application.CometApplication
	messageApp  application.MessageApplication
	chatroomApp application.ChatroomApplication

	heartbeatInterval time.Duration
}

func NewHandler(heartbeatInterval time.Duration) *Handler {
	return &Handler{
		heartbeatInterval: heartbeatInterval,
		routes:            map[api.OperateType]framework.HandleFunc{},

		userApp:     application.NewUserApplication(),
		cometApp:    application.GetCometApplication(),
		messageApp:  application.NewMessageApplication(),
		chatroomApp: application.NewChatroomApplication(),
	}
}

func (handler *Handler) filter(ctx *framework.Context) {
	if ctx.Operate() != api.LoginOperate {
		// check user login state
		if uid := stateful.GetUidFromConnection(ctx.Conn()); uid == "" {
			ctx.Abort()
		}
	}

	ctx.Next()
}

func (handler *Handler) OnConnect(conn *framework.Connection) {
	component.Provider().Logger().Infof("[%s] Connected", conn.UUID())
	handler.InitializeConn(conn)
}

func (handler *Handler) OnDisconnect(conn *framework.Connection) {
	component.Provider().Logger().Infof("[%s] Disconnected", conn.UUID())
	handler.FinalizeConn(conn)
}

func (handler *Handler) Install(router *framework.Router) {
	router.Route(api.LoginOperate, framework.StandardHandler(handler.Login))
	router.Route(api.HeartbeatOperate, framework.StandardHandler(handler.Heartbeat))
	router.Route(api.LogoutOperate, framework.StandardHandler(handler.Logout))
	router.Route(api.MessageSendOperate, framework.StandardHandler(handler.SendMessage))
	router.Route(api.MessageQueryOperate, framework.StandardHandler(handler.QueryMessage))
	router.Route(api.SessionListOperate, framework.StandardHandler(handler.ListSession))
	router.Route(api.ChatroomJoinOperate, framework.StandardHandler(handler.JoinChatroom))
}

// InitializeConn initializes connection
func (handler *Handler) InitializeConn(conn *framework.Connection) {
	// start release timer
	handler.startReleaseTimer(conn)
}

// FinalizeConn finalizes connection
func (handler *Handler) FinalizeConn(conn *framework.Connection) {
	// stop release timer
	handler.stopReleaseTimer(conn)
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
func (handler *Handler) startReleaseTimer(conn *framework.Connection) {
	// close the connection if client don't send heartbeat request.
	timer := handler.buildReleaseTimer(func() {
		conn.Close()

		handler.cometApp.RemoveUserConnection(stateful.GetUidFromConnection(conn))

	})

	stateful.SetConnectionTimer(conn, timer)
}

// leaseReleaseTimer
func (handler *Handler) leaseReleaseTimer(conn *framework.Connection, duration time.Duration) {
	runtime.HandleNil[time.Timer](stateful.GetTimerFromConnection(conn), func(timer *time.Timer) {
		timer.Reset(duration)
	})
}

// stopReleaseTimer
func (handler *Handler) stopReleaseTimer(conn *framework.Connection) {
	runtime.HandleNil[time.Timer](stateful.GetTimerFromConnection(conn), func(timer *time.Timer) {
		timer.Stop()
	})
}

// Login handle user login request
func (handler *Handler) Login(ctx *framework.Context, req *dto.UserLoginRequest) (resp *dto.UserLoginResponse, err error) {
	resp, err = handler.userApp.Login(ctx, req)
	if err != nil {
		return
	}

	stateful.SetConnectionUid(ctx.Conn(), req.ID)
	handler.cometApp.SetUserConnection(req.ID, ctx.Conn())
	return
}

// Logout handle user logout request

func (handler *Handler) Logout(ctx *framework.Context, req *dto.UserLogoutRequest) (resp *dto.UserLogoutResponse, err error) {
	resp, err = handler.userApp.Logout(ctx, req)

	if err == nil {
		handler.cometApp.RemoveUserConnection(stateful.UserFromContext(ctx))
	}
	return
}

// Heartbeat handle user heartbeat request
func (handler *Handler) Heartbeat(ctx *framework.Context, req *dto.SocketHeartbeatRequest) (resp *dto.SocketHeartbeatResponse, err error) {
	resp = &dto.SocketHeartbeatResponse{ServerTime: time.Now().UnixMilli()}

	// lease timer for close connection
	handler.leaseReleaseTimer(ctx.Conn(), handler.heartbeatInterval)
	return
}

// FindUser handle user query request
func (handler *Handler) FindUser(ctx *framework.Context, req *dto.UserFindRequest) (resp *dto.UserFindResponse, err error) {
	return handler.userApp.Find(ctx, req)
}

// SendMessage handle user send message request
func (handler *Handler) SendMessage(ctx *framework.Context, req *dto.MessageSendRequest) (resp *dto.MessageSendResponse, err error) {
	uid := stateful.UserFromContext(ctx)
	return handler.messageApp.Send(ctx, uid, req)
}

// JoinChatroom handle user join chatroom request
func (handler *Handler) JoinChatroom(ctx *framework.Context, req *dto.ChatroomJoinRequest) (resp *dto.ChatroomJoinResponse, err error) {
	uid := stateful.UserFromContext(ctx)
	return handler.chatroomApp.Join(ctx, uid, req)
}

// QueryMessage handle user query session message request
func (handler *Handler) QueryMessage(ctx *framework.Context, req *dto.MessageQueryRequest) (resp *dto.MessageQueryResponse, err error) {
	return handler.messageApp.Query(ctx, req)
}

// ListSession handle use list session request
func (handler *Handler) ListSession(ctx *framework.Context, req *dto.SessionQueryRequest) (resp *dto.SessionQueryResponse, err error) {
	uid := stateful.UserFromContext(ctx)
	return handler.messageApp.ListSession(ctx, uid, req)
}
