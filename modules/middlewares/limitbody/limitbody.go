package limitbody

import (
	"errors"
	"net/http"

	"github.com/lafriakh/kira"
)

var (
	ErrRequestTooLarge = kira.E("Request too large", http.StatusRequestEntityTooLarge)
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
	limit := ctx.Config().GetInt64("server.body_limit", 32) * MB

	// Early rejection if Content-Length is known and exceeds limit
	if ctx.Request().ContentLength > 0 && ctx.Request().ContentLength > limit {
		return ErrRequestTooLarge
	}

	ctx.Request().Body = http.MaxBytesReader(
		ctx.Response(),
		ctx.Request().Body, limit,
	)

	err := next(ctx)

	// Convert http.MaxBytesError to ErrRequestTooLarge
	var maxBytesErr *http.MaxBytesError
	if errors.As(err, &maxBytesErr) {
		return ErrRequestTooLarge
	}

	return err
}
