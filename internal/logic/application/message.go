package application

import (
	"context"
	"fmt"
	"gim/api"
	"gim/api/client"
	"gim/api/protocol"
	"gim/internal/logic/domain/entity"
	"gim/internal/logic/domain/repository"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/protobuf/proto"
	"strconv"
	"time"
)

type MessageApp struct {
	userRepo    repository.UserRepo
	messageRepo repository.MessageRepo
}

func NewMessage(userRepo repository.UserRepo, messageRepo repository.MessageRepo) *MessageApp {
	return &MessageApp{userRepo: userRepo, messageRepo: messageRepo}
}

func (app *MessageApp) Query(ctx context.Context, uid string, proto *protocol.Proto) (res proto.Message, err error) {
	request := new(client.MessageQueryRequest)
	if err = proto.Bind(request); err != nil {
		return
	}

	messages := app.messageRepo.Query(request.SessionId, request.LastMsgId, int(request.Count))

	items := make([]*client.MessageItem, len(messages))
	for i, message := range messages {
		items[i] = &client.MessageItem{
			Id:          message.Id,
			Sender:      nil,
			Type:        message.Type,
			Content:     message.Content,
			CreatedAt:   message.Time,
			ClientMsgId: message.ClientMsgId,
			Sequence:    message.Sequence,
		}
	}
	res = &client.MessageQueryResponse{List: items}

	return
}
func (app *MessageApp) Send(ctx context.Context, uid string, proto *protocol.Proto) (res proto.Message, err error) {
	sender, err := app.userRepo.Find(ctx, uid)
	if err != nil {
		return
	}

	request := new(client.MessageSendRequest)
	if err = proto.Bind(request); err != nil {
		return
	}

	message := app.generate(request)
	message.FromUser = sender

	if request.SessionType == api.SessionTypePrivate {
		var receiver *entity.User
		receiver, err = app.userRepo.Find(ctx, request.TargetId)
		if err != nil {
			return
		}
		err = app.sendPrivate(message, receiver)
	} else {
		group := &entity.Group{
			Id: request.TargetId,
		}
		err = app.sendGroup(message, group)
	}

	return
}

func (app *MessageApp) generate(request *client.MessageSendRequest) (msg *entity.Message) {
	msg = &entity.Message{
		Type:        request.Type,
		Content:     request.Content,
		Time:        time.Now().UnixNano(),
		ClientMsgId: request.ClientMsgId,
		SessionId:   fmt.Sprintf("%d:%s", request.SessionType, request.TargetId),
	}
	msg.Id = uuid.NewV4().String() + strconv.FormatInt(msg.Time, 10)
	return msg
}

func (app *MessageApp) sendPrivate(message *entity.Message, receiver *entity.User) error {
	app.messageRepo.Save(message)
	return nil
}

func (app *MessageApp) sendGroup(message *entity.Message, group *entity.Group) error {
	app.messageRepo.Save(message)
	return nil
}
