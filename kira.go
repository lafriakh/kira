package kira

// TODO:
//  - Implement "plugin" mechanism.
//  - We can use "plugin" to provide additional functionalities to the user like: Auth, Cache, Database ORM...
//  - Error wrapper: Error{op: "op.name", err: Error}...

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/julienschmidt/httprouter"
	"github.com/lafriakh/kira/modules/config"
	"github.com/lafriakh/kira/modules/log"
)

var hero = `   __ __   _             
  / //_/  (_)  ____ ___ _
 / ,<    / /  / __// _  /
/_/|_|  /_/  /_/   \_,_/ 
`

// some bytes :)
const (
	KB = 1 << 10
	MB = 1 << 20
	GB = 1 << 30
)

// Map a type to represent map, this will be used alot in the internal statusCode.
type Map map[string]any

// App hold the framework options
type App struct {
	Routes      []*Route
	Middlewares []Middleware
	Router      *httprouter.Router
	Configs     *config.Config
	Env         string
	lifecyles   struct {
		onStart    func()
		onShutdown func()
	}
	// Not found handler
	NotFoundHandler HandlerFunc

	// Logger
	logger *log.Logger

	// Context pool
	pool *sync.Pool
}

// New init the framework
func New() *App {
	app := &App{}
	app.Env = getEnv()
	app.Configs = getConfig()
	app.Router = httprouter.New()
	app.logger = setupLogger(app.Configs, setupWriter(app.Configs), log.Fields{})

	// Context pool
	app.pool = &sync.Pool{
		New: func() any {
			return &Context{
				logger:     app.logger,
				configs:    app.Configs,
				data:       make(map[string]any),
				env:        app.Env,
				statusCode: http.StatusOK,
				requestID:  uuid.New().String(),
				startAt:    time.Now().UTC(),
			}
		},
	}

	// return App instance
	return app
}
func (app *App) Logger() *log.Logger {
	return app.logger
}

// Run the framework
func (app *App) Run(args ...any) *App {
	fmt.Printf("%v", hero)

	// Register the application routes
	app.RegisterRoutes()

	// Timezone
	tz := app.Configs.GetString("app.timezone")
	if tz != "" {
		if err := os.Setenv("TZ", tz); err != nil {
			fmt.Fprint(os.Stderr, err)
		}
	}

	// Server
	server := &http.Server{
		Handler: app.Router,
		Addr:    getServerAddr(app.Configs),
	}

	app.StartServer(server)

	// App instance
	return app
}

func (app *App) OnStart(handler func()) {
	app.lifecyles.onStart = handler
}

func (app *App) OnShutdown(handler func()) {
	app.lifecyles.onShutdown = handler
}

// NotFound custom not found handler.
func (app *App) NotFound(ctx HandlerFunc) {
	app.NotFoundHandler = ctx
}

// default not found handler.
func defaultNotFound(ctx *Context) error {
	if ctx.WantsJSON() {
		ctx.Response().Header().Set("Content-Type", "application/json")
	} else {
		ctx.Response().Header().Set("Content-Type", "text/html")
	}
	ctx.SetStatusCode(http.StatusNotFound)

	return &Error{
		Status:  http.StatusNotFound,
		Message: http.StatusText(http.StatusNotFound),
	}
}

func noopHandler(ctx *Context) error {
	return nil
}
