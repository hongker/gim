package main

import (
	"gim/framework"
	"gim/framework/codec"
	"gim/logic/bucket"
	"github.com/ebar-go/ego/errors"
	uuid "github.com/satori/go.uuid"
	"log"
)

type Controller struct {
	bucket *bucket.Bucket
}

func (c *Controller) Install(router *framework.Router) {
	router.WithCodec(codec.Default()).OnNotFound(func(ctx *framework.Context) {
		log.Println("operation not found")
	}).OnError(func(ctx *framework.Context, err error) {
		log.Println("operation error: ", ctx.Operate(), err)
	})

	router.Route(1, framework.StandardHandler[LoginRequest, LoginResponse](c.Login))
	router.Route(2, framework.StandardHandler[SubscribeChannelRequest, SubscribeChannelResponse](c.SubscribeChannel))
	router.Route(3, framework.StandardHandler[SendMessageRequest, SendMessageResponse](c.SendMessage))

}

func (c *Controller) Login(ctx *framework.Context, req *LoginRequest) (*LoginResponse, error) {
	id := uuid.NewV4().String()
	c.bucket.AddSession(bucket.NewSession(id, ctx.Conn()))
	ctx.Conn().Property().Set("uid", id)
	ctx.Conn().Property().Set("name", req.Name)
	return &LoginResponse{ID: id}, nil
}

func (c *Controller) SubscribeChannel(ctx *framework.Context, req *SubscribeChannelRequest) (*SubscribeChannelResponse, error) {
	channel := c.bucket.GetOrCreate(req.ID)
	session := c.bucket.GetSession(GetUIDFromContext(ctx))
	c.bucket.SubscribeChannel(channel, session)
	log.Println("SubscribeChannel:", channel.ID, session.ID)

	return &SubscribeChannelResponse{}, nil
}

func (c *Controller) SendMessage(ctx *framework.Context, req *SendMessageRequest) (*SendMessageResponse, error) {
	msgId := uuid.NewV4().String()
	channel := c.bucket.GetChannel(req.ChannelID)
	if channel == nil {
		return nil, errors.NotFound("channel not found")
	}

	bytes, err := codec.Default().Pack(&codec.Packet{
		Operate:     5,
		ContentType: codec.ContentTypeJSON,
		Seq:         0,
		Body:        nil,
	}, Message{
		ID: msgId, Content: req.Content,
		Sender: MessageUser{
			ID:   GetUIDFromContext(ctx),
			Name: GetNameFromContext(ctx),
		}})

	if err != nil {
		return nil, err
	}
	c.bucket.BroadcastChannel(channel, bytes)

	return &SendMessageResponse{MsgID: msgId}, nil
}

func NewController() *Controller {
	return &Controller{bucket: bucket.NewBucket()}
}

func GetUIDFromContext(ctx *framework.Context) string {
	return ctx.Conn().Property().GetString("uid")
}
func GetNameFromContext(ctx *framework.Context) string {
	return ctx.Conn().Property().GetString("name")
}
