package kira

// Middleware interface
type Middleware interface {
	// Name() string
	Middleware(*Context, HandlerFunc) error
}

// Middleware - add the middleware
func (app *App) Middleware(middlewares ...Middleware) {
	app.Middlewares = append(app.Middlewares, middlewares...)
}

func defaultMiddlewares() (mds []Middleware) {
	mds = append(mds, NewRecover())
	mds = append(mds, NewRequestID())
	mds = append(mds, NewLogger())

	return mds
}
