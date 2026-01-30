package cors_test

import (
	"net/http/httptest"
	"testing"

	"github.com/lafriakh/kira"
	cors "github.com/lafriakh/kira/modules/middlewares/cors"
	rsCors "github.com/rs/cors"
)

func Example() {
	app := kira.New()
	// Use default CORS middleware (allows all origins)
	app.Middleware(cors.AllowAll())

	app.Get("/", func(c *kira.Context) error {
		c.WriteString("Hello, world!")
		return nil
	})

	// Start server (in real usage)
	// app.Run(":8080")
}

func Example_withOptions() {
	app := kira.New()
	// Use custom CORS options
	app.Middleware(cors.New(rsCors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "PUT"},
		AllowedHeaders: []string{"Content-Type"},
		Debug:          true,
	}))

	app.Get("/api", func(c *kira.Context) error {
		c.JSON(map[string]string{"message": "ok"}, 200)
		return nil
	})

	// Start server (in real usage)
	// app.Run(":8080")
}

func TestCORSExample(t *testing.T) {
	app := kira.New()
	app.Middleware(cors.AllowAll())

	app.Get("/test", func(c *kira.Context) error {
		c.WriteString("OK")
		return nil
	})

	// Simulate a request with Origin header
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/test", nil)
	r.Header.Set("Origin", "http://example.com")

	// Serve the request
	handler := app.RegisterRoutes()
	handler.ServeHTTP(w, r)

	// Check CORS header
	if allowOrigin := w.Header().Get("Access-Control-Allow-Origin"); allowOrigin != "*" && allowOrigin != "http://example.com" {
		t.Errorf("CORS header missing, got %q", allowOrigin)
	}
}
