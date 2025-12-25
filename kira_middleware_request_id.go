package kira

type RequestID struct{}

func NewRequestID() *RequestID {
	return &RequestID{}
}

func (rq *RequestID) Middleware(ctx *Context, next HandlerFunc) error {
	headerName := ctx.Config().GetString("server.request_id", "X-Request-Id")

	// Set header.
	ctx.Response().Header().Set(headerName, ctx.RequestID())

	// Move to the next handler.
	return next(ctx)
}
