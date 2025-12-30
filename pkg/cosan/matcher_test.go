package cosan

import (
	"testing"
)

func TestSimpleMatcher_Register(t *testing.T) {
	m := newSimpleMatcher().(*simpleMatcher)

	handler := func(ctx Context) error { return nil }
	err := m.Register("GET", "/users", handler)
	if err != nil {
		t.Errorf("Failed to register route: %v", err)
	}

	if len(m.routes) != 1 {
		t.Errorf("Expected 1 route, got %d", len(m.routes))
	}
}

func TestSimpleMatcher_Compile(t *testing.T) {
	m := newSimpleMatcher().(*simpleMatcher)

	handler := func(ctx Context) error { return nil }
	m.Register("GET", "/users", handler)

	err := m.Compile()
	if err != nil {
		t.Errorf("Compile failed: %v", err)
	}
}

func TestSimpleMatcher_MatchExact(t *testing.T) {
	m := newSimpleMatcher()
	handler := func(ctx Context) error { return nil }

	m.Register("GET", "/users", handler)
	m.Compile()

	route, params, ok := m.Match("GET", "/users")
	if !ok {
		t.Error("Failed to match exact route")
	}
	if route == nil {
		t.Error("Route is nil")
	}
	if len(params) != 0 {
		t.Errorf("Expected no params, got %d", len(params))
	}
}

func TestSimpleMatcher_MatchNotFound(t *testing.T) {
	m := newSimpleMatcher()
	handler := func(ctx Context) error { return nil }

	m.Register("GET", "/users", handler)
	m.Compile()

	_, _, ok := m.Match("GET", "/posts")
	if ok {
		t.Error("Should not match non-existent route")
	}
}

func TestSimpleMatcher_MatchWrongMethod(t *testing.T) {
	m := newSimpleMatcher()
	handler := func(ctx Context) error { return nil }

	m.Register("GET", "/users", handler)
	m.Compile()

	_, _, ok := m.Match("POST", "/users")
	if ok {
		t.Error("Should not match with wrong method")
	}
}

func TestSimpleMatcher_MultipleRoutes(t *testing.T) {
	m := newSimpleMatcher()
	handler := func(ctx Context) error { return nil }

	m.Register("GET", "/users", handler)
	m.Register("POST", "/users", handler)
	m.Register("GET", "/posts", handler)
	m.Compile()

	// Test all routes
	tests := []struct {
		method  string
		path    string
		shouldMatch bool
	}{
		{"GET", "/users", true},
		{"POST", "/users", true},
		{"GET", "/posts", true},
		{"DELETE", "/users", false},
		{"GET", "/comments", false},
	}

	for _, tt := range tests {
		_, _, ok := m.Match(tt.method, tt.path)
		if ok != tt.shouldMatch {
			t.Errorf("Match(%s, %s) = %v, want %v", tt.method, tt.path, ok, tt.shouldMatch)
		}
	}
}
