// Package main demonstrates path parameter handling in Cosan router.
package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/toutaio/toutago-cosan-router"
)

type Product struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
}

func main() {
	router := cosan.New()

	// Simple path parameters
	router.GET("/users/:id", GetUserHandler)
	router.GET("/posts/:slug", GetPostHandler)

	// Multiple path parameters
	router.GET("/users/:userId/posts/:postId", GetUserPostHandler)
	router.GET("/categories/:category/products/:productId", GetProductHandler)

	// Wildcard parameters
	router.GET("/files/*filepath", GetFileHandler)
	router.GET("/static/*path", StaticFileHandler)

	// RESTful API with path parameters
	api := router.Group("/api/v1")
	api.GET("/products/:id", GetProductByIDHandler)
	api.PUT("/products/:id", UpdateProductHandler)
	api.DELETE("/products/:id", DeleteProductHandler)

	// Nested resources
	router.GET("/orgs/:org/repos/:repo/issues/:issue", GetIssueHandler)

	// Query parameters combined with path parameters
	router.GET("/search/:category", SearchHandler)

	log.Println("Server starting on http://localhost:8080")
	log.Fatal(router.Listen(":8080"))
}

// GetUserHandler demonstrates simple path parameter
func GetUserHandler(ctx cosan.Context) error {
	id := ctx.Param("id")

	// Validate and convert parameter
	userID, err := strconv.Atoi(id)
	if err != nil {
		return ctx.JSON(400, map[string]string{
			"error": "Invalid user ID - must be a number",
		})
	}

	return ctx.JSON(200, map[string]interface{}{
		"user_id": userID,
		"name":    fmt.Sprintf("User %d", userID),
		"email":   fmt.Sprintf("user%d@example.com", userID),
	})
}

// GetPostHandler demonstrates slug parameter
func GetPostHandler(ctx cosan.Context) error {
	slug := ctx.Param("slug")

	return ctx.JSON(200, map[string]interface{}{
		"slug":    slug,
		"title":   "Blog Post: " + slug,
		"content": "This is the content for " + slug,
	})
}

// GetUserPostHandler demonstrates multiple path parameters
func GetUserPostHandler(ctx cosan.Context) error {
	userID := ctx.Param("userId")
	postID := ctx.Param("postId")

	return ctx.JSON(200, map[string]interface{}{
		"user_id": userID,
		"post_id": postID,
		"message": fmt.Sprintf("Post %s by user %s", postID, userID),
	})
}

// GetProductHandler demonstrates nested path parameters
func GetProductHandler(ctx cosan.Context) error {
	category := ctx.Param("category")
	productID := ctx.Param("productId")

	return ctx.JSON(200, map[string]interface{}{
		"category":   category,
		"product_id": productID,
		"name":       fmt.Sprintf("Product %s in %s", productID, category),
	})
}

// GetFileHandler demonstrates wildcard parameter
func GetFileHandler(ctx cosan.Context) error {
	filepath := ctx.Param("filepath")

	return ctx.JSON(200, map[string]interface{}{
		"filepath": filepath,
		"message":  "File path captured: " + filepath,
	})
}

// StaticFileHandler demonstrates serving static files with wildcard
func StaticFileHandler(ctx cosan.Context) error {
	path := ctx.Param("path")

	// In a real app, you'd serve actual files here
	return ctx.JSON(200, map[string]interface{}{
		"static_path": path,
		"message":     "Static file: " + path,
	})
}

// GetProductByIDHandler demonstrates RESTful API with path parameter
func GetProductByIDHandler(ctx cosan.Context) error {
	id := ctx.Param("id")

	productID, err := strconv.Atoi(id)
	if err != nil {
		return ctx.JSON(400, map[string]string{
			"error": "Invalid product ID",
		})
	}

	product := Product{
		ID:       productID,
		Name:     fmt.Sprintf("Product %d", productID),
		Category: "Electronics",
	}

	return ctx.JSON(200, product)
}

// UpdateProductHandler demonstrates PUT with path parameter
func UpdateProductHandler(ctx cosan.Context) error {
	id := ctx.Param("id")

	var product Product
	if err := ctx.Bind(&product); err != nil {
		return err
	}

	product.ID, _ = strconv.Atoi(id)

	return ctx.JSON(200, map[string]interface{}{
		"message": "Product updated",
		"product": product,
	})
}

// DeleteProductHandler demonstrates DELETE with path parameter
func DeleteProductHandler(ctx cosan.Context) error {
	id := ctx.Param("id")

	return ctx.JSON(200, map[string]string{
		"message": fmt.Sprintf("Product %s deleted", id),
	})
}

// GetIssueHandler demonstrates deeply nested resources
func GetIssueHandler(ctx cosan.Context) error {
	org := ctx.Param("org")
	repo := ctx.Param("repo")
	issue := ctx.Param("issue")

	return ctx.JSON(200, map[string]interface{}{
		"organization": org,
		"repository":   repo,
		"issue_number": issue,
		"url":          fmt.Sprintf("/%s/%s/issues/%s", org, repo, issue),
	})
}

// SearchHandler demonstrates combining path and query parameters
func SearchHandler(ctx cosan.Context) error {
	category := ctx.Param("category")
	query := ctx.Query("q")
	page := ctx.Query("page")
	limit := ctx.Query("limit")

	// Set defaults
	if page == "" {
		page = "1"
	}
	if limit == "" {
		limit = "10"
	}

	return ctx.JSON(200, map[string]interface{}{
		"category": category,
		"query":    query,
		"page":     page,
		"limit":    limit,
		"results":  []string{"result1", "result2", "result3"},
	})
}
