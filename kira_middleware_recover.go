package kira

import "net/http"

type Recover struct{}

func NewRecover() *Recover {
	return &Recover{}
}

func (rec *Recover) Middleware(ctx *Context, next HandlerFunc) error {
	err := next(ctx)

	if err != nil {
		if ctx.WantsJSON() {
			ctx.Response().Header().Set("Content-Type", "application/json")
			var defaultErr = &Error{
				Message: http.StatusText(http.StatusInternalServerError),
				Status:  http.StatusInternalServerError,
			}

			if e, ok := AsError(err); ok {
				ctx.SetStatusCode(int(e.Status))
				ctx.Status(int(e.Status))
				_ = ctx.JSON(e)
			} else {
				ctx.SetStatusCode(int(defaultErr.Status))
				ctx.Status(int(defaultErr.Status))
				_ = ctx.JSON(defaultErr)
			}
			ctx.Log().Error(err)
		}

		return err
	}

	return nil
}
