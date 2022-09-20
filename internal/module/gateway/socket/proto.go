package socket

import "encoding/json"

type OperateType int

const (
	ConnectOperate    OperateType = 100
	DisconnectOperate OperateType = 101
	HeartbeatOperate  OperateType = 102
)

type Proto struct {
	Operate OperateType `json:"operate"`
	Body    []byte      `json:"body"`
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
	return json.Unmarshal(p.Body, container)
}

func (p *Proto) Marshal(container interface{}) (err error) {
	p.Body, err = json.Marshal(container)
	p.Seq++
	return
}

type Validatable interface {
	Validate() error
}
