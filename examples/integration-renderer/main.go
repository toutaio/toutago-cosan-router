// Package main demonstrates integration with toutago-fith-renderer for HTML rendering.
package main

import (
	"log"

	cosan "github.com/toutaio/toutago-cosan-router"
)

type PageData struct {
	Title   string
	Content string
	User    User
}

type User struct {
	ID       int
	Username string
	Email    string
}

type Product struct {
	ID    int
	Name  string
	Price float64
}

func main() {
	router := cosan.New()

	// Configure router with renderer integration
	// Note: In production, you'd initialize the renderer with template directory
	// renderer := fith.New(fith.Config{
	//     TemplateDir: "templates",
	//     Extension:   ".html",
	//     Cache:       true,
	// })
	// router.SetRenderer(renderer)

	// HTML page rendering
	router.GET("/", HomePageHandler)
	router.GET("/about", AboutPageHandler)

	// Dynamic content rendering
	router.GET("/users/:id", UserProfileHandler)
	router.GET("/products", ProductListHandler)

	// Forms
	router.GET("/login", LoginFormHandler)
	router.POST("/login", LoginHandler)
	router.GET("/register", RegisterFormHandler)

	// Partial rendering (HTMX/AJAX support)
	router.GET("/partial/user/:id", UserPartialHandler)

	// API endpoints returning JSON
	api := router.Group("/api")
	api.GET("/users/:id", GetUserAPIHandler)
	api.GET("/products", GetProductsAPIHandler)

	log.Println("Server starting on http://localhost:8080")
	log.Println("Visit http://localhost:8080 in your browser")
	log.Fatal(router.Listen(":8080"))
}

// HomePageHandler renders the home page
func HomePageHandler(ctx cosan.Context) error {
	// With fith renderer:
	// return ctx.Render(200, "home", PageData{
	//     Title:   "Welcome",
	//     Content: "Welcome to Cosan Router with Fith Renderer",
	// })

	// Without renderer, return HTML directly
	html := `<!DOCTYPE html>
<html>
<head><title>Home</title></head>
<body>
	<h1>Welcome to Cosan Router</h1>
	<p>This is the home page.</p>
	<nav>
		<a href="/about">About</a> |
		<a href="/products">Products</a> |
		<a href="/login">Login</a>
	</nav>
</body>
</html>`
	return ctx.HTML(200, html)
}

// AboutPageHandler renders the about page
func AboutPageHandler(ctx cosan.Context) error {
	// With fith renderer:
	// return ctx.Render(200, "about", map[string]interface{}{
	//     "title": "About Us",
	//     "description": "Learn more about our application",
	// })

	html := `<!DOCTYPE html>
<html>
<head><title>About</title></head>
<body>
	<h1>About Us</h1>
	<p>This example demonstrates HTML rendering with Cosan Router.</p>
	<a href="/">Home</a>
</body>
</html>`
	return ctx.HTML(200, html)
}

// UserProfileHandler renders a user profile page
func UserProfileHandler(ctx cosan.Context) error {
	id := ctx.Param("id")

	user := User{
		ID:       1,
		Username: "john_doe",
		Email:    "john@example.com",
	}

	// With fith renderer:
	// return ctx.Render(200, "user/profile", map[string]interface{}{
	//     "title": "User Profile",
	//     "user":  user,
	// })

	html := `<!DOCTYPE html>
<html>
<head><title>User Profile</title></head>
<body>
	<h1>User Profile: ` + id + `</h1>
	<p>Username: ` + user.Username + `</p>
	<p>Email: ` + user.Email + `</p>
	<a href="/">Home</a>
</body>
</html>`
	return ctx.HTML(200, html)
}

// ProductListHandler renders a list of products
func ProductListHandler(ctx cosan.Context) error {
	// With fith renderer using loops:
	// return ctx.Render(200, "products/list", map[string]interface{}{
	//     "title":    "Products",
	//     "products": products,
	// })

	html := `<!DOCTYPE html>
<html>
<head><title>Products</title></head>
<body>
	<h1>Products</h1>
	<ul>
		<li>Laptop - $999.99</li>
		<li>Mouse - $29.99</li>
		<li>Keyboard - $79.99</li>
	</ul>
	<a href="/">Home</a>
</body>
</html>`
	return ctx.HTML(200, html)
}

