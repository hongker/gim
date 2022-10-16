package task

import (
	"gim/api"
	"gim/framework/codec"
	"gim/internal/application"
	"gim/internal/domain/dto"
	"gim/internal/domain/types"
	"github.com/ebar-go/ego/component"
	"github.com/ebar-go/ego/utils/runtime"
	"sync"
	"time"
)

type MessageTask struct {
	cometApp application.CometApplication
	mu       sync.Mutex
	queues   map[string]*Queue[*types.Message]

	pollInterval time.Duration
	pollCount    int
	codec        codec.Codec
}

func NewMessageTask(queuePollInterval time.Duration, queuePollCount int) *MessageTask {
	return &MessageTask{
		queues:       make(map[string]*Queue[*types.Message]),
		cometApp:     application.GetCometApplication(),
		pollInterval: queuePollInterval,
		pollCount:    queuePollCount,
	}
}

func (task *MessageTask) Start() {
	// listen delivery message event
	component.ListenEvent[*types.SessionMessage](dto.EventDeliveryMessage, func(item *types.SessionMessage) {
		task.getOrCreateQueue(item.Session).Offer(item.Message)
	})
}

// getOrCreateQueue returns a queue for the given session
func (task *MessageTask) getOrCreateQueue(session *types.Session) *Queue[*types.Message] {
	task.mu.Lock()
	defer task.mu.Unlock()
	if queue, ok := task.queues[session.Id]; ok {
		return queue
	}

	queue := task.createSessionQueue(session)
	task.queues[session.Id] = queue
	return queue
}

// newSessionQueue return a new queue for the given session
func (task *MessageTask) createSessionQueue(session *types.Session) *Queue[*types.Message] {
	queue := NewQueue[*types.Message](task.pollCount, true)
	go func() {
		defer runtime.HandleCrash()

		// poll with handler
		queue.Poll(task.pollInterval, func(items []*types.Message) {
			task.pushMessages(session, items)
		})
	}()
	return queue
}

// pushMessages push session messages to client
func (task *MessageTask) pushMessages(session *types.Session, messages []*types.Message) {
	packet := &codec.Packet{
		Operate:     api.MessagePushOperate,
		ContentType: codec.ContentTypeJSON,
		Seq:         0,
	}

	bytes, err := task.codec.Pack(packet, &types.SessionMessageItems{Session: session, Items: messages})
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
