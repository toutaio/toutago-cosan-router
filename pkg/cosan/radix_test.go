package cosan

import (
"net/http"
"net/http/httptest"
"testing"
)

// TestPathParameters tests basic path parameter extraction.
func TestPathParameters(t *testing.T) {
router := New()

router.GET("/users/:id", func(ctx Context) error {
id := ctx.Param("id")
return ctx.JSON(200, map[string]string{"id": id})
})

tests := []struct {
path       string
expectID   string
expectCode int
}{
{"/users/123", "123", 200},
{"/users/abc", "abc", 200},
{"/users/test-user", "test-user", 200},
{"/users/", "", 404}, // No ID provided
{"/users", "", 404},  // No trailing slash
}

for _, tt := range tests {
req := httptest.NewRequest(http.MethodGet, tt.path, nil)
w := httptest.NewRecorder()
router.ServeHTTP(w, req)

if w.Code != tt.expectCode {
t.Errorf("Path %s: expected status %d, got %d", tt.path, tt.expectCode, w.Code)
}
}
}

// TestMultipleParameters tests routes with multiple parameters.
func TestMultipleParameters(t *testing.T) {
router := New()

router.GET("/users/:userID/posts/:postID", func(ctx Context) error {
userID := ctx.Param("userID")
postID := ctx.Param("postID")
return ctx.JSON(200, map[string]string{
"userID": userID,
"postID": postID,
})
})

req := httptest.NewRequest(http.MethodGet, "/users/123/posts/456", nil)
w := httptest.NewRecorder()
router.ServeHTTP(w, req)

if w.Code != 200 {
t.Errorf("Expected status 200, got %d", w.Code)
}
}

// TestWildcardParameters tests wildcard/catch-all routes.
func TestWildcardParameters(t *testing.T) {
router := New()

router.GET("/files/*filepath", func(ctx Context) error {
filepath := ctx.Param("filepath")
return ctx.JSON(200, map[string]string{"filepath": filepath})
})

tests := []struct {
path         string
expectPath   string
expectCode   int
}{
{"/files/docs/readme.md", "docs/readme.md", 200},
{"/files/a/b/c/d.txt", "a/b/c/d.txt", 200},
{"/files/single.txt", "single.txt", 200},
}

for _, tt := range tests {
req := httptest.NewRequest(http.MethodGet, tt.path, nil)
w := httptest.NewRecorder()
router.ServeHTTP(w, req)

if w.Code != tt.expectCode {
t.Errorf("Path %s: expected status %d, got %d", tt.path, tt.expectCode, w.Code)
}
}
}

// TestStaticVsParamPriority tests that static routes have priority over params.
func TestStaticVsParamPriority(t *testing.T) {
router := New()

staticCalled := false
paramCalled := false

router.GET("/users/me", func(ctx Context) error {
staticCalled = true
return ctx.String(200, "static")
})

router.GET("/users/:id", func(ctx Context) error {
paramCalled = true
return ctx.String(200, "param")
})

// Test static route
req := httptest.NewRequest(http.MethodGet, "/users/me", nil)
w := httptest.NewRecorder()
router.ServeHTTP(w, req)

if !staticCalled {
t.Error("Expected static route to be called")
}
if paramCalled {
t.Error("Did not expect param route to be called for /users/me")
}

// Reset
staticCalled = false
paramCalled = false

// Test param route
req = httptest.NewRequest(http.MethodGet, "/users/123", nil)
w = httptest.NewRecorder()
router.ServeHTTP(w, req)

if staticCalled {
t.Error("Did not expect static route to be called for /users/123")
}
if !paramCalled {
t.Error("Expected param route to be called")
}
}

// TestNestedParameters tests parameters in nested route groups.
func TestNestedParameters(t *testing.T) {
router := New()

api := router.Group("/api/v1")
users := api.Group("/users")

users.GET("/:id", func(ctx Context) error {
id := ctx.Param("id")
return ctx.JSON(200, map[string]string{"id": id})
})

users.GET("/:id/posts/:postID", func(ctx Context) error {
userID := ctx.Param("id")
postID := ctx.Param("postID")
return ctx.JSON(200, map[string]string{
"userID": userID,
"postID": postID,
})
})

// Test first route
req := httptest.NewRequest(http.MethodGet, "/api/v1/users/123", nil)
w := httptest.NewRecorder()
router.ServeHTTP(w, req)

if w.Code != 200 {
t.Errorf("Expected status 200, got %d", w.Code)
}

// Test nested route
req = httptest.NewRequest(http.MethodGet, "/api/v1/users/123/posts/456", nil)
w = httptest.NewRecorder()
router.ServeHTTP(w, req)

if w.Code != 200 {
t.Errorf("Expected status 200, got %d", w.Code)
}
}

