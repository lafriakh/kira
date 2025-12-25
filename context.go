package kira

import (
	"net/http"
	"time"

	"github.com/lafriakh/kira/modules/log"

	"github.com/lafriakh/kira/modules/config"
)

// HandlerFunc - Type to define context function
type HandlerFunc func(*Context) error

// Context ...
type Context struct {
	request  *http.Request
	response http.ResponseWriter
	logger   *log.Logger
	configs  *config.Config
	// The data associated with the request.
	data       map[string]any
	statusCode int
	requestID  string
	startAt    time.Time
	// environment
	env string
}

// SetRequest change the current request with the given one.
func (c *Context) SetRequest(r *http.Request) {
	c.request = r
}

// SetResponse change the current response with the given one.
func (c *Context) SetResponse(w http.ResponseWriter) {
	c.response = w
}

// Request a Request represents an HTTP request received by a server.
func (c *Context) Request() *http.Request {
	return c.request
}

// Response is used by an HTTP handler to construct an HTTP response.
func (c *Context) Response() http.ResponseWriter {
	return c.response
}

// Redirect replies to the request with a redirect to url,
func (c *Context) Redirect(url string, code int) {
	http.Redirect(c.Response(), c.Request(), url, code)
}

// Log gets the Log instance.
func (c *Context) Log() *log.Logger {
	return setupLogger(c.Config(), c.logger.Writer, log.Fields{
		"method":     c.Request().Method,
		"path":       c.Request().RequestURI,
		"duration":   time.Since(c.startAt).String(),
		"request_id": c.RequestID(),
	})
}

// Config gets the application configs.
func (c *Context) Config() *config.Config {
	return c.configs
}

// Code sets response status statusCode.
func (c *Context) SetStatusCode(code int) {
	c.statusCode = code
}

// Code gets response status statusCode.
func (c *Context) StatusCode() int {
	return c.statusCode
}

// Code sets response status statusCode.
func (c *Context) SetRequestID(id string) {
	c.requestID = id
}

// Code gets response status statusCode.
func (c *Context) RequestID() string {
	return c.requestID
}

// Env gets the application environment.
func (c *Context) Env() string {
	return c.env
}
