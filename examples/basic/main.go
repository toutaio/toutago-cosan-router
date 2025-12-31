// Package main demonstrates basic usage of the Cosan router.
package main

import (
	"log"

	"github.com/toutaio/toutago-cosan-router/pkg/cosan"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	// Create a new router
	router := cosan.New()

	// Register routes
	router.GET("/", HomeHandler)
	router.GET("/users/:id", GetUser)
	router.POST("/users", CreateUser)

	// Route groups
	api := router.Group("/api/v1")
	api.GET("/status", StatusHandler)
	api.GET("/info", InfoHandler)

	// Start the server
	log.Println("Server starting on http://localhost:8080")
	log.Fatal(router.Listen(":8080"))
}

// HomeHandler handles the root path.
func HomeHandler(ctx cosan.Context) error {
	return ctx.JSON(200, map[string]string{
		"message": "Welcome to Cosan Router!",
	})
}

// GetUser demonstrates JSON responses.
func GetUser(ctx cosan.Context) error {
	id := ctx.Param("id")
	user := User{
		ID:   1,
		Name: "User " + id,
	}
	return ctx.JSON(200, user)
}

// CreateUser demonstrates JSON request binding.
func CreateUser(ctx cosan.Context) error {
	var user User
	if err := ctx.Bind(&user); err != nil {
		return err
	}
	// In a real app, you'd save to database here
	return ctx.JSON(201, user)
}

// StatusHandler handles status checks.
func StatusHandler(ctx cosan.Context) error {
	return ctx.JSON(200, map[string]string{
		"status": "ok",
	})
}

// InfoHandler handles info requests.
func InfoHandler(ctx cosan.Context) error {
	return ctx.JSON(200, map[string]interface{}{
		"name":    "Cosan Router",
		"version": "v0.1.0",
	})
}
