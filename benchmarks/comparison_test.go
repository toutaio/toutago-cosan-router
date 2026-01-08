package benchmarks

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/toutaio/toutago-cosan-router"
)

// Benchmark comparison suite against stdlib
// To compare with other routers (Chi, Gin, Echo), uncomment and add dependencies

// BenchmarkCosanRouter_SingleRoute benchmarks Cosan router with single route
func BenchmarkCosanRouter_SingleRoute(b *testing.B) {
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

// BenchmarkStdlib_SingleRoute benchmarks stdlib ServeMux with single route
func BenchmarkStdlib_SingleRoute(b *testing.B) {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("Hello World"))
	})

	req := httptest.NewRequest("GET", "/hello", nil)
	w := httptest.NewRecorder()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mux.ServeHTTP(w, req)
	}
}

// BenchmarkCosanRouter_MultiRoute benchmarks Cosan with multiple routes
func BenchmarkCosanRouter_MultiRoute(b *testing.B) {
	r := cosan.New()

	routes := []string{"/a", "/b", "/c", "/d", "/e", "/f", "/g", "/h", "/i", "/j"}
	for _, route := range routes {
		r.GET(route, func(ctx cosan.Context) error {
			return ctx.String(200, "OK")
		})
	}

	req := httptest.NewRequest("GET", "/j", nil)
	w := httptest.NewRecorder()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r.ServeHTTP(w, req)
	}
}

// BenchmarkStdlib_MultiRoute benchmarks stdlib with multiple routes
func BenchmarkStdlib_MultiRoute(b *testing.B) {
	mux := http.NewServeMux()

	routes := []string{"/a", "/b", "/c", "/d", "/e", "/f", "/g", "/h", "/i", "/j"}
	for _, route := range routes {
		mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			w.Write([]byte("OK"))
		})
	}

	req := httptest.NewRequest("GET", "/j", nil)
	w := httptest.NewRecorder()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mux.ServeHTTP(w, req)
	}
}

/*
// Uncomment to benchmark against Chi router
// Requires: go get -u github.com/go-chi/chi/v5

import "github.com/go-chi/chi/v5"

func BenchmarkChi_SingleRoute(b *testing.B) {
	r := chi.NewRouter()
	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("Hello World"))
	})

	req := httptest.NewRequest("GET", "/hello", nil)
	w := httptest.NewRecorder()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r.ServeHTTP(w, req)
	}
}
*/

/*
// Uncomment to benchmark against Gin
// Requires: go get -u github.com/gin-gonic/gin

import "github.com/gin-gonic/gin"

func BenchmarkGin_SingleRoute(b *testing.B) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.GET("/hello", func(c *gin.Context) {
		c.String(200, "Hello World")
	})

	req := httptest.NewRequest("GET", "/hello", nil)
	w := httptest.NewRecorder()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r.ServeHTTP(w, req)
	}
}
*/
