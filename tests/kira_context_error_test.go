package tests

import (
	"net/http"
	"testing"

	"github.com/lafriakh/kira"
	"github.com/stretchr/testify/assert"
)

var ErrContextError = kira.NewKiraError("context error").WithPath("/path")

func TestContextError(t *testing.T) {
	s := endpoint("GET", "/error", func(c *kira.Context) {
		c.Error(ErrContextError)
	})

	// Request
	req, _ := http.NewRequest(http.MethodGet, url(s, "/error"), nil)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json")
	// Send request.
	res, _ := http.DefaultClient.Do(req)

	content := contentS(res.Body)
	assert.Equal(t, `{"message":"context error","path":"/path"}`, content)
}
