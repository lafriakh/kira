package kira

// TODO:
//  - Implement "plugin" mechanism.
//  - We can use "plugin" to provide additional functionalities to the user like: Auth, Cache, Database ORM...
//  - Error wrapper: Error{op: "op.name", err: Error}...

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"os"
	"sync"
	"time"

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
	// Not found handler
	NotFoundHandler HandlerFunc

	// Logger
	logger *log.Logger

	// Context pool
	pool  *sync.Pool
	mutex sync.Mutex
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

// Run the framework
func (app *App) Run(args ...any) *App {
	fmt.Printf("%v", hero)

	// Register the application routes
	app.RegisterRoutes()

	// Timezone
	tz := app.Configs.GetString("app.timezone")
	if tz != "" {
		os.Setenv("TZ", tz)
	}

	// Server
	server := &http.Server{
		Handler: app.Router,
	}

	var config any
	if len(args) > 0 {
		config = args[0]
	} else {
		config = nil
	}

	switch config := config.(type) {
	case *http.Server:
		server = config
		server.Handler = app.Router
	case string:
		server.Addr = serverAddr(app.Configs, config)
	default:
		server.Addr = serverAddr(app.Configs)
	}

	if !app.Configs.GetBool("server.tls", false) {
		app.StartServer(server)
	} else {
		app.StartTLSServer(server)
	}

	// App instance
	return app
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
		Status:    http.StatusNotFound,
		Message:   http.StatusText(http.StatusNotFound),
	}
}
