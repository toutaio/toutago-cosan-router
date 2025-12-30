package cosan_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/toutaio/toutago-cosan-router/pkg/cosan"
)

// TestBasicRouterCreation tests router instantiation.
func TestBasicRouterCreation(t *testing.T) {
	router := cosan.New()
	if router == nil {
		t.Fatal("Expected router to be created, got nil")
	}
}

// TestMethodBasedRouting tests all HTTP method registration.
func TestMethodBasedRouting(t *testing.T) {
	router := cosan.New()

	// Register routes for all HTTP methods
	router.GET("/get", func(ctx cosan.Context) error {
		return ctx.String(200, "GET")
	})
	router.POST("/post", func(ctx cosan.Context) error {
		return ctx.String(200, "POST")
	})
	router.PUT("/put", func(ctx cosan.Context) error {
		return ctx.String(200, "PUT")
	})
	router.DELETE("/delete", func(ctx cosan.Context) error {
		return ctx.String(200, "DELETE")
	})
	router.PATCH("/patch", func(ctx cosan.Context) error {
		return ctx.String(200, "PATCH")
	})
	router.OPTIONS("/options", func(ctx cosan.Context) error {
		return ctx.String(200, "OPTIONS")
	})
	router.HEAD("/head", func(ctx cosan.Context) error {
		ctx.Status(200)
		return nil
	})

	tests := []struct {
		method string
		path   string
		want   string
	}{
		{http.MethodGet, "/get", "GET"},
		{http.MethodPost, "/post", "POST"},
		{http.MethodPut, "/put", "PUT"},
		{http.MethodDelete, "/delete", "DELETE"},
		{http.MethodPatch, "/patch", "PATCH"},
		{http.MethodOptions, "/options", "OPTIONS"},
		{http.MethodHead, "/head", ""},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != 200 {
				t.Errorf("Expected status 200, got %d", w.Code)
			}

			if tt.method != http.MethodHead && w.Body.String() != tt.want {
				t.Errorf("Expected body %q, got %q", tt.want, w.Body.String())
			}
		})
	}
}

// TestExactPathMatching tests exact path matching.
func TestExactPathMatching(t *testing.T) {
	router := cosan.New()

	router.GET("/users", func(ctx cosan.Context) error {
		return ctx.String(200, "users list")
	})
	router.GET("/users/123", func(ctx cosan.Context) error {
		return ctx.String(200, "user 123")
	})

	tests := []struct {
		path       string
		wantStatus int
		wantBody   string
	}{
		{"/users", 200, "users list"},
		{"/users/123", 200, "user 123"},
		{"/users/456", 404, "404 page not found\n"},
		{"/nonexistent", 404, "404 page not found\n"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Expected status %d, got %d", tt.wantStatus, w.Code)
			}

			if w.Body.String() != tt.wantBody {
				t.Errorf("Expected body %q, got %q", tt.wantBody, w.Body.String())
			}
		})
	}
}

