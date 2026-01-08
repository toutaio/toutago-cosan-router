package cosan

import (
	"net/http/httptest"
	"testing"
)

// TestContext_Set_Get tests context value storage
func TestContext_Set_Get(t *testing.T) {
	router := New()

	router.GET("/test", func(ctx Context) error {
		ctx.Set("user", "john")
		ctx.Set("role", "admin")

		if ctx.Get("user") != "john" {
			return ctx.String(500, "Get failed")
		}
		if ctx.Get("role") != "admin" {
			return ctx.String(500, "Get failed")
		}
		if ctx.Get("missing") != nil {
			return ctx.String(500, "Should return nil for missing key")
		}

		return ctx.String(200, "OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

// TestRouterGroup_Methods tests all HTTP methods on groups
func TestRouterGroup_Methods(t *testing.T) {
	router := New()
	api := router.Group("/api")

	api.GET("/test", func(ctx Context) error { return ctx.String(200, "GET") })
	api.POST("/test", func(ctx Context) error { return ctx.String(200, "POST") })
	api.PUT("/test", func(ctx Context) error { return ctx.String(200, "PUT") })
	api.DELETE("/test", func(ctx Context) error { return ctx.String(200, "DELETE") })
	api.PATCH("/test", func(ctx Context) error { return ctx.String(200, "PATCH") })
	api.OPTIONS("/test", func(ctx Context) error { return ctx.String(200, "OPTIONS") })
	api.HEAD("/test", func(ctx Context) error { ctx.Status(200); return nil })

	tests := []struct {
		method string
		expect string
	}{
		{"GET", "GET"},
		{"POST", "POST"},
		{"PUT", "PUT"},
		{"DELETE", "DELETE"},
		{"PATCH", "PATCH"},
		{"OPTIONS", "OPTIONS"},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(tt.method, "/api/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Errorf("%s: expected 200, got %d", tt.method, w.Code)
		}
		if w.Body.String() != tt.expect {
			t.Errorf("%s: expected %s, got %s", tt.method, tt.expect, w.Body.String())
		}
	}

	// Test HEAD
	req := httptest.NewRequest("HEAD", "/api/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("HEAD: expected 200, got %d", w.Code)
	}
}