// LoginFormHandler renders the login form
func LoginFormHandler(ctx cosan.Context) error {
	// With fith renderer:
	// return ctx.Render(200, "auth/login", map[string]interface{}{
	//     "title": "Login",
	//     "csrf":  generateCSRFToken(),
	// })

	html := `<!DOCTYPE html>
<html>
<head><title>Login</title></head>
<body>
	<h1>Login</h1>
	<form method="POST" action="/login">
		<div>
			<label>Username:</label>
			<input type="text" name="username" required>
		</div>
		<div>
			<label>Password:</label>
			<input type="password" name="password" required>
		</div>
		<button type="submit">Login</button>
	</form>
	<p><a href="/register">Register</a> | <a href="/">Home</a></p>
</body>
</html>`
	return ctx.HTML(200, html)
}

// LoginHandler processes login form
func LoginHandler(ctx cosan.Context) error {
	if err := ctx.Request().ParseForm(); err != nil {
		return ctx.JSON(400, map[string]string{"error": "Invalid form"})
	}

	username := ctx.Request().FormValue("username")
	password := ctx.Request().FormValue("password")

	// Validate credentials (simplified)
	if username == "admin" && password == "admin123" {
		// With fith renderer, redirect or render success:
		// return ctx.Render(200, "auth/success", map[string]interface{}{
		//     "username": username,
		// })

		html := `<!DOCTYPE html>
<html>
<head><title>Login Success</title></head>
<body>
	<h1>Welcome, ` + username + `!</h1>
	<p>You are now logged in.</p>
	<a href="/">Home</a>
</body>
</html>`
		return ctx.HTML(200, html)
	}

	// With fith renderer, show error:
	// return ctx.Render(401, "auth/login", map[string]interface{}{
	//     "title": "Login",
	//     "error": "Invalid credentials",
	// })

	html := `<!DOCTYPE html>
<html>
<head><title>Login Failed</title></head>
<body>
	<h1>Login Failed</h1>
	<p>Invalid username or password.</p>
	<a href="/login">Try Again</a>
</body>
</html>`
	return ctx.HTML(401, html)
}

// RegisterFormHandler renders registration form
func RegisterFormHandler(ctx cosan.Context) error {
	html := `<!DOCTYPE html>
<html>
<head><title>Register</title></head>
<body>
	<h1>Register</h1>
	<form method="POST" action="/register">
		<div>
			<label>Username:</label>
			<input type="text" name="username" required>
		</div>
		<div>
			<label>Email:</label>
			<input type="email" name="email" required>
		</div>
		<div>
			<label>Password:</label>
			<input type="password" name="password" required>
		</div>
		<button type="submit">Register</button>
	</form>
	<p><a href="/login">Login</a> | <a href="/">Home</a></p>
</body>
</html>`
	return ctx.HTML(200, html)
}

// UserPartialHandler renders a partial template (for HTMX/AJAX)
func UserPartialHandler(ctx cosan.Context) error {
	id := ctx.Param("id")

	// With fith renderer:
	// return ctx.Render(200, "partials/user", map[string]interface{}{
	//     "user": user,
	// })

	html := `<div class="user-card">
	<h3>User ` + id + `</h3>
	<p>Username: john_doe</p>
	<p>Email: john@example.com</p>
</div>`
	return ctx.HTML(200, html)
}

// API Handlers (JSON responses)

func GetUserAPIHandler(ctx cosan.Context) error {
	id := ctx.Param("id")
	user := User{
		ID:       1,
		Username: "john_doe",
		Email:    "john@example.com",
	}

	return ctx.JSON(200, map[string]interface{}{
		"id":   id,
		"user": user,
	})
}

func GetProductsAPIHandler(ctx cosan.Context) error {
	products := []Product{
		{ID: 1, Name: "Laptop", Price: 999.99},
		{ID: 2, Name: "Mouse", Price: 29.99},
		{ID: 3, Name: "Keyboard", Price: 79.99},
	}

	return ctx.JSON(200, map[string]interface{}{
		"products": products,
		"total":    len(products),
	})
}
