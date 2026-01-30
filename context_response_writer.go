package kira

import "net/http"

// responseWriter for kira framework.
type responseWriter struct {
	http.ResponseWriter
	ctx      *Context
}

// WriteHeader - store the header to use it later.
func (r *responseWriter) WriteHeader(code int) {
	r.ctx.SetStatusCode(code)
	r.ResponseWriter.WriteHeader(code)
}
