package cosan

import (
	"bytes"
	"net/http/httptest"
	"strings"
	"testing"
)

// FuzzRouterPath tests router with various path inputs
func FuzzRouterPath(f *testing.F) {
	// Seed corpus with known paths
	f.Add("/users")
	f.Add("/users/123")
	f.Add("/api/v1/posts")
	f.Add("/")
	f.Add("")

	f.Fuzz(func(t *testing.T, path string) {
		// Skip empty paths as httptest.NewRequest doesn't accept them
		if path == "" {
			path = "/"
		}

		router := New()
		
		router.GET("/users", func(ctx Context) error {
			return ctx.String(200, "users")
		})
		
		router.GET("/users/:id", func(ctx Context) error {
			return ctx.String(200, "user")
		})

		// Try to match the path - should not panic
		req := httptest.NewRequest("GET", path, nil)
		w := httptest.NewRecorder()
		
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Router panicked on path %q: %v", path, r)
			}
		}()
		
		router.ServeHTTP(w, req)
	})
}

// FuzzJSONInput tests JSON parsing with various inputs
func FuzzJSONInput(f *testing.F) {
	// Seed corpus
	f.Add(`{"name":"test"}`)
	f.Add(`{"id":123}`)
	f.Add(`[]`)
	f.Add(`null`)

	f.Fuzz(func(t *testing.T, jsonData string) {
		router := New()
		
		router.POST("/data", func(ctx Context) error {
			var data map[string]interface{}
			_ = ctx.Bind(&data) // Ignore error, just test for panics
			return ctx.String(200, "OK")
		})

		req := httptest.NewRequest("POST", "/data", bytes.NewBufferString(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Router panicked on JSON %q: %v", jsonData, r)
			}
		}()
		
		router.ServeHTTP(w, req)
	})
}

// FuzzParamNames tests parameter names
func FuzzParamNames(f *testing.F) {
	// Seed corpus
	f.Add("id")
	f.Add("userId")
	f.Add("postID")
	f.Add("a")

	f.Fuzz(func(t *testing.T, paramName string) {
		// Skip invalid param names
		if paramName == "" || strings.ContainsAny(paramName, "/:*") {
			return
		}

		router := New()
		
		pattern := "/items/:" + paramName
		router.GET(pattern, func(ctx Context) error {
			_ = ctx.Param(paramName)
			return ctx.String(200, "OK")
		})

		req := httptest.NewRequest("GET", "/items/123", nil)
		w := httptest.NewRecorder()
		
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Router panicked with param name %q: %v", paramName, r)
			}
		}()
		
		router.ServeHTTP(w, req)
	})
}
