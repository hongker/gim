package framework

import "context"

type Handler func(ctx *Context)

func StandardHandler[Request, Response any](action func(ctx context.Context, request *Request) (*Response, error)) Handler {
	return func(ctx *Context) {
		request := new(Request)
		if err := ctx.Bind(request); err != nil {
			ctx.Render().Failure(err)
		}

		response, err := action(ctx, request)
		if err != nil {
			ctx.Render().Failure(err)
		}

		ctx.Render().Success(response)

	}
}
