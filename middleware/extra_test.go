package middleware

import (
	"testing"

	"github.com/toutaio/toutago-cosan-router"
)

// TestMiddleware_CORSAllOrigins tests CORS with all origins
func TestMiddleware_CORSAllOrigins(t *testing.T) {
	mw := CORS()

	handler := cosan.MiddlewareFunc(func(next cosan.HandlerFunc) cosan.HandlerFunc {
		return func(ctx cosan.Context) error {
			return ctx.String(200, "OK")
		}
	})

	processed := mw.Process(handler.Process(nil))
	// Just test that it doesn't panic
	if processed == nil {
		t.Error("Expected non-nil handler")
	}
}

// TestMiddleware_CORSCustomConfig tests CORS with custom config
func TestMiddleware_CORSCustomConfig(t *testing.T) {
	config := CORSConfig{
		AllowOrigins: []string{"https://example.com", "https://api.example.com"},
		AllowMethods: []string{"GET", "POST"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
	}
	mw := CORSWithConfig(config)

	handler := cosan.MiddlewareFunc(func(next cosan.HandlerFunc) cosan.HandlerFunc {
		return func(ctx cosan.Context) error {
			return ctx.String(200, "OK")
		}
	})

	processed := mw.Process(handler.Process(nil))
	if processed == nil {
		t.Error("Expected non-nil handler")
	}
}
