package filter

import (
	"gim/api"
	"gim/internal/domain/types"
	"gim/internal/presentation/helper"
	"gim/pkg/errors"
	"gim/pkg/network"
	"log"
)

func Unpack(ctx *network.Context) {
	packet := api.NewPacket()
	if err := packet.Decode(ctx.Request().Body()); err != nil {
		helper.Failure(ctx, errors.InvalidParameter(err.Error()))
		ctx.Abort()
		return
	}

	helper.SetContextPacket(ctx, packet)
	ctx.Next()
}

func Auth(ctx *network.Context) {
	packet := helper.GetContextPacket(ctx)
	if packet.Op == api.OperateAuth {
		ctx.Next()
		return
	}

	user := types.GetCollection().GetUser(ctx.Connection())
	if user == nil {
		helper.Failure(ctx, errors.New(errors.CodeForbidden, "auth is required"))
		ctx.Abort()
		return
	}
	helper.SetContextUser(ctx, user)
	ctx.Next()
}

func Recover(ctx *network.Context) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("recover: err=%v\n", err)
		}
	}()

	ctx.Next()
}
