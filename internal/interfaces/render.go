package interfaces

import "gim/pkg/network"

func Success(ctx *network.Context, p []byte)  {
	ctx.Connection().Push(p)
}

func Failure(ctx *network.Context, err error)  {
	
}
