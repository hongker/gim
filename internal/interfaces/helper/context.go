package helper

import (
	"gim/api"
	"gim/pkg/network"
)

func SetContextPacket(ctx *network.Context, packet *api.Packet)  {
	ctx.WithValue("packet", packet)
}

func GetContextPacket(ctx *network.Context) *api.Packet  {
	return ctx.Value("packet").(*api.Packet)
}

func Bind(ctx *network.Context, container interface{}) error  {
	return GetContextPacket(ctx).Bind(container)
}
