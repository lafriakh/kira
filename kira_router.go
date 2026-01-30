package kira

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

type Method string

const (
	GET     Method = "GET"
	HEAD    Method = "HEAD"
	POST    Method = "POST"
	PUT     Method = "PUT"
	PATCH   Method = "PATCH"
	DELETE  Method = "DELETE"
	OPTIONS Method = "OPTIONS"
)

// Route represent a route.
type Route struct {
	Method      Method
	Path        string
	HandlerFunc HandlerFunc
	Middlewares []Middleware
}

// Middleware add a middleware to the route.
func (r *Route) Middleware(midd ...Middleware) {
	r.Middlewares = append(r.Middlewares, midd...)
}

// Use is an alias of Middleware method.
func (r *Route) Use(midd ...Middleware) {
	r.Middleware(midd...)
}

// RegisterRoutes it's simply register the routes into the router.
func (app *App) RegisterRoutes() *httprouter.Router {
	// build the routes and attach the middlewares to every route.
	for _, route := range app.Routes {
		// Register the route.
		app.Router.Handler(
			// Method
			string(route.Method),
			// Path
			route.Path,
			// Handler
			buildRoute(app, route.HandlerFunc, route.Middlewares),
		)
	}

	// Handle not found requests
	if app.NotFoundHandler == nil {
		app.Router.NotFound = buildRoute(app, defaultNotFound, nil)
	} else {
		app.Router.NotFound = buildRoute(app, app.NotFoundHandler, nil)
	}

	// Handle options requests
	app.Router.GlobalOPTIONS = buildRoute(app, noopHandler, nil)

	// Handle panic
	app.Router.PanicHandler = func(w http.ResponseWriter, r *http.Request, err any) {
		var panicHandler = func(ctx *Context) error {
			ctx.Log().Warn("Panic")
			ctx.Log().Error(err)

			return &Error{
				Status:  http.StatusInternalServerError,
				Message: http.StatusText(http.StatusInternalServerError),
			}
		}
		buildRoute(app, panicHandler, nil)(w, r)
	}

	return app.Router
}

// Change the middleware to support middleware chain.
// This function will take the middleware and the next handler as a parameters.
// Then return a handler that accept the next handler as a parameter.
func buildMiddleware(middleware Middleware, next HandlerFunc) HandlerFunc {
	return func(ctx *Context) error {
		return middleware.Middleware(ctx, next)
	}
}

// buildRoute create the context for the route and attach the middlewares to it if exists.
func buildRoute(app *App, handler HandlerFunc, rm []Middleware) http.HandlerFunc {
	// Route middlewares
	if len(rm) > 0 {
		for _, m := range rm {
			handler = buildMiddleware(m, handler)
		}
	}

	// Global Middlewares
	if len(app.Middlewares) > 0 {
		for _, m := range app.Middlewares {
			handler = buildMiddleware(m, handler)
		}
	}

	// Assign default middlewares to all handlers.
	for _, defaultMiddleware := range defaultMiddlewares() {
		handler = buildMiddleware(defaultMiddleware, handler)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// Root context.
		ctx := app.pool.Get().(*Context)

		// Reset context before use
		ctx.response = &responseWriter{w, ctx}
		ctx.request = r
		ctx.data = make(map[string]any)
		ctx.statusCode = http.StatusOK
		ctx.startAt = time.Now().UTC()
		ctx.requestID = uuid.New().String()

		// Release the pool
		defer app.pool.Put(ctx)

		// Run the chain
		err := handler(ctx)
		if err != nil {
			handleError(ctx, err)
		}
	}
}

// create new route instance.
func createRoute(app *App, method Method, path string, handler HandlerFunc, middlewares ...Middleware) *Route {
	route := &Route{
		Method:      method,
		Path:        path,
		Middlewares: middlewares,
		HandlerFunc: handler,
	}

	// Append the route
	app.Routes = append(app.Routes, route)

	return route
}

// Get Handle GET requests.
func (app *App) Get(path string, ctx HandlerFunc, middlewares ...Middleware) *Route {
	return createRoute(app, GET, path, ctx, middlewares...)
}

// Head Handle HEAD requests.
func (app *App) Head(path string, ctx HandlerFunc, middlewares ...Middleware) *Route {
	return createRoute(app, HEAD, path, ctx, middlewares...)
}

// Post Handle POST requests.
func (app *App) Post(path string, ctx HandlerFunc, middlewares ...Middleware) *Route {
	return createRoute(app, POST, path, ctx, middlewares...)
}

// Put Handle PUT requests.
func (app *App) Put(path string, ctx HandlerFunc, middlewares ...Middleware) *Route {
	return createRoute(app, PUT, path, ctx, middlewares...)
}

// Patch Handle PATCH requests.
func (app *App) Patch(path string, ctx HandlerFunc, middlewares ...Middleware) *Route {
	return createRoute(app, PATCH, path, ctx, middlewares...)
}

// Delete Handle DELETE requests.
func (app *App) Delete(path string, ctx HandlerFunc, middlewares ...Middleware) *Route {
	return createRoute(app, DELETE, path, ctx, middlewares...)
}

// Options Handle OPTIONS requests.
func (app *App) Options(path string, ctx HandlerFunc, middlewares ...Middleware) *Route {
	return createRoute(app, OPTIONS, path, ctx, middlewares...)
}

// ServeFiles serve files in the given root.
func (app *App) ServeFiles(path string, root http.FileSystem) {
	app.Router.ServeFiles(path, root)
}
