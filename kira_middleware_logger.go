package kira

import (
	"net/http"
	"strconv"
)

// Log - log middleware
type Log struct{}

// New ...
func NewLogger() *Log {
	return &Log{}
}

// Middleware handler.
func (l *Log) Middleware(ctx *Context, next HandlerFunc) error {
	err := next(ctx)
	status := ctx.StatusCode()
	if status == 0 {
		status = http.StatusOK
	}

	ctx.Log().WithField(
		"status", strconv.Itoa(status),
	).Info("Request")

	return err
}
