package job

import (
	"gim/internal/application"
	"gim/internal/domain/dto"
	"gim/internal/domain/types"
	"gim/pkg/queue"
	"github.com/ebar-go/ego/component"
	"github.com/ebar-go/ego/utils/runtime"
	"sync"
	"time"
)

// Controller represents cronjob controller.
type Controller struct {
	name string
	once sync.Once

	config   *Config
	cometApp application.CometApplication
	mu       sync.Mutex
	queues   map[string]*queue.GenericQueue[*types.Message]
}

func (c *Controller) Run(stopCh <-chan struct{}) {
	c.once.Do(c.initialize)
	c.run()

	runtime.WaitClose(stopCh, c.shutdown)
}

func (c *Controller) WithName(name string) *Controller {
	c.name = name
	return c
}

func (c *Controller) initialize() {
	c.queues = make(map[string]*queue.GenericQueue[*types.Message])
	component.ListenEvent[*types.SessionMessage](dto.EventDeliveryMessage, func(item *types.SessionMessage) {
		q := c.prepareQueue(item.Session)
		q.Offer(item.Message)
	})

}

func (c *Controller) prepareQueue(session *types.Session) *queue.GenericQueue[*types.Message] {
	c.mu.Lock()
	defer c.mu.Unlock()
	if q, ok := c.queues[session.Id]; ok {
		return q
	}
	q := queue.NewGenericQueue[*types.Message](10, true)
	go func() {
		defer runtime.HandleCrash()
		q.Poll(time.Second, func(items []*types.Message) {
			packet := &types.MessagePacket{Session: session, Items: items}
			// send private message
			if session.IsPrivate() {
				uid := session.GetPrivateUid()

				runtime.HandlerError(c.cometApp.PushUserMessage(uid, packet.Encode()), func(err error) {
					component.Provider().Logger().Errorf("push user message: %v, %v", uid, err)
				})
			} else if session.IsChatroom() {
				roomId := session.GetChatroomId()
				runtime.HandlerError(c.cometApp.PushChatroomMessage(roomId, packet.Encode()), func(err error) {
					component.Provider().Logger().Errorf("push room message: %v, %v", roomId, err)
				})
			}

		})
	}()
	c.queues[session.Id] = q
	return q
}

func (c *Controller) run() {
	component.Provider().Logger().Infof("controller running: [%s]", c.name)

}
func (c *Controller) shutdown() {
	component.Provider().Logger().Infof("controller shutdown: [%s]", c.name)
}
