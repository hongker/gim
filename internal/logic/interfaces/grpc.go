package interfaces

import (
	"context"
	"gim/api"
	"gim/api/protocol"
	"gim/api/server"
	"gim/internal/logic/application"
	"gim/pkg/grpc"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/protobuf/proto"
)

type Server struct {
	userApp    *application.UserApp
	messageApp *application.MessageApp
	handlers   map[int32]HandleFunc // 各个操作的处理函数，key为操作类型
}

type HandleFunc func(ctx context.Context, mid string, proto *protocol.Proto) (proto.Message, error)

func NewServer(userApp *application.UserApp, messageApp *application.MessageApp) *Server {
	srv := &Server{handlers: map[int32]HandleFunc{}}
	srv.userApp = userApp
	srv.messageApp = messageApp

	// 发送消息
	srv.handlers[api.OperateMessageSend] = srv.messageApp.Send
	// 查询消息
	srv.handlers[api.OperateMessageQuery] = srv.messageApp.Query
	return srv
}

// Start 启动tcp服务
func (srv *Server) Start(conf grpc.ServerConfig) error {
	grpcSrv := grpc.NewServer(conf)
	server.RegisterLogicServer(grpcSrv.Instance(), srv)
	return grpcSrv.Start()
}

// Auth 对客户端进行授权
func (srv *Server) Auth(ctx context.Context, request *server.AuthRequest) (*server.AuthResponse, error) {
	uid := uuid.NewV4().String()
	err := srv.userApp.Auth(ctx, uid, request.Name)
	if err != nil {
		return nil, err
	}
	return &server.AuthResponse{Uid: uid}, nil
}

// Heartbeat 心跳逻辑
func (srv *Server) Heartbeat(ctx context.Context, request *server.HeartbeatRequest) (*server.HeartbeatResponse, error) {
	return &server.HeartbeatResponse{}, nil
}

// Receive 处理其他逻辑
func (srv *Server) Receive(ctx context.Context, request *server.ReceiveRequest) (response *server.ReceiveResponse, err error) {
	handler, ok := srv.handlers[request.Proto.Op]
	if !ok {
		err = errors.New("invalid operate")
		return
	}

	// 通过handler对各个操作进行处理
	result, err := handler(ctx, request.Mid, request.Proto)
	if err != nil {
		err = errors.WithMessage(err, "handle failed")
		return
	}

	// 对结果进行encode
	data, err := proto.Marshal(result)
	if err != nil {
		err = errors.WithMessage(err, "marshal failed")
		return
	}
	response = &server.ReceiveResponse{Data: data}
	return

}
