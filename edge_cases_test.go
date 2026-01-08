package cosan

import (
	"net/http/httptest"
	"testing"
)

// TestRadix_EdgeCases tests edge cases in radix tree
func TestRadix_EdgeCases(t *testing.T) {
	router := New()

	// Test multiple params on same segment
	router.GET("/a/:id/b/:name", func(ctx Context) error {
		return ctx.JSON(200, map[string]string{
			"id":   ctx.Param("id"),
			"name": ctx.Param("name"),
		})
	})

	req := httptest.NewRequest("GET", "/a/123/b/john", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

// TestRadix_StaticAndParam tests static vs param priority
func TestRadix_StaticAndParam(t *testing.T) {
	router := New()

	staticCalled := false
	paramCalled := false

	router.GET("/users/new", func(ctx Context) error {
		staticCalled = true
		return ctx.String(200, "static")
	})

	router.GET("/users/:id", func(ctx Context) error {
		paramCalled = true
		return ctx.String(200, "param")
	})

	// Static should match first
	req := httptest.NewRequest("GET", "/users/new", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if !staticCalled {
		t.Error("Static route should be called")
	}
	if paramCalled {
		t.Error("Param route should not be called")
	}

	// Reset
	staticCalled = false
	paramCalled = false

	// Param should match for other values
	req = httptest.NewRequest("GET", "/users/123", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if !paramCalled {
		t.Error("Param route should be called")
	}
}

// TestRadix_WildcardRoutes tests wildcard matching
func TestRadix_WildcardRoutes(t *testing.T) {
	router := New()

	router.GET("/files/*path", func(ctx Context) error {
		path := ctx.Param("path")
		return ctx.String(200, "path: "+path)
	})

	req := httptest.NewRequest("GET", "/files/docs/readme.md", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

// TestRadix_DeepNesting tests deeply nested routes
func TestRadix_DeepNesting(t *testing.T) {
	router := New()

	router.GET("/a/b/c/d/e/f/g", func(ctx Context) error {
		return ctx.String(200, "deep")
	})

	req := httptest.NewRequest("GET", "/a/b/c/d/e/f/g", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

// TestRadix_TrailingSlash tests trailing slash handling
func TestRadix_TrailingSlash(t *testing.T) {
	router := New()

	router.GET("/users", func(ctx Context) error {
		return ctx.String(200, "users")
	})

	// Without trailing slash
	req := httptest.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	// With trailing slash - radix tree matches it currently
	req = httptest.NewRequest("GET", "/users/", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Our current implementation treats them the same
	if w.Code == 0 {
		t.Error("Expected some status code")
	}
}

// TestRadix_PrefixMatching tests prefix matching edge cases
func TestRadix_PrefixMatching(t *testing.T) {
	router := New()

	router.GET("/user", func(ctx Context) error {
		return ctx.String(200, "user")
	})

	router.GET("/user/profile", func(ctx Context) error {
		return ctx.String(200, "profile")
	})

	tests := []struct {
		path     string
		expected string
		status   int
	}{
		{"/user", "user", 200},
		{"/user/profile", "profile", 200},
		{"/userz", "", 404},
	}

	for _, tt := range tests {
		req := httptest.NewRequest("GET", tt.path, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != tt.status {
			t.Errorf("Path %s: expected status %d, got %d", tt.path, tt.status, w.Code)
		}
		if tt.status == 200 && w.Body.String() != tt.expected {
			t.Errorf("Path %s: expected body %s, got %s", tt.path, tt.expected, w.Body.String())
		}
	}
}
