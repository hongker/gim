package socket

import "encoding/json"

type OperateType int

const (
	LoginOperate      OperateType = 101
	LoginOperateReply OperateType = 102

	HeartbeatOperate      OperateType = 103
	HeartbeatOperateReply OperateType = 104
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
	if validator, ok := container.(Validatable); ok {
		if err := validator.Validate(); err != nil {
			return err
		}
	}
	return json.Unmarshal([]byte(p.Body), container)
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

type Validatable interface {
	Validate() error
}
