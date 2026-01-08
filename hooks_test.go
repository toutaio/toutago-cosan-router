package cosan

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouterHooks_BeforeRequest(t *testing.T) {
	r := New()
	called := false

	r.BeforeRequest(func(req *http.Request) error {
		called = true
		return nil
	})

	r.GET("/test", func(ctx Context) error {
		return ctx.String(200, "OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if !called {
		t.Error("BeforeRequest hook was not called")
	}
}

func TestRouterHooks_BeforeRequestError(t *testing.T) {
	r := New()
	hookErr := errors.New("hook error")

	r.BeforeRequest(func(req *http.Request) error {
		return hookErr
	})

	r.GET("/test", func(ctx Context) error {
		return ctx.String(200, "OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 500 {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestRouterHooks_MultipleBeforeHooks(t *testing.T) {
	r := New()
	var order []int

	r.BeforeRequest(func(req *http.Request) error {
		order = append(order, 1)
		return nil
	})
	r.BeforeRequest(func(req *http.Request) error {
		order = append(order, 2)
		return nil
	})

	r.GET("/test", func(ctx Context) error {
		return ctx.String(200, "OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if len(order) != 2 || order[0] != 1 || order[1] != 2 {
		t.Errorf("Hooks executed in wrong order: %v", order)
	}
}

func TestRouterHooks_AfterResponse(t *testing.T) {
	r := New()
	called := false
	var capturedStatus int

	r.AfterResponse(func(req *http.Request, statusCode int) {
		called = true
		capturedStatus = statusCode
	})

	r.GET("/test", func(ctx Context) error {
		return ctx.String(201, "Created")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if !called {
		t.Error("AfterResponse hook was not called")
	}
	if capturedStatus != 201 {
		t.Errorf("Expected status 201, got %d", capturedStatus)
	}
}

func TestRouterHooks_MultipleAfterHooks(t *testing.T) {
	r := New()
	var order []int

	r.AfterResponse(func(req *http.Request, statusCode int) {
		order = append(order, 1)
	})
	r.AfterResponse(func(req *http.Request, statusCode int) {
		order = append(order, 2)
	})

	r.GET("/test", func(ctx Context) error {
		return ctx.String(200, "OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if len(order) != 2 || order[0] != 1 || order[1] != 2 {
		t.Errorf("Hooks executed in wrong order: %v", order)
	}
}

func TestRouterHooks_CustomErrorHandler(t *testing.T) {
	r := New()
	called := false
	var capturedErr error

	r.SetErrorHandler(func(ctx Context, err error) {
		called = true
		capturedErr = err
		ctx.String(418, "Custom Error: "+err.Error())
	})

	testErr := errors.New("test error")
	r.GET("/test", func(ctx Context) error {
		return testErr
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if !called {
		t.Error("Custom error handler was not called")
	}
	if capturedErr != testErr {
		t.Errorf("Expected error %v, got %v", testErr, capturedErr)
	}
	if w.Code != 418 {
		t.Errorf("Expected status 418, got %d", w.Code)
	}
}

func TestRouterHooks_DefaultErrorHandler(t *testing.T) {
	r := New()
	testErr := errors.New("test error")

	r.GET("/test", func(ctx Context) error {
		return testErr
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 500 {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
	if w.Body.String() != "Internal Server Error: test error" {
		t.Errorf("Unexpected error message: %s", w.Body.String())
	}
}

func TestRouterHooks_NoHooks(t *testing.T) {
	r := New()

	r.GET("/test", func(ctx Context) error {
		return ctx.String(200, "OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}
