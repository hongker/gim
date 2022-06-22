package gate

import (
	"context"
	"gim/api"
	"gim/api/client"
	"gim/api/protocol"
	"gim/api/server"
	"gim/pkg/network"
	"github.com/pkg/errors"
	"log"
	"time"
)

func (s *Server) HandleConnect(conn *network.Connection) {
	log.Println("welcome:", conn.IP())
}

func (s *Server) HandleAuth(ctx *network.Context) {
	// 判断是否已认证过
	ch := s.bucket.GetChannel(ctx.Connection().ID())
	if ch != nil {
		ctx.Next()
		return
	}

	// 解析数据包
	proto := new(protocol.Proto)
	if err := proto.Unpack(ctx.Request().Body()); err != nil {
		ctx.Output(proto.MustPackFromError(1001, err))
		ctx.Abort()
		return
	}
	log.Println("auth:", proto.Op)

	res, err := s.auth(ctx, proto)
	proto.Op += 1
	if err != nil {
		ctx.Output(proto.MustPackFromError(1001, err))
	} else {
		ctx.Output(proto.MustPackSuccess(&client.AuthResponse{Uid: res.Uid}))
	}
	ctx.Abort()
	return
}

func (s *Server) HandleRequest(ctx *network.Context) {
	// 解析数据包
	proto := new(protocol.Proto)
	if err := proto.Unpack(ctx.Request().Body()); err != nil {
		ctx.Output(nil)
		return
	}

	// 判断是否已认证过
	ch := s.bucket.GetChannel(ctx.Connection().ID())

	// 处理其他请求
	response, err := s.handleReceive(ch, proto)
	if err != nil {
		return
	}
	proto.Body = response.Data
	// 响应包的Op都是在请求包的Op上+1
	proto.Op += 1
	msg, err := proto.Pack()
	if err != nil {
		return
	}
	ctx.Output(msg)
}

func (s *Server) HandleDisconnect(conn *network.Connection) {
	log.Println("goodbye")
}

// auth 处理认证请求
func (s *Server) auth(ctx *network.Context, proto *protocol.Proto) (res *server.AuthResponse, err error) {
	if proto.Op != api.OperateAuth {
		err = errors.New("invalid operate")
		return
	}

	authRequest := new(client.AuthRequest)
	if err = proto.Bind(authRequest); err != nil {
		err = errors.WithMessage(err, "invalid request")
		return
	}
	// 验证用户
	res, err = s.logic.Auth(context.Background(), &server.AuthRequest{
		AppId: authRequest.AppId,
		Name:  authRequest.Name,
	})
	if err != nil {
		return
	}

	// 验证通过后，将channel加入到bucket,用于服务端广播消息
	s.bucket.AddChannel(NewChannel(res.Uid, ctx.Connection()))
	return
}

func (s *Server) handleReceive(channel *Channel, proto *protocol.Proto) (res *server.ReceiveResponse, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	switch proto.Op {
	case api.OperateGroupJoin: // 加群
		groupJoinRequest := new(client.GroupJoinRequest)
		if err = proto.Bind(groupJoinRequest); err != nil {
			err = errors.WithMessage(err, "invalid request")
			return
		}
		s.bucket.GetRoom(groupJoinRequest.GroupId).Add(channel)
	case api.OperateGroupQuit: // 退群
		groupQuitRequest := new(client.GroupQuitRequest)
		if err = proto.Bind(groupQuitRequest); err != nil {
			err = errors.WithMessage(err, "invalid request")
			return
		}
		s.bucket.GetRoom(groupQuitRequest.GroupId).Remove(channel)
	}

	log.Println("receive:", proto.Op)
	return s.logic.Receive(ctx, &server.ReceiveRequest{
		Mid:   channel.key,
		Proto: proto,
	})
}