// TestJSONResponse tests JSON response handling.
func TestJSONResponse(t *testing.T) {
	router := cosan.New()

	type User struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	router.GET("/user", func(ctx cosan.Context) error {
		return ctx.JSON(200, User{ID: 1, Name: "Alice"})
	})

	req := httptest.NewRequest(http.MethodGet, "/user", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	var user User
	if err := json.NewDecoder(w.Body).Decode(&user); err != nil {
		t.Fatalf("Failed to decode JSON: %v", err)
	}

	if user.ID != 1 || user.Name != "Alice" {
		t.Errorf("Expected user {1, Alice}, got {%d, %s}", user.ID, user.Name)
	}
}

// TestJSONBinding tests JSON request body binding.
func TestJSONBinding(t *testing.T) {
	router := cosan.New()

	type CreateUser struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	router.POST("/users", func(ctx cosan.Context) error {
		var user CreateUser
		if err := ctx.Bind(&user); err != nil {
			return err
		}
		return ctx.JSON(201, user)
	})

	body := bytes.NewBufferString(`{"name":"Bob","email":"bob@example.com"}`)
	req := httptest.NewRequest(http.MethodPost, "/users", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 201 {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var user CreateUser
	if err := json.NewDecoder(w.Body).Decode(&user); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if user.Name != "Bob" || user.Email != "bob@example.com" {
		t.Errorf("Expected user {Bob, bob@example.com}, got {%s, %s}", user.Name, user.Email)
	}
}

// TestQueryParameters tests query parameter handling.
func TestQueryParameters(t *testing.T) {
	router := cosan.New()

	router.GET("/search", func(ctx cosan.Context) error {
		q := ctx.Query("q")
		tags := ctx.QueryAll("tag")
		return ctx.JSON(200, map[string]interface{}{
			"query": q,
			"tags":  tags,
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/search?q=golang&tag=web&tag=api", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if result["query"] != "golang" {
		t.Errorf("Expected query 'golang', got %v", result["query"])
	}

	tags, ok := result["tags"].([]interface{})
	if !ok || len(tags) != 2 {
		t.Errorf("Expected 2 tags, got %v", result["tags"])
	}
}

// TestRouteConflictDetection tests duplicate route detection.
func TestRouteConflictDetection(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for duplicate route, but didn't panic")
		}
	}()

	router := cosan.New()
	router.GET("/users", func(ctx cosan.Context) error {
		return ctx.String(200, "first")
	})
	// This should panic
	router.GET("/users", func(ctx cosan.Context) error {
		return ctx.String(200, "duplicate")
	})
}

// TestRouteGroups tests route grouping functionality.
func TestRouteGroups(t *testing.T) {
	router := cosan.New()

	// Create API v1 group
	v1 := router.Group("/api/v1")
	v1.GET("/users", func(ctx cosan.Context) error {
		return ctx.String(200, "v1 users")
	})
	v1.POST("/users", func(ctx cosan.Context) error {
		return ctx.String(201, "v1 create user")
	})

	// Create API v2 group
	v2 := router.Group("/api/v2")
	v2.GET("/users", func(ctx cosan.Context) error {
		return ctx.String(200, "v2 users")
	})

	// Nested groups
	admin := v1.Group("/admin")
	admin.GET("/stats", func(ctx cosan.Context) error {
		return ctx.String(200, "admin stats")
	})

	tests := []struct {
		method     string
		path       string
		wantStatus int
		wantBody   string
	}{
		{http.MethodGet, "/api/v1/users", 200, "v1 users"},
		{http.MethodPost, "/api/v1/users", 201, "v1 create user"},
		{http.MethodGet, "/api/v2/users", 200, "v2 users"},
		{http.MethodGet, "/api/v1/admin/stats", 200, "admin stats"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Expected status %d, got %d", tt.wantStatus, w.Code)
			}

			if w.Body.String() != tt.wantBody {
				t.Errorf("Expected body %q, got %q", tt.wantBody, w.Body.String())
			}
		})
	}
}

// TestContextValueStorage tests context value storage.
func TestContextValueStorage(t *testing.T) {
	router := cosan.New()

	router.GET("/context", func(ctx cosan.Context) error {
		ctx.Set("user", "Alice")
		ctx.Set("role", "admin")

		user := ctx.Get("user")
		role := ctx.Get("role")
		missing := ctx.Get("missing")

		return ctx.JSON(200, map[string]interface{}{
			"user":    user,
			"role":    role,
			"missing": missing,
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/context", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var result map[string]interface{}
	json.NewDecoder(w.Body).Decode(&result)

	if result["user"] != "Alice" {
		t.Errorf("Expected user 'Alice', got %v", result["user"])
	}
	if result["role"] != "admin" {
		t.Errorf("Expected role 'admin', got %v", result["role"])
	}
	if result["missing"] != nil {
		t.Errorf("Expected missing to be nil, got %v", result["missing"])
	}
}

// TestHTTPHandlerCompliance tests that router implements http.Handler.
func TestHTTPHandlerCompliance(t *testing.T) {
	router := cosan.New()
	router.GET("/", func(ctx cosan.Context) error {
		return ctx.String(200, "Hello, World!")
	})

	// Use with httptest.Server
	server := httptest.NewServer(router)
	defer server.Close()

	resp, err := http.Get(server.URL + "/")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if string(body) != "Hello, World!" {
		t.Errorf("Expected body 'Hello, World!', got %q", string(body))
	}
}

// TestErrorHandling tests error handling in handlers.
func TestErrorHandling(t *testing.T) {
	router := cosan.New()

	router.GET("/error", func(ctx cosan.Context) error {
		return fmt.Errorf("something went wrong")
	})

	req := httptest.NewRequest(http.MethodGet, "/error", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 500 {
		t.Errorf("Expected status 500, got %d", w.Code)
	}

	body := w.Body.String()
	if !bytes.Contains([]byte(body), []byte("something went wrong")) {
		t.Errorf("Expected error message in body, got %q", body)
	}
}

// Benchmark tests
func BenchmarkSimpleRoute(b *testing.B) {
	router := cosan.New()
	router.GET("/users", func(ctx cosan.Context) error {
		return ctx.String(200, "users")
	})

	req := httptest.NewRequest(http.MethodGet, "/users", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}

func BenchmarkJSONResponse(b *testing.B) {
	router := cosan.New()

	type Response struct {
		ID      int    `json:"id"`
		Message string `json:"message"`
	}

	router.GET("/json", func(ctx cosan.Context) error {
		return ctx.JSON(200, Response{ID: 1, Message: "Hello"})
	})

	req := httptest.NewRequest(http.MethodGet, "/json", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}
