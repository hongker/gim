package application

import (
	"context"
	"gim/internal/domain/dto"
	"gim/internal/domain/entity"
	repository2 "gim/internal/domain/repository"
	types2 "gim/internal/domain/types"
	"github.com/ebar-go/ego/component"
	"github.com/ebar-go/ego/errors"
	"github.com/ebar-go/ego/utils/runtime"
)

type MessageApplication interface {
	Send(ctx context.Context, uid string, req *dto.MessageSendRequest) (resp *dto.MessageSendResponse, err error)
}

type messageApplication struct {
	userRepo     repository2.UserRepository
	sessionRepo  repository2.SessionRepository
	msgRepo      repository2.MessageRepository
	chatroomRepo repository2.ChatroomRepository
}

func (app messageApplication) Send(ctx context.Context, uid string, req *dto.MessageSendRequest) (resp *dto.MessageSendResponse, err error) {
	sender, err := app.userRepo.Find(ctx, uid)
	if err != nil {
		return nil, errors.WithMessage(err, "find sender")
	}

	if req.Type == string(types2.SessionPrivate) {
		err = app.sendPrivate(ctx, sender, req)
	} else {
		err = app.sendChatroom(ctx, sender, req)
	}

	return
}

func (app messageApplication) sendPrivate(ctx context.Context, sender *entity.User, req *dto.MessageSendRequest) (err error) {
	// find receiver info
	receiver, err := app.userRepo.Find(ctx, req.TargetId)
	if err != nil {
		return errors.WithMessage(err, "find receiver")
	}

	// save source message
	msg := types2.NewTextMessage(req.Content)
	msg.SenderId = sender.Id
	err = app.msgRepo.Save(ctx, msg)
	if err != nil {
		return errors.WithMessage(err, "save message")
	}

	// save session message of sender and receiver.
	err = runtime.Call(func() error {
		senderSession := types2.NewPrivateSession(sender.Id, receiver.Id, receiver.Name)
		go app.deliverySessionMessage(senderSession, msg)
		return app.sessionRepo.SaveMessage(ctx, senderSession, msg)
	}, func() error {
		receiverSession := types2.NewPrivateSession(receiver.Id, sender.Id, sender.Name)
		go app.deliverySessionMessage(receiverSession, msg)
		return app.sessionRepo.SaveMessage(ctx, receiverSession, msg)
	})

	return
}

func (app messageApplication) sendChatroom(ctx context.Context, sender *entity.User, req *dto.MessageSendRequest) (err error) {
	chatroom, err := app.chatroomRepo.Find(ctx, req.TargetId)
	if err != nil {
		return
	}

	// save source message
	msg := types2.NewTextMessage(req.Content)
	msg.SenderId = sender.Id
	err = app.msgRepo.Save(ctx, msg)
	if err != nil {
		return errors.WithMessage(err, "save message")
	}

	chatroomSession := types2.NewChatroomSession(chatroom.Id, chatroom.Name)
	go app.deliverySessionMessage(chatroomSession, msg)
	return app.sessionRepo.SaveMessage(ctx, chatroomSession, msg)

}

func (app messageApplication) pushUid(uid string, msg *types2.Message) error {
	conn, err := GetCometApplication().GetUserConnection(uid)
	if err != nil {
		return err
	}

	bytes, err := msg.Encode()
	if err != nil {
		return err
	}
	return conn.Push(bytes)

}
func (app messageApplication) deliverySessionMessage(session *types2.Session, msg *types2.Message) {
	var err error
	if session.IsPrivate() {
		uid := session.GetPrivateUid()
		err = app.pushUid(uid, msg)

	} else if session.IsChatroom() {
		members, lastErr := app.chatroomRepo.GetMember(context.Background(), session.GetChatroomId())
		for _, member := range members {
			_ = app.pushUid(member, msg)
		}
		err = lastErr
	}
	runtime.HandlerError(err, func(err error) {
		component.Provider().Logger().Errorf("deliverySessionMessage: %v", err)
	})
}

func NewMessageApplication() MessageApplication {
	return &messageApplication{
		userRepo:     repository2.NewUserRepository(),
		msgRepo:      repository2.NewMessageRepository(),
		sessionRepo:  repository2.NewSessionRepository(),
		chatroomRepo: repository2.NewChatroomRepository(),
	}
}
