package cors

import (
	"net/http"

	"github.com/lafriakh/kira"
	"github.com/rs/cors"
)

type CORS struct {
	cors *cors.Cors
}

func New(options ...cors.Options) *CORS {
	if len(options) == 0 {
		return AllowAll()
	}
	return &CORS{
		cors: cors.New(options[0]),
	}
}

func AllowAll() *CORS {
	return &CORS{
		cors: cors.AllowAll(),
	}
}

func (c *CORS) Middleware(ctx *kira.Context, next kira.HandlerFunc) error {
	var err error
	c.cors.ServeHTTP(ctx.Response(), ctx.Request(), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx.SetResponse(w)
		ctx.SetRequest(r)

		err = next(ctx)
	}))
	return err
}