// TestComplexRoutePatterns tests complex route combinations.
func TestComplexRoutePatterns(t *testing.T) {
router := New()

// Various route patterns
router.GET("/", func(ctx Context) error {
return ctx.String(200, "root")
})

router.GET("/users", func(ctx Context) error {
return ctx.String(200, "users-list")
})

router.GET("/users/:id", func(ctx Context) error {
return ctx.String(200, "user-detail")
})

router.GET("/users/:id/profile", func(ctx Context) error {
return ctx.String(200, "user-profile")
})

router.GET("/users/:id/posts/:postID", func(ctx Context) error {
return ctx.String(200, "user-post")
})

router.GET("/static/files/*path", func(ctx Context) error {
return ctx.String(200, "static-file")
})

tests := []struct {
path         string
expectedBody string
expectedCode int
}{
{"/", "root", 200},
{"/users", "users-list", 200},
{"/users/123", "user-detail", 200},
{"/users/123/profile", "user-profile", 200},
{"/users/123/posts/456", "user-post", 200},
{"/static/files/css/main.css", "static-file", 200},
{"/notfound", "", 404},
}

for _, tt := range tests {
req := httptest.NewRequest(http.MethodGet, tt.path, nil)
w := httptest.NewRecorder()
router.ServeHTTP(w, req)

if w.Code != tt.expectedCode {
t.Errorf("Path %s: expected status %d, got %d", tt.path, tt.expectedCode, w.Code)
}

if tt.expectedCode == 200 && w.Body.String() != tt.expectedBody {
t.Errorf("Path %s: expected body %q, got %q", tt.path, tt.expectedBody, w.Body.String())
}
}
}

// TestParamExtraction verifies parameter extraction works correctly.
func TestParamExtraction(t *testing.T) {
router := New()

var capturedParams map[string]string

router.GET("/users/:userID/posts/:postID/comments/:commentID", func(ctx Context) error {
// Make a copy since context will be pooled
params := ctx.Params()
capturedParams = make(map[string]string)
for k, v := range params {
capturedParams[k] = v
}
return ctx.String(200, "OK")
})

req := httptest.NewRequest(http.MethodGet, "/users/u123/posts/p456/comments/c789", nil)
w := httptest.NewRecorder()
router.ServeHTTP(w, req)

if w.Code != 200 {
t.Fatalf("Expected status 200, got %d", w.Code)
}

expected := map[string]string{
"userID":    "u123",
"postID":    "p456",
"commentID": "c789",
}

for key, expectedVal := range expected {
if val := capturedParams[key]; val != expectedVal {
t.Errorf("Expected param %s=%s, got %s", key, expectedVal, val)
}
}
}

// TestMethodSeparation verifies different methods don't interfere.
func TestMethodSeparation(t *testing.T) {
router := New()

router.GET("/users/:id", func(ctx Context) error {
return ctx.String(200, "GET")
})

router.POST("/users/:id", func(ctx Context) error {
return ctx.String(200, "POST")
})

router.DELETE("/users/:id", func(ctx Context) error {
return ctx.String(200, "DELETE")
})

tests := []struct {
method string
body   string
}{
{http.MethodGet, "GET"},
{http.MethodPost, "POST"},
{http.MethodDelete, "DELETE"},
}

for _, tt := range tests {
req := httptest.NewRequest(tt.method, "/users/123", nil)
w := httptest.NewRecorder()
router.ServeHTTP(w, req)

if w.Code != 200 {
t.Errorf("Method %s: expected status 200, got %d", tt.method, w.Code)
}

if w.Body.String() != tt.body {
t.Errorf("Method %s: expected body %q, got %q", tt.method, tt.body, w.Body.String())
}
}
}

// BenchmarkRadixStaticRoute benchmarks static route matching.
func BenchmarkRadixStaticRoute(b *testing.B) {
router := New()
router.GET("/users/profile", func(ctx Context) error {
return ctx.String(200, "OK")
})

req := httptest.NewRequest(http.MethodGet, "/users/profile", nil)
b.ResetTimer()
b.ReportAllocs()

for i := 0; i < b.N; i++ {
w := httptest.NewRecorder()
router.ServeHTTP(w, req)
}
}

// BenchmarkRadixParamRoute benchmarks parameterized route matching.
func BenchmarkRadixParamRoute(b *testing.B) {
router := New()
router.GET("/users/:id", func(ctx Context) error {
return ctx.String(200, "OK")
})

req := httptest.NewRequest(http.MethodGet, "/users/123", nil)
b.ResetTimer()
b.ReportAllocs()

for i := 0; i < b.N; i++ {
w := httptest.NewRecorder()
router.ServeHTTP(w, req)
}
}

// BenchmarkRadixComplexRoute benchmarks complex parameterized routes.
func BenchmarkRadixComplexRoute(b *testing.B) {
router := New()
router.GET("/api/v1/users/:userID/posts/:postID/comments/:commentID", func(ctx Context) error {
return ctx.String(200, "OK")
})

req := httptest.NewRequest(http.MethodGet, "/api/v1/users/123/posts/456/comments/789", nil)
b.ResetTimer()
b.ReportAllocs()

for i := 0; i < b.N; i++ {
w := httptest.NewRecorder()
router.ServeHTTP(w, req)
}
}
