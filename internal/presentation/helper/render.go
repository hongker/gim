package helper

import (
	"gim/api"
	"gim/pkg/network"
)

type Render struct {
}

func (render *Render) write(ctx *network.Context, b []byte) {
	ctx.Connection().Push(b)
}

func NewRender() *Render {
	return &Render{}
}

var currentReader = NewRender()

func SetRender(render *Render) {
	currentReader = render
}

func Success(ctx *network.Context, data interface{}) {
	response := api.SuccessResponse(data)
	packet := BuildResponsePacket(ctx, response)
	currentReader.write(ctx, packet.Encode())
}

func Failure(ctx *network.Context, err error) {
	response := api.FailureResponse(err)
	packet := BuildResponsePacket(ctx, response)
	currentReader.write(ctx, packet.Encode())
}
