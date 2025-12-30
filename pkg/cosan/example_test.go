package cosan_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/toutaio/toutago-cosan-router/pkg/cosan"
)

// ExampleNew demonstrates basic router creation and usage
func ExampleNew() {
	router := cosan.New()

	router.GET("/hello", func(ctx cosan.Context) error {
		return ctx.String(200, "Hello, World!")
	})

	// Simulate a request
	req := httptest.NewRequest("GET", "/hello", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	fmt.Println(w.Body.String())
	// Output: Hello, World!
}

// ExampleRouter_GET demonstrates GET route registration
func ExampleRouter_GET() {
	router := cosan.New()

	router.GET("/users/:id", func(ctx cosan.Context) error {
		id := ctx.Param("id")
		return ctx.String(200, fmt.Sprintf("User ID: %s", id))
	})

	req := httptest.NewRequest("GET", "/users/123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	fmt.Println(w.Body.String())
	// Output: User ID: 123
}

// ExampleRouter_POST demonstrates POST route registration
func ExampleRouter_POST() {
	router := cosan.New()

	router.POST("/users", func(ctx cosan.Context) error {
		return ctx.String(201, "User created")
	})

	req := httptest.NewRequest("POST", "/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	fmt.Println(w.Code, w.Body.String())
	// Output: 201 User created
}

// ExampleRouter_Use demonstrates middleware usage
func ExampleRouter_Use() {
	router := cosan.New()

	// Add logging middleware
	logger := cosan.MiddlewareFunc(func(next cosan.HandlerFunc) cosan.HandlerFunc {
		return func(ctx cosan.Context) error {
			fmt.Println("Request received")
			return next(ctx)
		}
	})

	router.Use(logger)

	router.GET("/test", func(ctx cosan.Context) error {
		return ctx.String(200, "OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	fmt.Println(w.Body.String())
	// Output:
	// Request received
	// OK
}

// ExampleRouter_Group demonstrates route grouping
func ExampleRouter_Group() {
	router := cosan.New()

	api := router.Group("/api/v1")
	api.GET("/users", func(ctx cosan.Context) error {
		return ctx.String(200, "Users list")
	})

	req := httptest.NewRequest("GET", "/api/v1/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	fmt.Println(w.Body.String())
	// Output: Users list
}

// ExampleContext_Param demonstrates path parameter extraction
func ExampleContext_Param() {
	router := cosan.New()

	router.GET("/posts/:slug", func(ctx cosan.Context) error {
		slug := ctx.Param("slug")
		return ctx.String(200, slug)
	})

	req := httptest.NewRequest("GET", "/posts/hello-world", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	fmt.Println(w.Body.String())
	// Output: hello-world
}

// ExampleContext_Query demonstrates query parameter access
func ExampleContext_Query() {
	router := cosan.New()

	router.GET("/search", func(ctx cosan.Context) error {
		q := ctx.Query("q")
		return ctx.String(200, fmt.Sprintf("Search: %s", q))
	})

	req := httptest.NewRequest("GET", "/search?q=golang", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	fmt.Println(w.Body.String())
	// Output: Search: golang
}

// ExampleContext_JSON demonstrates JSON response
func ExampleContext_JSON() {
	router := cosan.New()

	router.GET("/user", func(ctx cosan.Context) error {
		return ctx.JSON(200, map[string]string{
			"name": "John",
			"role": "Admin",
		})
	})

	req := httptest.NewRequest("GET", "/user", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	fmt.Println(w.Body.String())
	// Output: {"name":"John","role":"Admin"}
}

// ExampleRouter_BeforeRequest demonstrates before-request hooks
func ExampleRouter_BeforeRequest() {
	router := cosan.New()

	router.BeforeRequest(func(req *http.Request) error {
		fmt.Println("Before:", req.Method, req.URL.Path)
		return nil
	})

	router.GET("/test", func(ctx cosan.Context) error {
		return ctx.String(200, "OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Output: Before: GET /test
}

// ExampleRouter_AfterResponse demonstrates after-response hooks
func ExampleRouter_AfterResponse() {
	router := cosan.New()

	router.AfterResponse(func(req *http.Request, statusCode int) {
		fmt.Printf("After: %d\n", statusCode)
	})

	router.GET("/test", func(ctx cosan.Context) error {
		return ctx.String(200, "OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Output: After: 200
}
