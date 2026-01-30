package cors

import (
	"net/http/httptest"
	"testing"

	"github.com/lafriakh/kira"
)

func TestCORSDefault(t *testing.T) {
	cors := Default()
	ctx := &kira.Context{}
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Origin", "http://example.com")
	w := httptest.NewRecorder()
	ctx.SetRequest(req)
	ctx.SetResponse(w)

	var nextCalled bool
	err := cors.Middleware(ctx, func(c *kira.Context) error {
		nextCalled = true
		return nil
	})
	if err != nil {
		t.Fatalf("middleware error: %v", err)
	}
	if !nextCalled {
		t.Error("next handler not called")
	}
	if allowOrigin := w.Header().Get("Access-Control-Allow-Origin"); allowOrigin != "*" && allowOrigin != "http://example.com" {
		t.Errorf("expected CORS header '*' or 'http://example.com', got %q, headers: %v", allowOrigin, w.Header())
	}
}

func TestCORSPreflight(t *testing.T) {
	cors := Default()
	ctx := &kira.Context{}
	req := httptest.NewRequest("OPTIONS", "/", nil)
	req.Header.Set("Origin", "http://example.com")
	req.Header.Set("Access-Control-Request-Method", "GET")
	w := httptest.NewRecorder()
	ctx.SetRequest(req)
	ctx.SetResponse(w)

	var nextCalled bool
	err := cors.Middleware(ctx, func(c *kira.Context) error {
		nextCalled = true
		return nil
	})
	if err != nil {
		t.Fatalf("middleware error: %v", err)
	}
	if nextCalled {
		t.Error("next handler should not be called for preflight")
	}
	if allowOrigin := w.Header().Get("Access-Control-Allow-Origin"); allowOrigin != "*" && allowOrigin != "http://example.com" {
		t.Errorf("expected CORS header '*' or 'http://example.com', got %q, headers: %v", allowOrigin, w.Header())
	}
}
