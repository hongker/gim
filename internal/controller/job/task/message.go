package task

import (
	"gim/api"
	"gim/framework/codec"
	"gim/internal/application"
	"gim/internal/domain/dto"
	"gim/internal/domain/types"
	"gim/pkg/queue"
	"github.com/ebar-go/ego/component"
	"github.com/ebar-go/ego/utils/runtime"
	"sync"
	"time"
)

type MessageTask struct {
	cometApp application.CometApplication
	mu       sync.Mutex
	queues   map[string]*queue.GenericQueue[*types.Message]

	pollInterval time.Duration
	pollCount    int
	codec        codec.Codec
}

func NewMessageTask(queuePollInterval time.Duration, queuePollCount int) *MessageTask {
	return &MessageTask{
		queues:       make(map[string]*queue.GenericQueue[*types.Message]),
		cometApp:     application.GetCometApplication(),
		pollInterval: queuePollInterval,
		pollCount:    queuePollCount,
	}
}

func (task *MessageTask) Start() {
	task.initialize()
}

func (task *MessageTask) initialize() {
	// listen event.
	task.listenEvent()
}

func (task *MessageTask) listenEvent() {
	component.ListenEvent[*types.SessionMessage](dto.EventDeliveryMessage, func(item *types.SessionMessage) {
		task.getOrInitQueue(item.Session).Offer(item.Message)
	})
}

func (task *MessageTask) getOrInitQueue(session *types.Session) *queue.GenericQueue[*types.Message] {
	task.mu.Lock()
	defer task.mu.Unlock()
	if q, ok := task.queues[session.Id]; ok {
		return q
	}

	q := task.initQueue(session)
	task.queues[session.Id] = q
	return q
}

func (task *MessageTask) initQueue(session *types.Session) *queue.GenericQueue[*types.Message] {
	q := queue.NewGenericQueue[*types.Message](task.pollCount, true)
	go func() {
		defer runtime.HandleCrash()
		q.Poll(task.pollInterval, func(items []*types.Message) {
			task.handleSessionMessage(session, items)
		})
	}()
	return q
}

func (task *MessageTask) handleSessionMessage(session *types.Session, messages []*types.Message) {
	packet := &types.MessagePacket{Session: session, Items: messages}

	bytes, err := task.codec.Pack(&codec.Packet{
		Operate:     api.MessagePushOperate,
		ContentType: codec.ContentTypeJSON,
		Seq:         0,
	}, packet)
	if err != nil {
		return
	}
	// send private message
	if session.IsPrivate() {
		uid := session.GetPrivateUid()

		runtime.HandleError(task.cometApp.PushUserMessage(uid, bytes), func(err error) {
			component.Provider().Logger().Errorf("push user message: %v, %v", uid, err)
		})
	} else if session.IsChatroom() {
		roomId := session.GetChatroomId()
		runtime.HandleError(task.cometApp.PushChatroomMessage(roomId, bytes), func(err error) {
			component.Provider().Logger().Errorf("push room message: %v, %v", roomId, err)
		})
	}
}
