package cosan

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestIntegration_FullRESTAPI tests a complete REST API scenario
func TestIntegration_FullRESTAPI(t *testing.T) {
	router := New()

	// In-memory store
	users := make(map[string]map[string]interface{})

	// List users
	router.GET("/users", func(ctx Context) error {
		list := make([]map[string]interface{}, 0, len(users))
		for _, u := range users {
			list = append(list, u)
		}
		return ctx.JSON(200, list)
	})

	// Create user
	router.POST("/users", func(ctx Context) error {
		var user map[string]interface{}
		if err := ctx.Bind(&user); err != nil {
			return err
		}
		id := user["id"].(string)
		users[id] = user
		return ctx.JSON(201, user)
	})

	// Get user
	router.GET("/users/:id", func(ctx Context) error {
		id := ctx.Param("id")
		user, ok := users[id]
		if !ok {
			return ctx.String(404, "Not Found")
		}
		return ctx.JSON(200, user)
	})

	// Update user
	router.PUT("/users/:id", func(ctx Context) error {
		id := ctx.Param("id")
		var user map[string]interface{}
		if err := ctx.Bind(&user); err != nil {
			return err
		}
		users[id] = user
		return ctx.JSON(200, user)
	})

	// Delete user
	router.DELETE("/users/:id", func(ctx Context) error {
		id := ctx.Param("id")
		delete(users, id)
		ctx.Status(204)
		return nil
	})

	// Test CREATE
	userData := map[string]interface{}{"id": "1", "name": "John"}
	body, _ := json.Marshal(userData)
	req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 201 {
		t.Errorf("POST expected 201, got %d", w.Code)
	}

	// Test READ
	req = httptest.NewRequest("GET", "/users/1", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("GET expected 200, got %d", w.Code)
	}

	// Test UPDATE
	updateData := map[string]interface{}{"id": "1", "name": "Jane"}
	body, _ = json.Marshal(updateData)
	req = httptest.NewRequest("PUT", "/users/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("PUT expected 200, got %d", w.Code)
	}

	// Test DELETE
	req = httptest.NewRequest("DELETE", "/users/1", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 204 {
		t.Errorf("DELETE expected 204, got %d", w.Code)
	}

	// Verify deletion
	req = httptest.NewRequest("GET", "/users/1", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 404 {
		t.Errorf("GET after DELETE expected 404, got %d", w.Code)
	}
}

// TestIntegration_NestedRoutes tests complex nested routing
func TestIntegration_NestedRoutes(t *testing.T) {
	router := New()

	api := router.Group("/api")
	v1 := api.Group("/v1")
	users := v1.Group("/users")

	users.GET("/:id/posts/:postId/comments/:commentId", func(ctx Context) error {
		return ctx.JSON(200, map[string]string{
			"userId":    ctx.Param("id"),
			"postId":    ctx.Param("postId"),
			"commentId": ctx.Param("commentId"),
		})
	})

	req := httptest.NewRequest("GET", "/api/v1/users/1/posts/2/comments/3", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	var result map[string]string
	json.Unmarshal(w.Body.Bytes(), &result)

	if result["userId"] != "1" || result["postId"] != "2" || result["commentId"] != "3" {
		t.Errorf("Parameters not extracted correctly: %v", result)
	}
}

// TestIntegration_MiddlewareChain tests middleware execution order
func TestIntegration_MiddlewareChain(t *testing.T) {
	router := New()

	var order []string

	mw1 := MiddlewareFunc(func(next HandlerFunc) HandlerFunc {
		return func(ctx Context) error {
			order = append(order, "mw1-before")
			err := next(ctx)
			order = append(order, "mw1-after")
			return err
		}
	})

	mw2 := MiddlewareFunc(func(next HandlerFunc) HandlerFunc {
		return func(ctx Context) error {
			order = append(order, "mw2-before")
			err := next(ctx)
			order = append(order, "mw2-after")
			return err
		}
	})

	router.Use(mw1, mw2)

	router.GET("/test", func(ctx Context) error {
		order = append(order, "handler")
		return ctx.String(200, "OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	expected := []string{"mw1-before", "mw2-before", "handler", "mw2-after", "mw1-after"}
	if len(order) != len(expected) {
		t.Fatalf("Expected %d elements, got %d", len(expected), len(order))
	}

	for i, exp := range expected {
		if order[i] != exp {
			t.Errorf("At position %d: expected %s, got %s", i, exp, order[i])
		}
	}
}

// TestIntegration_ContentNegotiation tests different content types
func TestIntegration_ContentNegotiation(t *testing.T) {
	router := New()

	router.POST("/data", func(ctx Context) error {
		contentType := ctx.Request().Header.Get("Content-Type")

		var data map[string]interface{}
		if err := ctx.Bind(&data); err != nil {
			return err
		}

		if strings.Contains(contentType, "application/json") {
			return ctx.JSON(200, data)
		}
		return ctx.String(200, "OK")
	})

	// Test JSON
	jsonData := `{"name":"test","value":123}`
	req := httptest.NewRequest("POST", "/data", strings.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	if !strings.Contains(w.Header().Get("Content-Type"), "application/json") {
		t.Error("Expected JSON content type")
	}
}

// TestIntegration_ErrorHandling tests error propagation
func TestIntegration_ErrorHandling(t *testing.T) {
	router := New()

	var capturedError error
	router.SetErrorHandler(func(ctx Context, err error) {
		capturedError = err
		ctx.String(500, "Custom error: "+err.Error())
	})

	testErr := http.ErrNotSupported
	router.GET("/error", func(ctx Context) error {
		return testErr
	})

	req := httptest.NewRequest("GET", "/error", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if capturedError != testErr {
		t.Errorf("Error handler not called with correct error")
	}

	if !strings.Contains(w.Body.String(), "Custom error") {
		t.Error("Custom error handler not used")
	}
}
