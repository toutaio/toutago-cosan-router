package middleware_test

import (
"encoding/json"
"fmt"
"net/http"
"net/http/httptest"
"strings"
"testing"

"github.com/toutaio/toutago-cosan-router/pkg/cosan"
"github.com/toutaio/toutago-cosan-router/pkg/middleware"
)

func TestLogger(t *testing.T) {
router := cosan.New()
router.Use(middleware.Logger())
router.GET("/test", func(ctx cosan.Context) error {
return ctx.String(200, "OK")
})

req := httptest.NewRequest(http.MethodGet, "/test", nil)
w := httptest.NewRecorder()
router.ServeHTTP(w, req)

if w.Code != 200 {
t.Errorf("Expected status 200, got %d", w.Code)
}
}

func TestRecovery(t *testing.T) {
router := cosan.New()
router.Use(middleware.Recovery())
router.GET("/panic", func(ctx cosan.Context) error {
panic("test panic")
})

req := httptest.NewRequest(http.MethodGet, "/panic", nil)
w := httptest.NewRecorder()
router.ServeHTTP(w, req)

if w.Code != 500 {
t.Errorf("Expected status 500, got %d", w.Code)
}

var response map[string]string
json.Unmarshal(w.Body.Bytes(), &response)
if !strings.Contains(response["message"], "test panic") {
t.Errorf("Expected panic message in response, got %v", response)
}
}

func TestRequestID(t *testing.T) {
router := cosan.New()
router.Use(middleware.RequestID())
router.GET("/test", func(ctx cosan.Context) error {
if ctx.Get("requestID") == nil {
t.Error("Request ID not set in context")
}
return ctx.String(200, "OK")
})

req := httptest.NewRequest(http.MethodGet, "/test", nil)
w := httptest.NewRecorder()
router.ServeHTTP(w, req)

if w.Header().Get("X-Request-ID") == "" {
t.Error("Expected X-Request-ID header to be set")
}
}

func TestRequestIDPreserve(t *testing.T) {
router := cosan.New()
router.Use(middleware.RequestID())
router.GET("/test", func(ctx cosan.Context) error {
return ctx.String(200, "OK")
})

req := httptest.NewRequest(http.MethodGet, "/test", nil)
req.Header.Set("X-Request-ID", "existing-id-12345")
w := httptest.NewRecorder()
router.ServeHTTP(w, req)

if w.Header().Get("X-Request-ID") != "existing-id-12345" {
t.Errorf("Expected request ID 'existing-id-12345', got %s", w.Header().Get("X-Request-ID"))
}
}

func TestCORS(t *testing.T) {
router := cosan.New()
router.Use(middleware.CORS())
router.GET("/test", func(ctx cosan.Context) error {
return ctx.String(200, "OK")
})

req := httptest.NewRequest(http.MethodGet, "/test", nil)
req.Header.Set("Origin", "http://example.com")
w := httptest.NewRecorder()
router.ServeHTTP(w, req)

if w.Header().Get("Access-Control-Allow-Origin") != "*" {
t.Error("Expected CORS headers to be set")
}
}

func TestCORSPreflight(t *testing.T) {
router := cosan.New()
router.Use(middleware.CORS())
router.OPTIONS("/test", func(ctx cosan.Context) error {
return nil
})

req := httptest.NewRequest(http.MethodOptions, "/test", nil)
req.Header.Set("Origin", "http://example.com")
w := httptest.NewRecorder()
router.ServeHTTP(w, req)

if w.Code != 204 {
t.Errorf("Expected status 204 for OPTIONS, got %d", w.Code)
}
}

