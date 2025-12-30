// Package main demonstrates middleware usage with the Cosan router.
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/toutaio/toutago-cosan-router/pkg/cosan"
)

func main() {
	router := cosan.New()

	// Global middleware - applied to all routes
	router.Use(cosan.MiddlewareFunc(LoggerMiddleware))
	router.Use(cosan.MiddlewareFunc(RecoveryMiddleware))

	// Public routes (no auth required)
	router.GET("/", HomeHandler)
	router.GET("/public", PublicHandler)

	// Protected routes with auth middleware
	protected := router.Group("/api")
	protected.Use(cosan.MiddlewareFunc(AuthMiddleware))
	protected.GET("/profile", ProfileHandler)
	protected.POST("/data", DataHandler)

	// Admin routes with multiple middleware
	admin := router.Group("/admin")
	admin.Use(cosan.MiddlewareFunc(AuthMiddleware))
	admin.Use(cosan.MiddlewareFunc(AdminMiddleware))
	admin.GET("/dashboard", DashboardHandler)
	admin.DELETE("/users/:id", DeleteUserHandler)

	// Route-specific middleware - apply to handler
	slowHandler := cosan.MiddlewareFunc(TimeoutMiddleware(5*time.Second)).Process(SlowHandler)
	router.GET("/slow", slowHandler)

	log.Println("Server starting on http://localhost:8080")
	log.Fatal(router.Listen(":8080"))
}

// LoggerMiddleware logs request information
func LoggerMiddleware(next cosan.HandlerFunc) cosan.HandlerFunc {
	return func(ctx cosan.Context) error {
		start := time.Now()
		path := ctx.Request().URL.Path
		method := ctx.Request().Method

		err := next(ctx)

		duration := time.Since(start)
		log.Printf("[%s] %s - %v", method, path, duration)

		return err
	}
}

// RecoveryMiddleware recovers from panics
func RecoveryMiddleware(next cosan.HandlerFunc) cosan.HandlerFunc {
	return func(ctx cosan.Context) error {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic recovered: %v", r)
				ctx.JSON(500, map[string]string{
					"error": "Internal server error",
				})
			}
		}()
		return next(ctx)
	}
}

// AuthMiddleware checks for authentication
func AuthMiddleware(next cosan.HandlerFunc) cosan.HandlerFunc {
	return func(ctx cosan.Context) error {
		token := ctx.Request().Header.Get("Authorization")
		if token == "" {
			return ctx.JSON(401, map[string]string{
				"error": "Unauthorized - no token provided",
			})
		}

		// Simplified token validation
		if token != "Bearer valid-token" {
			return ctx.JSON(401, map[string]string{
				"error": "Unauthorized - invalid token",
			})
		}

		// Store user info in context for downstream handlers
		ctx.Set("user_id", "123")
		ctx.Set("username", "john.doe")

		return next(ctx)
	}
}

// AdminMiddleware checks for admin role
func AdminMiddleware(next cosan.HandlerFunc) cosan.HandlerFunc {
	return func(ctx cosan.Context) error {
		// In a real app, check user role from database
		userID := ctx.Get("user_id")
		if userID != "123" { // Simplified admin check
			return ctx.JSON(403, map[string]string{
				"error": "Forbidden - admin access required",
			})
		}

		return next(ctx)
	}
}

// TimeoutMiddleware adds timeout to handler
func TimeoutMiddleware(timeout time.Duration) cosan.MiddlewareFunc {
	return func(next cosan.HandlerFunc) cosan.HandlerFunc {
		return func(ctx cosan.Context) error {
			done := make(chan error, 1)

			go func() {
				done <- next(ctx)
			}()

			select {
			case err := <-done:
				return err
			case <-time.After(timeout):
				return ctx.JSON(408, map[string]string{
					"error": "Request timeout",
				})
			}
		}
	}
}

// Handlers

func HomeHandler(ctx cosan.Context) error {
	return ctx.JSON(200, map[string]string{
		"message": "Public home page",
	})
}

func PublicHandler(ctx cosan.Context) error {
	return ctx.JSON(200, map[string]string{
		"message": "Public endpoint - no auth required",
	})
}

func ProfileHandler(ctx cosan.Context) error {
	username := ctx.Get("username")
	return ctx.JSON(200, map[string]interface{}{
		"message":  "Protected profile endpoint",
		"username": username,
	})
}

func DataHandler(ctx cosan.Context) error {
	var data map[string]interface{}
	if err := ctx.Bind(&data); err != nil {
		return err
	}
	return ctx.JSON(200, map[string]interface{}{
		"message": "Data received",
		"data":    data,
	})
}

func DashboardHandler(ctx cosan.Context) error {
	return ctx.JSON(200, map[string]string{
		"message": "Admin dashboard - requires admin role",
	})
}

func DeleteUserHandler(ctx cosan.Context) error {
	id := ctx.Param("id")
	return ctx.JSON(200, map[string]string{
		"message": fmt.Sprintf("User %s deleted (admin only)", id),
	})
}

func SlowHandler(ctx cosan.Context) error {
	// Simulate slow operation
	time.Sleep(3 * time.Second)
	return ctx.JSON(200, map[string]string{
		"message": "Slow operation completed",
	})
}
