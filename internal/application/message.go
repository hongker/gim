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

// MessaegApplication represents message application
type MessageApplication interface {
	// Send sends a message to the receiver or chatroom.
	Send(ctx context.Context, uid string, req *dto.MessageSendRequest) (resp *dto.MessageSendResponse, err error)

	// Query return message response.
	Query(ctx context.Context, req *dto.MessageQueryRequest) (*dto.MessageQueryResponse, error)

	// ListSession returns user sessions.
	ListSession(ctx context.Context, uid string, req *dto.SessionQueryRequest) (*dto.SessionQueryResponse, error)
}

// messageApplication implements the MessageApplication interface
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

	switch req.Type {
	case types.SessionPrivate:
		err = app.sendPrivate(ctx, sender, req)
	case types.SessionChatroom:
		err = app.sendChatroom(ctx, sender, req)
	default:
		err = errors.InvalidParam("unknown session type")

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

// -------------------------private methods------------------------
func (app messageApplication) sendPrivate(ctx context.Context, sender *entity.User, req *dto.MessageSendRequest) (err error) {
	// find receiver info
	receiver, err := app.userRepo.Find(ctx, req.TargetId)
	if err != nil {
		return errors.WithMessage(err, "find receiver")
	}

	// save source message
	msg := &types.Message{
		Id:        uuid.NewV4().String(),
		SenderId:  sender.Id,
		Content:   req.Content,
		Category:  types.MessageCategory(req.Category),
		Status:    0,
		CreatedAt: time.Now().UnixMilli(),
	}

	// save session message of sender and receiver.
	err = runtime.Call(func() error {
		senderSession := types.NewPrivateSession(sender.Id, receiver.Id, receiver.Name)
		return senderSession.SaveAndDelivery(ctx, app.sessionRepo, sender, msg)
	}, func() error {
		receiverSession := types.NewPrivateSession(receiver.Id, sender.Id, sender.Name)
		return receiverSession.SaveAndDelivery(ctx, app.sessionRepo, receiver, msg)
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
	msg := &types.Message{
		Id:        uuid.NewV4().String(),
		SenderId:  sender.Id,
		Content:   req.Content,
		Category:  types.MessageCategory(req.Category),
		Status:    0,
		CreatedAt: time.Now().UnixMilli(),
	}

	chatroomSession := types.NewChatroomSession(chatroom.Id, chatroom.Name)
	return chatroomSession.SaveAndDelivery(ctx, app.sessionRepo, sender, msg)

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
	runtime.HandleError(err, func(err error) {
		component.Provider().Logger().Errorf("deliverySessionMessage: %v", err)
	})
}

func NewMessageApplication() MessageApplication {
	return &messageApplication{
		userRepo:     repository.NewUserRepository(),
		sessionRepo:  repository.NewSessionRepository(),
		chatroomRepo: repository.NewChatroomRepository(),
	}
}
