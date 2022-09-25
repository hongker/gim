package application

import (
	"context"
	"gim/internal/module/gateway/domain/dto"
	"gim/internal/module/gateway/domain/entity"
	"gim/internal/module/gateway/domain/repository"
	"gim/internal/module/gateway/domain/types"
	"github.com/ebar-go/ego/errors"
	"github.com/ebar-go/ego/utils/runtime"
)

type MessageApplication interface {
	Send(ctx context.Context, uid string, req *dto.MessageSendRequest) (resp *dto.MessageSendResponse, err error)
}

type messageApplication struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
	msgRepo     repository.MessageRepository
}

func (app messageApplication) Send(ctx context.Context, uid string, req *dto.MessageSendRequest) (resp *dto.MessageSendResponse, err error) {
	sender, err := app.userRepo.Find(ctx, uid)
	if err != nil {
		return nil, errors.WithMessage(err, "find sender")
	}

	err = app.SendUserSessionMessage(ctx, sender, req)
	return
}

func (app messageApplication) SendUserSessionMessage(ctx context.Context, sender *entity.User, req *dto.MessageSendRequest) (err error) {
	// find receiver info
	receiver, err := app.userRepo.Find(ctx, req.TargetId)
	if err != nil {
		return errors.WithMessage(err, "find receiver")
	}

	// save source message
	msg := types.NewTextMessage(req.Content)
	err = app.msgRepo.Save(ctx, msg)
	if err != nil {
		return errors.WithMessage(err, "save message")
	}

	// save session message of sender and receiver.
	err = runtime.Call(func() error {
		senderSession := types.NewPrivateSession(sender.Id, receiver.Id, receiver.Name)
		go app.deliverySessionMessage(senderSession, msg)
		return app.sessionRepo.SaveMessage(ctx, senderSession, msg)
	}, func() error {
		receiverSession := types.NewPrivateSession(receiver.Id, sender.Id, sender.Name)
		go app.deliverySessionMessage(receiverSession, msg)
		return app.sessionRepo.SaveMessage(ctx, receiverSession, msg)
	})

	return
}

func (app messageApplication) deliverySessionMessage(session *types.Session, msg *types.Message) {
	if session.IsPrivate() {
		uid := session.GetPrivateUid()
		conn, err := GetCometApplication().GetUserConnection(uid)
		if err != nil {
			return
		}
		bytes, err := msg.Encode()
		if err != nil {
			return
		}
		_ = conn.Push(bytes)
	} else if session.IsChatroom() {

	}
}

func NewMessageApplication() MessageApplication {
	return &messageApplication{}
}
