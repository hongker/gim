package helper

import (
	"gim/api"
	"gim/pkg/network"
)

func Success(ctx *network.Context, data interface{})  {
	response := api.SuccessResponse(data)
	packet := GetContextPacket(ctx)
	_ = packet.Marshal(response)
	packet.Op += 1
	ctx.Connection().Push(packet.Encode())
}

func Failure(ctx *network.Context, err error)  {
	response := api.FailureResponse(err)
	packet := GetContextPacket(ctx)
	_ = packet.Marshal(response)
	packet.Op += 1
	ctx.Connection().Push(packet.Encode())
}
