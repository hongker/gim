package gateway

import "encoding/json"

type OperateType int

const (
	LoginOperate      OperateType = 101
	LoginOperateReply OperateType = 102

	LogoutOperate      OperateType = 103
	LogoutOperateReply OperateType = 104

	HeartbeatOperate      OperateType = 105
	HeartbeatOperateReply OperateType = 106

	MessageSendOperate      OperateType = 201
	MessageSendOperateReply OperateType = 202

	ChatroomJoinOperate      OperateType = 301
	ChatroomJoinOperateReply OperateType = 302
)

type Proto struct {
	Operate OperateType `json:"op"`
	Body    string      `json:"body"`
	Seq     int         `json:"seq"`
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
