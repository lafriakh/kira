package limitbody

import (
	"net/http"

	"github.com/lafriakh/kira"
)

var (
	ErrRequestTooLarge = kira.E("Request too large", kira.StatusCode(http.StatusRequestEntityTooLarge))
)

// MB - one MB.
const MB = 1 << 20

// Limitbody - Middleware.
type Limitbody struct{}

// New - return Limitbody instance
func New() *Limitbody {
	return &Limitbody{}
}

// Middleware handler.
func (l *Limitbody) Middleware(ctx *kira.Context, next kira.HandlerFunc) error {
	if ctx.Request().ContentLength > ctx.Config().GetInt64("server.body_limit", 32)*MB {
		return ErrRequestTooLarge
	}

	ctx.Request().Body = http.MaxBytesReader(
		ctx.Response(),
		ctx.Request().Body, ctx.Config().GetInt64("server.body_limit", 32)*MB,
	)

	return next(ctx)
}
