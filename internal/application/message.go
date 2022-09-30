package application

import (
	"context"
	"gim/internal/domain/dto"
	"gim/internal/domain/entity"
	"gim/internal/domain/repository"
	"gim/internal/domain/types"
	"github.com/ebar-go/ego/component"
	"github.com/ebar-go/ego/errors"
	"github.com/ebar-go/ego/utils/runtime"
	uuid "github.com/satori/go.uuid"
	"time"
)

type MessageApplication interface {
	Send(ctx context.Context, uid string, req *dto.MessageSendRequest) (resp *dto.MessageSendResponse, err error)
	Query(ctx context.Context, req *dto.MessageQueryRequest) (*dto.MessageQueryResponse, error)
	ListSession(ctx context.Context, uid string, req *dto.SessionQueryRequest) (*dto.SessionQueryResponse, error)
}

type messageApplication struct {
	userRepo     repository.UserRepository
	sessionRepo  repository.SessionRepository
	chatroomRepo repository.ChatroomRepository
}

func (app messageApplication) Send(ctx context.Context, uid string, req *dto.MessageSendRequest) (resp *dto.MessageSendResponse, err error) {
	sender, err := app.userRepo.Find(ctx, uid)
	if err != nil {
		return nil, errors.WithMessage(err, "find sender")
	}

	if req.Type == types.SessionPrivate {
		err = app.sendPrivate(ctx, sender, req)
	} else {
		err = app.sendChatroom(ctx, sender, req)
	}

	return
}

func (app messageApplication) Query(ctx context.Context, req *dto.MessageQueryRequest) (*dto.MessageQueryResponse, error) {
	session := &entity.Session{Id: req.SessionId}
	messages, err := app.sessionRepo.QueryMessage(ctx, session)
	if err != nil {
		return nil, err
	}

	res := &dto.MessageQueryResponse{Items: make([]dto.MessageItem, 0, len(messages))}
	for _, message := range messages {
		sender, lastErr := app.userRepo.Find(ctx, message.SenderId)
		if lastErr != nil {
			component.Provider().Logger().Errorf("user not found: %v", lastErr)
			continue
		}
		res.Items = append(res.Items, dto.MessageItem{
			Id:      message.Id,
			Content: message.Content,
			Sender: dto.MessageUser{
				Id:   sender.Id,
				Name: sender.Name,
			},
		})
	}

	return res, nil
}

func (app messageApplication) ListSession(ctx context.Context, uid string, req *dto.SessionQueryRequest) (*dto.SessionQueryResponse, error) {
	items, err := app.sessionRepo.List(ctx, uid)
	if err != nil {
		return nil, err
	}

	res := &dto.SessionQueryResponse{Items: make([]dto.Session, 0, len(items))}
	for _, item := range items {
		msg, lastErr := app.sessionRepo.FindMessage(ctx, item.Last)
		if lastErr != nil {
			continue
		}
		sender, lastErr := app.userRepo.Find(ctx, msg.SenderId)
		if lastErr != nil {
			continue
		}
		res.Items = append(res.Items, dto.Session{
			Id:    item.Id,
			Title: item.Title,
			Type:  types.CategoryFromSessionId(item.Id),
			Last: &dto.MessageItem{
				Id:      msg.Id,
				Content: msg.Content,
				Sender:  dto.MessageUser{Id: sender.Id, Name: sender.Name},
			}})
	}
	return res, nil
}
func (app messageApplication) sendPrivate(ctx context.Context, sender *entity.User, req *dto.MessageSendRequest) (err error) {
	// find receiver info
	receiver, err := app.userRepo.Find(ctx, req.TargetId)
	if err != nil {
		return errors.WithMessage(err, "find receiver")
	}

	// save source message
	msg := &entity.Message{
		Id:        uuid.NewV4().String(),
		SenderId:  sender.Id,
		Content:   req.Content,
		Category:  req.Category,
		Status:    0,
		CreatedAt: time.Now().UnixMilli(),
	}
	msg.SenderId = sender.Id

	// save session message of sender and receiver.
	err = runtime.Call(func() error {
		senderSession := types.NewPrivateSession(sender.Id, receiver.Id, receiver.Name)

		//go app.deliverySessionMessage(senderSession, msg)
		app.delivery(senderSession, msg)
		return app.sessionRepo.SaveMessage(ctx, sender.Id, senderSession.Entity(), msg)
	}, func() error {
		receiverSession := types.NewPrivateSession(receiver.Id, sender.Id, sender.Name)
		//go app.deliverySessionMessage(receiverSession, msg)
		app.delivery(receiverSession, msg)
		return app.sessionRepo.SaveMessage(ctx, receiver.Id, receiverSession.Entity(), msg)
	})

	return
}

func (app messageApplication) sendChatroom(ctx context.Context, sender *entity.User, req *dto.MessageSendRequest) (err error) {
	chatroom, err := app.chatroomRepo.Find(ctx, req.TargetId)
	if err != nil {
		return
	}

	// user is not allowed to send messages before join chatroom.
	if !app.chatroomRepo.HasMember(ctx, chatroom, sender) {
		return errors.Forbidden("cannot send message before join the chatroom")
	}

	// save source message
	msg := &entity.Message{
		Id:        uuid.NewV4().String(),
		SenderId:  sender.Id,
		Content:   req.Content,
		Category:  req.Category,
		Status:    0,
		CreatedAt: time.Now().UnixMilli(),
	}
	msg.SenderId = sender.Id

	chatroomSession := types.NewChatroomSession(chatroom.Id, chatroom.Name)
	go app.deliverySessionMessage(chatroomSession, msg)
	return app.sessionRepo.SaveMessage(ctx, sender.Id, chatroomSession.Entity(), msg)

}

func (app messageApplication) pushUid(uid string, msg *types.Message) error {
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
func (app messageApplication) deliverySessionMessage(session *types.Session, msg *entity.Message) {
	message := &types.Message{Id: msg.Id, SenderId: msg.SenderId, Category: types.MessageCategory(msg.Category), Content: msg.Content, CreatedAt: msg.CreatedAt}
	var err error
	if session.IsPrivate() {
		uid := session.GetPrivateUid()
		err = app.pushUid(uid, message)

	} else if session.IsChatroom() {
		chatroom, lastErr := app.chatroomRepo.Find(context.TODO(), session.GetChatroomId())
		if lastErr != nil {
			err = errors.WithMessage(lastErr, "find chatroom")
		}
		members, lastErr := app.chatroomRepo.GetMember(context.Background(), chatroom)
		for _, member := range members {
			_ = app.pushUid(member, message)
		}
		err = lastErr
	}
	runtime.HandlerError(err, func(err error) {
		component.Provider().Logger().Errorf("deliverySessionMessage: %v", err)
	})
}

func (app messageApplication) delivery(session *types.Session, msg *entity.Message) {
	message := &types.Message{Id: msg.Id, SenderId: msg.SenderId, Category: types.MessageCategory(msg.Category), Content: msg.Content, CreatedAt: msg.CreatedAt}
	sessionMessage := &types.SessionMessage{Session: session, Message: message}
	component.Provider().EventDispatcher().Trigger(dto.EventDeliveryMessage, sessionMessage)
}
func NewMessageApplication() MessageApplication {
	return &messageApplication{
		userRepo:     repository.NewUserRepository(),
		sessionRepo:  repository.NewSessionRepository(),
		chatroomRepo: repository.NewChatroomRepository(),
	}
}
