package framework

import "context"

type Handler func(ctx *Context, serializer Serializer) (any, error)

func StandardHandler[Request, Response any](action func(ctx context.Context, request *Request) (*Response, error)) Handler {
	return func(ctx *Context, serializer Serializer) (any, error) {
		request := new(Request)
		if err := serializer.Unmarshal(ctx.body, request); err != nil {
			return nil, err
		}
		return action(ctx, request)

	}
}
