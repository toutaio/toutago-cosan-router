package cosan

import (
	"testing"
)

func TestRouteMetadata_WithName(t *testing.T) {
	r := &route{}
	opt := WithName("test-route")
	opt(r)

	if r.metadata == nil {
		t.Fatal("Metadata was not initialized")
	}
	if r.metadata.Name != "test-route" {
		t.Errorf("Expected name 'test-route', got '%s'", r.metadata.Name)
	}
}

func TestRouteMetadata_WithDescription(t *testing.T) {
	r := &route{}
	opt := WithDescription("Test description")
	opt(r)

	if r.metadata == nil {
		t.Fatal("Metadata was not initialized")
	}
	if r.metadata.Description != "Test description" {
		t.Errorf("Expected description 'Test description', got '%s'", r.metadata.Description)
	}
}

func TestRouteMetadata_WithTags(t *testing.T) {
	r := &route{}
	opt := WithTags("tag1", "tag2")
	opt(r)

	if r.metadata == nil {
		t.Fatal("Metadata was not initialized")
	}
	if len(r.metadata.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(r.metadata.Tags))
	}
	if r.metadata.Tags[0] != "tag1" || r.metadata.Tags[1] != "tag2" {
		t.Errorf("Tags don't match: %v", r.metadata.Tags)
	}
}

func TestRouteMetadata_MultipleTags(t *testing.T) {
	r := &route{}
	opt1 := WithTags("tag1")
	opt2 := WithTags("tag2", "tag3")
	opt1(r)
	opt2(r)

	if len(r.metadata.Tags) != 3 {
		t.Errorf("Expected 3 tags, got %d", len(r.metadata.Tags))
	}
}

func TestRouteMetadata_Deprecated(t *testing.T) {
	r := &route{}
	opt := Deprecated()
	opt(r)

	if r.metadata == nil {
		t.Fatal("Metadata was not initialized")
	}
	if !r.metadata.Deprecated {
		t.Error("Route should be marked as deprecated")
	}
}

func TestRouteMetadata_WithVersion(t *testing.T) {
	r := &route{}
	opt := WithVersion("v1.0.0")
	opt(r)

	if r.metadata == nil {
		t.Fatal("Metadata was not initialized")
	}
	if r.metadata.Version != "v1.0.0" {
		t.Errorf("Expected version 'v1.0.0', got '%s'", r.metadata.Version)
	}
}

func TestRouteMetadata_MultipleOptions(t *testing.T) {
	r := &route{}
	WithName("api-endpoint")(r)
	WithDescription("API endpoint for testing")(r)
	WithTags("api", "test")(r)
	Deprecated()(r)
	WithVersion("v2.0.0")(r)

	if r.metadata.Name != "api-endpoint" {
		t.Errorf("Name mismatch")
	}
	if r.metadata.Description != "API endpoint for testing" {
		t.Errorf("Description mismatch")
	}
	if len(r.metadata.Tags) != 2 {
		t.Errorf("Expected 2 tags")
	}
	if !r.metadata.Deprecated {
		t.Error("Should be deprecated")
	}
	if r.metadata.Version != "v2.0.0" {
		t.Errorf("Version mismatch")
	}
}

func TestRouter_GetRoutes(t *testing.T) {
	router := New()

	router.GET("/users", func(ctx Context) error {
		return ctx.String(200, "users")
	})

	router.POST("/users", func(ctx Context) error {
		return ctx.String(200, "create user")
	})

	routes := router.GetRoutes()
	if len(routes) != 2 {
		t.Errorf("Expected 2 routes, got %d", len(routes))
	}

	// Check route methods and patterns
	found := make(map[string]bool)
	for _, r := range routes {
		key := r.Method + ":" + r.Pattern
		found[key] = true
	}

	if !found["GET:/users"] {
		t.Error("GET /users not found")
	}
	if !found["POST:/users"] {
		t.Error("POST /users not found")
	}
}

func TestRouter_GetRoutesWithMetadata(t *testing.T) {
	router := New().(*router)

	// Create route with metadata
	r := &route{
		method:  "GET",
		pattern: "/api/v1/users",
		handler: func(ctx Context) error { return nil },
		metadata: &RouteMetadata{
			Name:        "list-users",
			Description: "Lists all users",
			Tags:        []string{"users", "api"},
			Deprecated:  false,
			Version:     "v1.0.0",
		},
	}
	router.routes = append(router.routes, r)

	routes := router.GetRoutes()
	if len(routes) != 1 {
		t.Fatalf("Expected 1 route, got %d", len(routes))
	}

	info := routes[0]
	if info.Name != "list-users" {
		t.Errorf("Name mismatch: %s", info.Name)
	}
	if info.Description != "Lists all users" {
		t.Errorf("Description mismatch: %s", info.Description)
	}
	if len(info.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(info.Tags))
	}
	if info.Deprecated {
		t.Error("Should not be deprecated")
	}
	if info.Version != "v1.0.0" {
		t.Errorf("Version mismatch: %s", info.Version)
	}
}

func TestRouter_FindRoute(t *testing.T) {
	router := New().(*router)

	r := &route{
		method:  "GET",
		pattern: "/users",
		handler: func(ctx Context) error { return nil },
		metadata: &RouteMetadata{
			Name: "get-users",
		},
	}
	router.routes = append(router.routes, r)

	found := router.FindRoute("get-users")
	if found == nil {
		t.Fatal("Route not found")
	}
	if found.Name != "get-users" {
		t.Errorf("Name mismatch: %s", found.Name)
	}
	if found.Method != "GET" {
		t.Errorf("Method mismatch: %s", found.Method)
	}
	if found.Pattern != "/users" {
		t.Errorf("Pattern mismatch: %s", found.Pattern)
	}
}

func TestRouter_FindRoute_NotFound(t *testing.T) {
	router := New().(*router)

	found := router.FindRoute("nonexistent")
	if found != nil {
		t.Error("Should not find nonexistent route")
	}
}

func TestRouter_FindRoute_NoMetadata(t *testing.T) {
	router := New().(*router)

	r := &route{
		method:  "GET",
		pattern: "/users",
		handler: func(ctx Context) error { return nil },
	}
	router.routes = append(router.routes, r)

	found := router.FindRoute("get-users")
	if found != nil {
		t.Error("Should not find route without metadata name")
	}
}
