package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lafriakh/kira"
	"github.com/lafriakh/kira/modules/middlewares/cors"
)

func TestCORSDefault(t *testing.T) {
	app := kira.New()
	app.Use(cors.Default())

	app.Get("/test", func(c *kira.Context) error {
		c.WriteString("OK")
		return nil
	})

	s := httptest.NewServer(app.RegisterRoutes())
	defer s.Close()

	// Test regular request
	req, _ := http.NewRequest("GET", s.URL+"/test", nil)
	req.Header.Set("Origin", "http://example.com")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	// Check CORS headers
	if origin := res.Header.Get("Access-Control-Allow-Origin"); origin != "http://example.com" {
		t.Errorf("expected Access-Control-Allow-Origin header 'http://example.com', got %q", origin)
	}
}

func TestCORSPreflight(t *testing.T) {
	app := kira.New()
	app.Use(cors.Default())

	app.Get("/test", func(c *kira.Context) error {
		c.WriteString("OK")
		return nil
	})

	s := httptest.NewServer(app.RegisterRoutes())
	defer s.Close()

	// Preflight request
	req, _ := http.NewRequest("OPTIONS", s.URL+"/test", nil)
	req.Header.Set("Origin", "http://example.com")
	req.Header.Set("Access-Control-Request-Method", "GET")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 200 or 204, got %d", res.StatusCode)
	}
	if origin := res.Header.Get("Access-Control-Allow-Origin"); origin != "http://example.com" {
		t.Errorf("expected Access-Control-Allow-Origin header 'http://example.com', got %q", origin)
	}
}