func TestMiddlewareChain(t *testing.T) {
router := cosan.New()
var order []string

mw1 := cosan.MiddlewareFunc(func(next cosan.HandlerFunc) cosan.HandlerFunc {
return func(ctx cosan.Context) error {
order = append(order, "mw1-before")
err := next(ctx)
order = append(order, "mw1-after")
return err
}
})

mw2 := cosan.MiddlewareFunc(func(next cosan.HandlerFunc) cosan.HandlerFunc {
return func(ctx cosan.Context) error {
order = append(order, "mw2-before")
err := next(ctx)
order = append(order, "mw2-after")
return err
}
})

router.Use(mw1, mw2)
router.GET("/test", func(ctx cosan.Context) error {
order = append(order, "handler")
return ctx.String(200, "OK")
})

req := httptest.NewRequest(http.MethodGet, "/test", nil)
w := httptest.NewRecorder()
router.ServeHTTP(w, req)

expected := []string{"mw1-before", "mw2-before", "handler", "mw2-after", "mw1-after"}
if len(order) != len(expected) {
t.Fatalf("Expected %d execution steps, got %d", len(expected), len(order))
}

for i, step := range expected {
if order[i] != step {
t.Errorf("Expected step %d to be %s, got %s", i, step, order[i])
}
}
}

func TestMiddlewareContextValues(t *testing.T) {
router := cosan.New()
authMW := cosan.MiddlewareFunc(func(next cosan.HandlerFunc) cosan.HandlerFunc {
return func(ctx cosan.Context) error {
ctx.Set("user", "alice")
ctx.Set("role", "admin")
return next(ctx)
}
})

router.Use(authMW)
router.GET("/test", func(ctx cosan.Context) error {
user := ctx.Get("user")
role := ctx.Get("role")
if user != "alice" || role != "admin" {
t.Errorf("Expected user=alice, role=admin, got user=%v, role=%v", user, role)
}
return ctx.JSON(200, map[string]interface{}{"user": user, "role": role})
})

req := httptest.NewRequest(http.MethodGet, "/test", nil)
w := httptest.NewRecorder()
router.ServeHTTP(w, req)

if w.Code != 200 {
t.Errorf("Expected status 200, got %d", w.Code)
}
}

func TestMiddlewareErrorHandling(t *testing.T) {
router := cosan.New()
errorMW := cosan.MiddlewareFunc(func(next cosan.HandlerFunc) cosan.HandlerFunc {
return func(ctx cosan.Context) error {
err := next(ctx)
if err != nil {
return ctx.JSON(400, map[string]string{"error": err.Error()})
}
return nil
}
})

router.Use(errorMW)
router.GET("/error", func(ctx cosan.Context) error {
return fmt.Errorf("custom error")
})

req := httptest.NewRequest(http.MethodGet, "/error", nil)
w := httptest.NewRecorder()
router.ServeHTTP(w, req)

var response map[string]string
json.Unmarshal(w.Body.Bytes(), &response)
if response["error"] != "custom error" {
t.Errorf("Expected error message, got %v", response)
}
}

func TestAllStandardMiddleware(t *testing.T) {
router := cosan.New()
router.Use(
middleware.Recovery(),
middleware.Logger(),
middleware.RequestID(),
middleware.CORS(),
)

router.GET("/test", func(ctx cosan.Context) error {
return ctx.JSON(200, map[string]string{"status": "ok"})
})

req := httptest.NewRequest(http.MethodGet, "/test", nil)
w := httptest.NewRecorder()
router.ServeHTTP(w, req)

if w.Code != 200 {
t.Errorf("Expected status 200, got %d", w.Code)
}
if w.Header().Get("X-Request-ID") == "" {
t.Error("Request ID middleware not applied")
}
if w.Header().Get("Access-Control-Allow-Origin") == "" {
t.Error("CORS middleware not applied")
}
}

func BenchmarkMiddlewareChain(b *testing.B) {
router := cosan.New()
router.Use(middleware.Recovery(), middleware.Logger(), middleware.RequestID())
router.GET("/bench", func(ctx cosan.Context) error {
return ctx.String(200, "OK")
})

req := httptest.NewRequest(http.MethodGet, "/bench", nil)
b.ResetTimer()
b.ReportAllocs()

for i := 0; i < b.N; i++ {
w := httptest.NewRecorder()
router.ServeHTTP(w, req)
}
}
