package job

import (
	"gim/api"
	"gim/internal/application"
	"gim/internal/domain/dto"
	"gim/internal/domain/types"
	"gim/pkg/queue"
	"github.com/ebar-go/ego/component"
	"github.com/ebar-go/ego/utils/runtime"
	"sync"
	"time"
)

type MessageJob struct {
	cometApp application.CometApplication
	mu       sync.Mutex
	queues   map[string]*queue.GenericQueue[*types.Message]
	provider *api.SharedProtoProvider

	pollInterval time.Duration
	pollCount    int
}

func NewMessageJob(queuePollInterval time.Duration, queuePollCount int) *MessageJob {
	return &MessageJob{
		queues:       make(map[string]*queue.GenericQueue[*types.Message]),
		cometApp:     application.GetCometApplication(),
		provider:     api.NewSharedProtoProvider(),
		pollInterval: queuePollInterval,
		pollCount:    queuePollCount,
	}
}

func (job *MessageJob) Prepare() {
	job.initialize()
}

func (job *MessageJob) initialize() {
	// listen event.
	component.ListenEvent[*types.SessionMessage](dto.EventDeliveryMessage, func(item *types.SessionMessage) {
		job.getOrInitQueue(item.Session).Offer(item.Message)
	})
}

func (job *MessageJob) getOrInitQueue(session *types.Session) *queue.GenericQueue[*types.Message] {
	job.mu.Lock()
	defer job.mu.Unlock()
	if q, ok := job.queues[session.Id]; ok {
		return q
	}
	q := queue.NewGenericQueue[*types.Message](job.pollCount, true)
	go func() {
		defer runtime.HandleCrash()
		q.Poll(job.pollInterval, func(items []*types.Message) {
			packet := &types.MessagePacket{Session: session, Items: items}

			// send private message
			if session.IsPrivate() {
				uid := session.GetPrivateUid()

				runtime.HandleError(job.cometApp.PushUserMessage(uid, packet.Encode()), func(err error) {
					component.Provider().Logger().Errorf("push user message: %v, %v", uid, err)
				})
			} else if session.IsChatroom() {
				roomId := session.GetChatroomId()
				runtime.HandleError(job.cometApp.PushChatroomMessage(roomId, packet.Encode()), func(err error) {
					component.Provider().Logger().Errorf("push room message: %v, %v", roomId, err)
				})
			}

		})
	}()
	job.queues[session.Id] = q
	return q
}
