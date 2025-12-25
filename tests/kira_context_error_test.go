package tests

import (
	"net/http"
	"testing"

	"github.com/lafriakh/kira"
	"github.com/stretchr/testify/assert"
)

var ErrContextError = &kira.Error{Message: "context error", Path: kira.Path("/path"), Status: 400}

func TestContextError(t *testing.T) {
	s := endpoint("GET", "/error", func(c *kira.Context) error {
		return ErrContextError
	})

	// Request
	req, _ := http.NewRequest(http.MethodGet, url(s, "/error"), nil)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json")
	// Send request.
	res, _ := http.DefaultClient.Do(req)

	content := contentS(res.Body)
	// Note: The exact JSON output depends on the JSON marshaler.
	// We'll check for the presence of the key fields.
	assert.Contains(t, content, `"status":400`)
	assert.Contains(t, content, `"path":"/path"`)
	assert.Contains(t, content, `"message":"context error"`)
}
