package benchmarks

import (
	"net/http/httptest"
	"testing"

	cosan "github.com/toutaio/toutago-cosan-router"
)

// BenchmarkRouterSimpleRoute benchmarks a simple route lookup
func BenchmarkRouterSimpleRoute(b *testing.B) {
	r := cosan.New()
	r.GET("/hello", func(ctx cosan.Context) error {
		return ctx.String(200, "Hello World")
	})

	req := httptest.NewRequest("GET", "/hello", nil)
	w := httptest.NewRecorder()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r.ServeHTTP(w, req)
	}
}

// BenchmarkRouterWithParams benchmarks route with parameters
func BenchmarkRouterWithParams(b *testing.B) {
	r := cosan.New()
	r.GET("/users/:id", func(ctx cosan.Context) error {
		id := ctx.Param("id")
		return ctx.String(200, "User: "+id)
	})

	req := httptest.NewRequest("GET", "/users/123", nil)
	w := httptest.NewRecorder()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r.ServeHTTP(w, req)
	}
}

// BenchmarkRouterWithMiddleware benchmarks route with middleware
func BenchmarkRouterWithMiddleware(b *testing.B) {
	r := cosan.New()

	// Add middleware
	r.Use(cosan.MiddlewareFunc(func(next cosan.HandlerFunc) cosan.HandlerFunc {
		return func(ctx cosan.Context) error {
			return next(ctx)
		}
	}))
	r.Use(cosan.MiddlewareFunc(func(next cosan.HandlerFunc) cosan.HandlerFunc {
		return func(ctx cosan.Context) error {
			return next(ctx)
		}
	}))

	r.GET("/api/users", func(ctx cosan.Context) error {
		return ctx.JSON(200, map[string]string{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/api/users", nil)
	w := httptest.NewRecorder()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r.ServeHTTP(w, req)
	}
}

// BenchmarkRouterStaticRoutes benchmarks multiple static routes
func BenchmarkRouterStaticRoutes(b *testing.B) {
	r := cosan.New()

	routes := []string{
		"/",
		"/about",
		"/contact",
		"/products",
		"/services",
		"/team",
		"/careers",
		"/blog",
		"/faq",
		"/privacy",
	}

	for _, route := range routes {
		r.GET(route, func(ctx cosan.Context) error {
			return ctx.String(200, "OK")
		})
	}

	req := httptest.NewRequest("GET", "/blog", nil)
	w := httptest.NewRecorder()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r.ServeHTTP(w, req)
	}
}

// BenchmarkRouterComplexRouting benchmarks complex route patterns
func BenchmarkRouterComplexRouting(b *testing.B) {
	r := cosan.New()

	// API routes with parameters
	r.GET("/api/v1/users/:id", func(ctx cosan.Context) error {
		return ctx.JSON(200, map[string]string{"id": ctx.Param("id")})
	})
	r.GET("/api/v1/users/:id/posts/:postId", func(ctx cosan.Context) error {
		return ctx.JSON(200, map[string]string{
			"userId": ctx.Param("id"),
			"postId": ctx.Param("postId"),
		})
	})
	r.POST("/api/v1/users", func(ctx cosan.Context) error {
		return ctx.JSON(201, map[string]string{"status": "created"})
	})

	req := httptest.NewRequest("GET", "/api/v1/users/42/posts/99", nil)
	w := httptest.NewRecorder()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r.ServeHTTP(w, req)
	}
}

// BenchmarkRouterParallel benchmarks concurrent requests
func BenchmarkRouterParallel(b *testing.B) {
	r := cosan.New()
	r.GET("/hello", func(ctx cosan.Context) error {
		return ctx.String(200, "Hello World")
	})

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		req := httptest.NewRequest("GET", "/hello", nil)
		w := httptest.NewRecorder()

		for pb.Next() {
			r.ServeHTTP(w, req)
		}
	})
}

// BenchmarkMemoryAllocations specifically tracks allocations
func BenchmarkMemoryAllocations(b *testing.B) {
	r := cosan.New()
	r.GET("/test/:id", func(ctx cosan.Context) error {
		_ = ctx.Param("id")
		return ctx.String(200, "OK")
	})

	req := httptest.NewRequest("GET", "/test/123", nil)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}
