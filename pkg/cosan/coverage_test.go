package cosan

import (
	"net/http/httptest"
	"strings"
	"testing"
)

// TestContext_BodyBytes tests reading raw body bytes
func TestContext_BodyBytes(t *testing.T) {
	router := New()

	var bodyContent []byte
	router.POST("/data", func(ctx Context) error {
		var err error
		bodyContent, err = ctx.BodyBytes()
		if err != nil {
			return err
		}
		return ctx.String(200, "OK")
	})

	expected := "test body content"
	req := httptest.NewRequest("POST", "/data", strings.NewReader(expected))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if string(bodyContent) != expected {
		t.Errorf("Expected body %q, got %q", expected, string(bodyContent))
	}
}

// TestContext_HTML tests HTML response
func TestContext_HTML(t *testing.T) {
	router := New()

	router.GET("/page", func(ctx Context) error {
		return ctx.HTML(200, "<h1>Hello</h1>")
	})

	req := httptest.NewRequest("GET", "/page", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if w.Header().Get("Content-Type") != "text/html; charset=utf-8" {
		t.Errorf("Expected HTML content type, got %s", w.Header().Get("Content-Type"))
	}
	if w.Body.String() != "<h1>Hello</h1>" {
		t.Errorf("Unexpected body: %s", w.Body.String())
	}
}

// TestContext_Header tests header access
func TestContext_Header(t *testing.T) {
	router := New()

	router.GET("/test", func(ctx Context) error {
		ctx.Header().Set("X-Custom", "value")
		return ctx.String(200, "OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Header().Get("X-Custom") != "value" {
		t.Error("Custom header not set")
	}
}

// TestContext_Write tests direct write
func TestContext_Write(t *testing.T) {
	router := New()

	router.GET("/binary", func(ctx Context) error {
		data := []byte{0x01, 0x02, 0x03}
		n, err := ctx.Write(data)
		if err != nil {
			return err
		}
		if n != 3 {
			return ctx.String(500, "write error")
		}
		return nil
	})

	req := httptest.NewRequest("GET", "/binary", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if len(w.Body.Bytes()) != 3 {
		t.Errorf("Expected 3 bytes, got %d", len(w.Body.Bytes()))
	}
}

// TestRoute_Interface tests Route interface methods
func TestRoute_Interface(t *testing.T) {
	handler := func(ctx Context) error { return nil }
	r := &route{
		method:  "GET",
		pattern: "/test",
		handler: handler,
	}

	if r.Pattern() != "/test" {
		t.Errorf("Pattern() = %s, want /test", r.Pattern())
	}

	if r.Method() != "GET" {
		t.Errorf("Method() = %s, want GET", r.Method())
	}

	if r.Handler() == nil {
		t.Error("Handler() returned nil")
	}
}

// TestNew_WithDefaults tests router initialization
func TestNew_WithDefaults(t *testing.T) {
	r := New()
	if r == nil {
		t.Fatal("New() returned nil")
	}

	// Try using the router
	r.GET("/test", func(ctx Context) error {
		return ctx.String(200, "OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

// TestMatcher_DuplicateRoute tests duplicate route detection
func TestMatcher_DuplicateRoute(t *testing.T) {
	m := newSimpleMatcher()
	handler := func(ctx Context) error { return nil }

	// Register first route
	err := m.Register("GET", "/test", handler)
	if err != nil {
		t.Errorf("First registration failed: %v", err)
	}

	// Try to register duplicate
	err = m.Register("GET", "/test", handler)
	if err == nil {
		t.Error("Expected error for duplicate route")
	}
}

// TestMatcher_RegisterAfterCompile tests registration after compilation
func TestMatcher_RegisterAfterCompile(t *testing.T) {
	m := newSimpleMatcher()
	handler := func(ctx Context) error { return nil }

	m.Register("GET", "/test", handler)
	m.Compile()

	// Try to register after compile
	err := m.Register("GET", "/new", handler)
	if err == nil {
		t.Error("Expected error when registering after compile")
	}
}

// TestStatusRecorder tests status code capture
func TestStatusRecorder(t *testing.T) {
	w := httptest.NewRecorder()
	sr := &statusRecorder{ResponseWriter: w, statusCode: 200}

	// Test WriteHeader
	sr.WriteHeader(404)
	if sr.statusCode != 404 {
		t.Errorf("Expected status 404, got %d", sr.statusCode)
	}

	// Multiple WriteHeader calls should only apply first
	sr.WriteHeader(500)
	if sr.statusCode != 404 {
		t.Errorf("Status should remain 404, got %d", sr.statusCode)
	}
}

// TestStatusRecorder_Write tests implicit status setting
func TestStatusRecorder_Write(t *testing.T) {
	w := httptest.NewRecorder()
	sr := &statusRecorder{ResponseWriter: w, statusCode: 200}

	// Write should set status to 200 if not set
	sr.Write([]byte("test"))
	if sr.statusCode != 200 {
		t.Errorf("Expected status 200, got %d", sr.statusCode)
	}
	if !sr.written {
		t.Error("written flag should be true")
	}
}
