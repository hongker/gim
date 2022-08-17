package helper

import (
	"gim/api"
	"gim/internal/domain/dto"
	"gim/pkg/network"
)

func SetContextPacket(ctx *network.Context, packet *api.Packet) {
	ctx.WithValue("packet", packet)
}

func GetContextPacket(ctx *network.Context) *api.Packet {
	return ctx.Value("packet").(*api.Packet)
}

func BuildResponsePacket(ctx *network.Context, data interface{}) *api.Packet {
	packet := GetContextPacket(ctx)
	packet.Op++
	packet.Marshal(data)
	return packet
}

func Bind(ctx *network.Context, container interface{}) error {
	return GetContextPacket(ctx).Bind(container)
}

func SetContextUser(ctx *network.Context, user *dto.User) {
	ctx.WithValue("user", user)
}
func GetContextUser(ctx *network.Context) *dto.User {
	return ctx.Value("user").(*dto.User)
}
