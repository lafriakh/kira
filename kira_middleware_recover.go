package kira

import "net/http"

type Recover struct{}

func NewRecover() *Recover {
	return &Recover{}
}

func (rec *Recover) Middleware(ctx *Context, next HandlerFunc) error {
	err := next(ctx)

	if err != nil {
		handleError(ctx, err)
	}

	return nil
}

func handleError(ctx *Context, err error) {
	// Log the error
	defer ctx.Log().Error(err)

	// Error
	var e error

	// Status code
	code := http.StatusInternalServerError
	if err, ok := AsError(err); ok {
		e = err
		code = int(err.Status)
	} else {
		e = &Error{
			Message: http.StatusText(code),
			Status:  code,
		}
	}
	ctx.SetStatusCode(code)
	ctx.Status(code)

	if ctx.WantsJSON() {
		ctx.JSON(e)
	} else {
		ctx.WriteString(http.StatusText(ctx.StatusCode()))
	}
}
