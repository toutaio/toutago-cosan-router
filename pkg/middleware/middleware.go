package middleware

import (
	"fmt"
	"log"
	"runtime/debug"
	"time"

	"github.com/toutaio/toutago-cosan-router/pkg/cosan"
)

// Logger returns a middleware that logs HTTP requests.
// It logs the method, path, status code, and duration.
//
// Example:
//
// router.Use(middleware.Logger())
func Logger() cosan.Middleware {
	return cosan.MiddlewareFunc(func(next cosan.HandlerFunc) cosan.HandlerFunc {
		return func(ctx cosan.Context) error {
			start := time.Now()
			method := ctx.Request().Method
			path := ctx.Request().URL.Path

			// Call next handler
			err := next(ctx)

			// Log after response
			duration := time.Since(start)

			log.Printf("[%s] %s %s (%v)",
				method,
				path,
				statusFromError(err),
				duration,
			)

			return err
		}
	})
}

// Recovery returns a middleware that recovers from panics.
// It logs the panic and stack trace, then returns a 500 error.
//
// Example:
//
// router.Use(middleware.Recovery())
func Recovery() cosan.Middleware {
	return cosan.MiddlewareFunc(func(next cosan.HandlerFunc) cosan.HandlerFunc {
		return func(ctx cosan.Context) error {
			defer func() {
				if r := recover(); r != nil {
					// Log the panic and stack trace
					log.Printf("PANIC: %v\n%s", r, debug.Stack())

					// Return 500 error
					ctx.Status(500)
					ctx.Header().Set("Content-Type", "application/json")
					_, _ = ctx.Write([]byte(fmt.Sprintf(`{"error":"Internal Server Error","message":"%v"}`, r)))
				}
			}()

			return next(ctx)
		}
	})
}

// RequestID returns a middleware that adds a unique request ID.
// The ID is stored in the context and added to response headers.
//
// Example:
//
// router.Use(middleware.RequestID())
// // In handler: id := ctx.Get("requestID").(string)
func RequestID() cosan.Middleware {
	return cosan.MiddlewareFunc(func(next cosan.HandlerFunc) cosan.HandlerFunc {
		return func(ctx cosan.Context) error {
			// Check if request ID already exists in header
			requestID := ctx.Request().Header.Get("X-Request-ID")
			if requestID == "" {
				// Generate a new request ID
				requestID = fmt.Sprintf("%d", time.Now().UnixNano())
			}

			// Store in context
			ctx.Set("requestID", requestID)

			// Add to response headers
			ctx.Header().Set("X-Request-ID", requestID)

			return next(ctx)
		}
	})
}

// CORS returns a middleware that handles CORS headers.
//
// Example:
//
// router.Use(middleware.CORS())
func CORS() cosan.Middleware {
	return CORSWithConfig(CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	})
}

// CORSConfig holds CORS configuration.
type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	MaxAge           int
	AllowCredentials bool
}

// CORSWithConfig returns a CORS middleware with custom configuration.
func CORSWithConfig(config CORSConfig) cosan.Middleware {
	return cosan.MiddlewareFunc(func(next cosan.HandlerFunc) cosan.HandlerFunc {
		return func(ctx cosan.Context) error {
			origin := ctx.Request().Header.Get("Origin")

			// Set CORS headers
			if len(config.AllowOrigins) > 0 {
				if contains(config.AllowOrigins, "*") {
					ctx.Header().Set("Access-Control-Allow-Origin", "*")
				} else if contains(config.AllowOrigins, origin) {
					ctx.Header().Set("Access-Control-Allow-Origin", origin)
				}
			}

			if len(config.AllowMethods) > 0 {
				ctx.Header().Set("Access-Control-Allow-Methods", join(config.AllowMethods, ", "))
			}

			if len(config.AllowHeaders) > 0 {
				ctx.Header().Set("Access-Control-Allow-Headers", join(config.AllowHeaders, ", "))
			}

			if len(config.ExposeHeaders) > 0 {
				ctx.Header().Set("Access-Control-Expose-Headers", join(config.ExposeHeaders, ", "))
			}

			if config.AllowCredentials {
				ctx.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			if config.MaxAge > 0 {
				ctx.Header().Set("Access-Control-Max-Age", fmt.Sprintf("%d", config.MaxAge))
			}

			// Handle preflight request
			if ctx.Request().Method == "OPTIONS" {
				ctx.Status(204)
				return nil
			}

			return next(ctx)
		}
	})
}

// Helper functions

func statusFromError(err error) string {
	if err != nil {
		return "500 ERROR"
	}
	return "200 OK"
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func join(slice []string, sep string) string {
	if len(slice) == 0 {
		return ""
	}
	result := slice[0]
	for i := 1; i < len(slice); i++ {
		result += sep + slice[i]
	}
	return result
}
