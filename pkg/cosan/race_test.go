package cosan

import (
	"net/http/httptest"
	"sync"
	"testing"
)

// TestConcurrentRequests tests router under high concurrency
func TestConcurrentRequests(t *testing.T) {
	r := New()
	r.GET("/hello", func(ctx Context) error {
		return ctx.String(200, "Hello World")
	})

	const goroutines = 1000
	const requestsPerGoroutine = 100

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < requestsPerGoroutine; j++ {
				req := httptest.NewRequest("GET", "/hello", nil)
				w := httptest.NewRecorder()
				r.ServeHTTP(w, req)

				if w.Code != 200 {
					t.Errorf("Expected status 200, got %d", w.Code)
				}
			}
		}()
	}

	wg.Wait()
}

// TestConcurrentRouteRegistration tests concurrent route registration safety
func TestConcurrentRouteRegistration(t *testing.T) {
	r := New()

	const goroutines = 100

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		i := i
		go func() {
			defer wg.Done()
			// Use unique pattern to avoid conflicts
			pattern := "/route" + string(rune('a'+i%26)) + string(rune('a'+(i/26)%26))
			r.GET(pattern, func(ctx Context) error {
				return ctx.String(200, "OK")
			})
		}()
	}

	wg.Wait()

	// Verify routes were registered (may have duplicates due to pattern collision)
	rt := r.(*router)
	if len(rt.routes) == 0 {
		t.Error("Expected routes to be registered")
	}
}

// TestConcurrentMiddleware tests concurrent middleware execution
func TestConcurrentMiddleware(t *testing.T) {
	r := New()
	
	counter := 0
	var mu sync.Mutex
	
	r.Use(MiddlewareFunc(func(next HandlerFunc) HandlerFunc {
		return func(ctx Context) error {
			mu.Lock()
			counter++
			mu.Unlock()
			return next(ctx)
		}
	}))
	
	r.GET("/test", func(ctx Context) error {
		return ctx.String(200, "OK")
	})

	const requests = 1000
	var wg sync.WaitGroup
	wg.Add(requests)

	for i := 0; i < requests; i++ {
		go func() {
			defer wg.Done()
			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
		}()
	}

	wg.Wait()

	if counter != requests {
		t.Errorf("Expected counter to be %d, got %d", requests, counter)
	}
}

// TestConcurrentContextAccess tests concurrent context parameter access
func TestConcurrentContextAccess(t *testing.T) {
	r := New()
	r.GET("/users/:id", func(ctx Context) error {
		id := ctx.Param("id")
		return ctx.String(200, "User: "+id)
	})

	const goroutines = 500

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		i := i
		go func() {
			defer wg.Done()
			req := httptest.NewRequest("GET", "/users/"+string(rune('a'+i%26)), nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
		}()
	}

	wg.Wait()
}

// TestHighConcurrencyLoad tests router under extreme load
func TestHighConcurrencyLoad(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping high concurrency test in short mode")
	}

	r := New()
	
	// Add multiple routes
	r.GET("/", func(ctx Context) error { return ctx.String(200, "home") })
	r.GET("/about", func(ctx Context) error { return ctx.String(200, "about") })
	r.GET("/contact", func(ctx Context) error { return ctx.String(200, "contact") })
	r.GET("/users/:id", func(ctx Context) error { return ctx.String(200, ctx.Param("id")) })

	const goroutines = 10000
	const requestsPerGoroutine = 10

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()
			
			routes := []string{"/", "/about", "/contact", "/users/123"}
			for j := 0; j < requestsPerGoroutine; j++ {
				route := routes[j%len(routes)]
				req := httptest.NewRequest("GET", route, nil)
				w := httptest.NewRecorder()
				r.ServeHTTP(w, req)
			}
		}(i)
	}

	wg.Wait()
}
