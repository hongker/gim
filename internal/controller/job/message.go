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
}

func NewMessageJob() *MessageJob {
	return &MessageJob{
		queues:   make(map[string]*queue.GenericQueue[*types.Message]),
		cometApp: application.GetCometApplication(),
		provider: api.NewSharedProtoProvider(),
	}
}

func (job *MessageJob) Prepare() {
	job.initialize()
}

func (job *MessageJob) initialize() {
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
	q := queue.NewGenericQueue[*types.Message](10, true)
	go func() {
		defer runtime.HandleCrash()
		q.Poll(time.Second, func(items []*types.Message) {
			packet := &types.MessagePacket{Session: session, Items: items}
			//proto := c.provider.Acquire()
			//if err := proto.Marshal(packet); err != nil {
			//	return
			//}

			// send private message
			if session.IsPrivate() {
				uid := session.GetPrivateUid()

				runtime.HandlerError(job.cometApp.PushUserMessage(uid, packet.Encode()), func(err error) {
					component.Provider().Logger().Errorf("push user message: %v, %v", uid, err)
				})
			} else if session.IsChatroom() {
				roomId := session.GetChatroomId()
				runtime.HandlerError(job.cometApp.PushChatroomMessage(roomId, packet.Encode()), func(err error) {
					component.Provider().Logger().Errorf("push room message: %v, %v", roomId, err)
				})
			}

		})
	}()
	job.queues[session.Id] = q
	return q
}
