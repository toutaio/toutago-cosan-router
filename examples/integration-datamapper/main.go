// Package main demonstrates integration with toutago-datamapper for parameter binding.
package main

import (
	"log"

	"github.com/toutaio/toutago-cosan-router/pkg/cosan"
)

// User represents a user with validation rules
type User struct {
	ID       int    `json:"id" validate:"required,min=1"`
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Age      int    `json:"age" validate:"min=18,max=120"`
}

// CreateUserRequest demonstrates complex request binding
type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Profile  struct {
		FirstName string `json:"first_name" validate:"required"`
		LastName  string `json:"last_name" validate:"required"`
		Age       int    `json:"age" validate:"min=18"`
	} `json:"profile"`
}

func main() {
	router := cosan.New()

	// Configure router with datamapper integration
	// Note: In production, you'd initialize the datamapper with proper configuration
	// router.SetBinder(datamapper.NewBinder())
	// router.SetValidator(datamapper.NewValidator())

	// Basic binding examples
	router.POST("/users", CreateUserHandler)
	router.PUT("/users/:id", UpdateUserHandler)

	// Complex binding with nested structures
	router.POST("/register", RegisterHandler)

	// Query parameter binding
	router.GET("/search", SearchUsersHandler)

	// Form data binding
	router.POST("/upload", UploadHandler)

	// Multiple source binding (path + body + query)
	router.PUT("/users/:id/profile", UpdateProfileHandler)

	log.Println("Server starting on http://localhost:8080")
	log.Println("Example requests:")
	log.Println("  POST /users -d '{\"username\":\"john\",\"email\":\"john@example.com\",\"age\":25}'")
	log.Println("  GET /search?username=john&min_age=18&max_age=65")
	log.Fatal(router.Listen(":8080"))
}

// CreateUserHandler demonstrates basic JSON binding with validation
func CreateUserHandler(ctx cosan.Context) error {
	var user User
	if err := ctx.Bind(&user); err != nil {
		return ctx.JSON(400, map[string]interface{}{
			"error": "Invalid request body",
			"details": err.Error(),
		})
	}

	// If using datamapper validator:
	// if err := ctx.Validate(&user); err != nil {
	//     return ctx.JSON(400, map[string]interface{}{
	//         "error": "Validation failed",
	//         "details": err,
	//     })
	// }

	// Simulate user creation
	user.ID = 1

	return ctx.JSON(201, map[string]interface{}{
		"message": "User created successfully",
		"user":    user,
	})
}

// UpdateUserHandler demonstrates binding with path parameters
func UpdateUserHandler(ctx cosan.Context) error {
	var user User
	if err := ctx.Bind(&user); err != nil {
		return ctx.JSON(400, map[string]interface{}{
			"error": "Invalid request body",
			"details": err.Error(),
		})
	}

	// Get ID from path parameter
	id := ctx.Param("id")

	return ctx.JSON(200, map[string]interface{}{
		"message": "User updated",
		"id":      id,
		"user":    user,
	})
}

// RegisterHandler demonstrates complex nested binding
func RegisterHandler(ctx cosan.Context) error {
	var req CreateUserRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(400, map[string]interface{}{
			"error": "Invalid registration data",
			"details": err.Error(),
		})
	}

	// Create user from request
	user := User{
		ID:       1,
		Username: req.Username,
		Email:    req.Email,
		Age:      req.Profile.Age,
	}

	return ctx.JSON(201, map[string]interface{}{
		"message": "Registration successful",
		"user":    user,
		"profile": req.Profile,
	})
}

// SearchUsersHandler demonstrates query parameter binding
func SearchUsersHandler(ctx cosan.Context) error {
	// Manual query parameter extraction
	username := ctx.Query("username")
	minAge := ctx.Query("min_age")
	maxAge := ctx.Query("max_age")
	page := ctx.Query("page")

	// With datamapper, you could bind to struct:
	// type SearchParams struct {
	//     Username string `query:"username"`
	//     MinAge   int    `query:"min_age"`
	//     MaxAge   int    `query:"max_age"`
	//     Page     int    `query:"page" default:"1"`
	// }
	// var params SearchParams
	// ctx.BindQuery(&params)

	return ctx.JSON(200, map[string]interface{}{
		"query": map[string]string{
			"username": username,
			"min_age":  minAge,
			"max_age":  maxAge,
			"page":     page,
		},
		"results": []User{
			{ID: 1, Username: "john", Email: "john@example.com", Age: 25},
			{ID: 2, Username: "jane", Email: "jane@example.com", Age: 30},
		},
	})
}

// UploadHandler demonstrates form data binding
func UploadHandler(ctx cosan.Context) error {
	// With datamapper, form binding would be:
	// type UploadForm struct {
	//     Title       string `form:"title" validate:"required"`
	//     Description string `form:"description"`
	//     File        []byte `form:"file" validate:"required"`
	// }
	// var form UploadForm
	// ctx.BindForm(&form)

	// Manual form handling for now
	if err := ctx.Request().ParseMultipartForm(10 << 20); err != nil {
		return ctx.JSON(400, map[string]string{
			"error": "Failed to parse form",
		})
	}

	title := ctx.Request().FormValue("title")
	description := ctx.Request().FormValue("description")

	return ctx.JSON(200, map[string]interface{}{
		"message": "File uploaded successfully",
		"title":   title,
		"description": description,
	})
}

// UpdateProfileHandler demonstrates binding from multiple sources
func UpdateProfileHandler(ctx cosan.Context) error {
	id := ctx.Param("id")

	type ProfileUpdate struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Age       int    `json:"age"`
	}

	var profile ProfileUpdate
	if err := ctx.Bind(&profile); err != nil {
		return ctx.JSON(400, map[string]interface{}{
			"error": "Invalid profile data",
			"details": err.Error(),
		})
	}

	// Get optional query parameters
	notify := ctx.Query("notify")

	return ctx.JSON(200, map[string]interface{}{
		"message": "Profile updated",
		"user_id": id,
		"profile": profile,
		"notify":  notify == "true",
	})
}
