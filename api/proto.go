package api

import (
	"encoding/json"
	"sync"
)

type OperateType int

const (
	LoginOperate      = 101
	LoginOperateReply = 102

	LogoutOperate                  = 103
	LogoutOperateReply OperateType = 104

	HeartbeatOperate      = 105
	HeartbeatOperateReply = 106

	MessageSendOperate                  = 201
	MessageSendOperateReply OperateType = 202

	MessageQueryOperate                  = 203
	MessageQueryOperateReply OperateType = 204

	MessagePushOperate                  = 205
	MessagePushOperateReply OperateType = 206

	SessionListOperate                  = 205
	SessionListOperateReply OperateType = 206

	ChatroomJoinOperate      = 301
	ChatroomJoinOperateReply = 302
)

type Proto struct {
	Operate OperateType `json:"op"`
	Body    string      `json:"body"`
	Seq     int         `json:"seq"`
}

func (p *Proto) reset() {
	p.Operate = 0
	p.Body = ""
	p.Seq = 0
}

func (p *Proto) OperateType() OperateType {
	return p.Operate
}

func (p *Proto) Bind(container any) error {
	if p.Body == "" {
		return nil
	}
	return json.Unmarshal([]byte(p.Body), container)
}

func (p *Proto) BindFunc(container any) func() error {
	return func() error {
		return p.Bind(container)
	}
}

func (p *Proto) Marshal(container interface{}) (err error) {
	b, err := json.Marshal(container)
	if err != nil {
		return err
	}
	p.Body = string(b)
	p.Seq++
	p.Operate++
	return
}

type ProtoProvider interface {
	Acquire() *Proto
	Release(p *Proto)
}

// TOODO 使用泛型
// SharedProvider[T any]
type SharedProtoProvider struct {
	pool sync.Pool
}

func NewSharedProtoProvider() *SharedProtoProvider {
	return &SharedProtoProvider{
		pool: sync.Pool{
			New: func() interface{} {
				return new(Proto)
			}},
	}
}

func (provider *SharedProtoProvider) Acquire() *Proto {
	p := provider.pool.Get().(*Proto)
	p.reset()
	return p
}

func (provider *SharedProtoProvider) Release(p *Proto) {
	provider.pool.Put(p)
}
