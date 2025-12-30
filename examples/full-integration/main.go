// Package main demonstrates full integration of Cosan router with all toutago components:
// - toutago-datamapper for parameter binding and validation
// - toutago-fith-renderer for HTML rendering
// - toutago-nasc-dependency-injector for dependency injection
package main

import (
	"fmt"
	"log"

	"github.com/toutaio/toutago-cosan-router/pkg/cosan"
)

// Domain models
type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username" validate:"required,min=3,max=50"`
	Email     string `json:"email" validate:"required,email"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,min=0"`
	Category    string  `json:"category" validate:"required"`
}

// Service layer interfaces
type UserRepository interface {
	FindByID(id int) (*User, error)
	Create(user *User) error
	Update(user *User) error
	Delete(id int) error
}

type ProductRepository interface {
	FindByID(id int) (*Product, error)
	FindByCategory(category string) ([]*Product, error)
	Create(product *Product) error
}

type Logger interface {
	Info(msg string)
	Error(msg string)
}

// Service implementations (would be in separate files)
type userRepository struct {
	logger Logger
	users  map[int]*User
}

func NewUserRepository(logger Logger) UserRepository {
	return &userRepository{
		logger: logger,
		users:  make(map[int]*User),
	}
}

func (r *userRepository) FindByID(id int) (*User, error) {
	r.logger.Info(fmt.Sprintf("Finding user %d", id))
	if user, ok := r.users[id]; ok {
		return user, nil
	}
	return &User{ID: id, Username: "demo", Email: "demo@example.com"}, nil
}

func (r *userRepository) Create(user *User) error {
	r.logger.Info(fmt.Sprintf("Creating user %s", user.Username))
	user.ID = len(r.users) + 1
	r.users[user.ID] = user
	return nil
}

func (r *userRepository) Update(user *User) error {
	r.logger.Info(fmt.Sprintf("Updating user %d", user.ID))
	r.users[user.ID] = user
	return nil
}

func (r *userRepository) Delete(id int) error {
	r.logger.Info(fmt.Sprintf("Deleting user %d", id))
	delete(r.users, id)
	return nil
}

type productRepository struct {
	logger   Logger
	products map[int]*Product
}

func NewProductRepository(logger Logger) ProductRepository {
	return &productRepository{
		logger:   logger,
		products: make(map[int]*Product),
	}
}

func (r *productRepository) FindByID(id int) (*Product, error) {
	if product, ok := r.products[id]; ok {
		return product, nil
	}
	return &Product{ID: id, Name: "Demo Product", Price: 99.99}, nil
}

func (r *productRepository) FindByCategory(category string) ([]*Product, error) {
	r.logger.Info(fmt.Sprintf("Finding products in category %s", category))
	return []*Product{
		{ID: 1, Name: "Product 1", Price: 99.99, Category: category},
		{ID: 2, Name: "Product 2", Price: 149.99, Category: category},
	}, nil
}

func (r *productRepository) Create(product *Product) error {
	r.logger.Info(fmt.Sprintf("Creating product %s", product.Name))
	product.ID = len(r.products) + 1
	r.products[product.ID] = product
	return nil
}

type simpleLogger struct{}

func NewLogger() Logger {
	return &simpleLogger{}
}

func (l *simpleLogger) Info(msg string)  { log.Printf("[INFO] %s", msg) }
func (l *simpleLogger) Error(msg string) { log.Printf("[ERROR] %s", msg) }

// Controllers
type UserController struct {
	repo UserRepository
}

func NewUserController(repo UserRepository) *UserController {
	return &UserController{repo: repo}
}

func (c *UserController) Show(ctx cosan.Context) error {
	id := ctx.Param("id")
	user, err := c.repo.FindByID(1) // Simplified
	if err != nil {
		return ctx.JSON(500, map[string]string{"error": err.Error()})
	}

	// With Fith renderer:
	// return ctx.Render(200, "users/show", map[string]interface{}{
	//     "title": "User Profile",
	//     "user":  user,
	// })

	html := fmt.Sprintf(`<!DOCTYPE html>
<html><head><title>User %s</title></head>
<body>
	<h1>User Profile</h1>
	<p>Username: %s</p>
	<p>Email: %s</p>
	<p>Name: %s %s</p>
	<a href="/users">Back to list</a>
</body></html>`, id, user.Username, user.Email, user.FirstName, user.LastName)

	return ctx.HTML(200, html)
}

func (c *UserController) Create(ctx cosan.Context) error {
	var user User

	// With datamapper:
	// if err := ctx.Bind(&user); err != nil {
	//     return ctx.JSON(400, map[string]interface{}{
	//         "error": "Validation failed",
	//         "details": err,
	//     })
	// }

	if err := ctx.Bind(&user); err != nil {
		return ctx.JSON(400, map[string]string{"error": err.Error()})
	}

	if err := c.repo.Create(&user); err != nil {
		return ctx.JSON(500, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(201, user)
}

func (c *UserController) Update(ctx cosan.Context) error {
	var user User
	if err := ctx.Bind(&user); err != nil {
		return ctx.JSON(400, map[string]string{"error": err.Error()})
	}

	if err := c.repo.Update(&user); err != nil {
		return ctx.JSON(500, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(200, user)
}

func (c *UserController) Delete(ctx cosan.Context) error {
	id := 1 // Simplified
	if err := c.repo.Delete(id); err != nil {
		return ctx.JSON(500, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(204, nil)
}

type ProductController struct {
	repo ProductRepository
}

func NewProductController(repo ProductRepository) *ProductController {
	return &ProductController{repo: repo}
}

func (c *ProductController) Index(ctx cosan.Context) error {
	category := ctx.Query("category")
	if category == "" {
		category = "all"
	}

	products, err := c.repo.FindByCategory(category)
	if err != nil {
		return ctx.JSON(500, map[string]string{"error": err.Error()})
	}

	// Check Accept header for content negotiation
	accept := ctx.Request().Header.Get("Accept")
	if accept == "application/json" {
		return ctx.JSON(200, map[string]interface{}{
			"products": products,
			"total":    len(products),
		})
	}

	// With Fith renderer:
	// return ctx.Render(200, "products/index", map[string]interface{}{
	//     "title":    "Products",
	//     "products": products,
	//     "category": category,
	// })

	html := `<!DOCTYPE html><html><head><title>Products</title></head><body>
		<h1>Products</h1>
		<ul><li>Product 1 - $99.99</li><li>Product 2 - $149.99</li></ul>
	</body></html>`

	return ctx.HTML(200, html)
}

func (c *ProductController) Show(ctx cosan.Context) error {
	id := 1 // Simplified
	product, err := c.repo.FindByID(id)
	if err != nil {
		return ctx.JSON(500, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(200, product)
}

func (c *ProductController) Create(ctx cosan.Context) error {
	var product Product
	if err := ctx.Bind(&product); err != nil {
		return ctx.JSON(400, map[string]string{"error": err.Error()})
	}

	if err := c.repo.Create(&product); err != nil {
		return ctx.JSON(500, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(201, product)
}

func main() {
	// Setup DI container (manual wiring for demo)
	// In production with NASC:
	// container := nasc.NewContainer()
	// container.Register(NewLogger)
	// container.Register(NewUserRepository)
	// container.Register(NewProductRepository)
	// container.Register(NewUserController)
	// container.Register(NewProductController)

	logger := NewLogger()
	userRepo := NewUserRepository(logger)
	productRepo := NewProductRepository(logger)
	userCtrl := NewUserController(userRepo)
	productCtrl := NewProductController(productRepo)

	// Setup router
	router := cosan.New()

	// With full integration:
	// router.SetBinder(datamapper.NewBinder())
	// router.SetValidator(datamapper.NewValidator())
	// router.SetRenderer(fith.New(fith.Config{TemplateDir: "templates"}))
	// router.SetContainer(container)

	// Global middleware
	router.Use(LoggerMiddleware(logger))
	router.Use(RecoveryMiddleware())

	// Web routes (HTML)
	router.GET("/", HomeHandler)
	router.GET("/users/:id", userCtrl.Show)
	router.GET("/products", productCtrl.Index)

	// API routes (JSON)
	api := router.Group("/api/v1")
	api.Use(JSONMiddleware())

	// Users API
	users := api.Group("/users")
	users.GET("/:id", userCtrl.Show)
	users.POST("", userCtrl.Create)
	users.PUT("/:id", userCtrl.Update)
	users.DELETE("/:id", userCtrl.Delete)

	// Products API
	products := api.Group("/products")
	products.GET("", productCtrl.Index)
	products.GET("/:id", productCtrl.Show)
	products.POST("", productCtrl.Create)

	log.Println("Full integration example starting on http://localhost:8080")
	log.Println("Try:")
	log.Println("  Web:  http://localhost:8080/")
	log.Println("  Web:  http://localhost:8080/users/1")
	log.Println("  API:  http://localhost:8080/api/v1/users/1")
	log.Fatal(router.Listen(":8080"))
}

func HomeHandler(ctx cosan.Context) error {
	html := `<!DOCTYPE html><html><head><title>Home</title></head><body>
		<h1>Cosan Router - Full Integration</h1>
		<nav>
			<a href="/users/1">Users</a> |
			<a href="/products">Products</a> |
			<a href="/api/v1/users/1">API</a>
		</nav>
	</body></html>`
	return ctx.HTML(200, html)
}

func LoggerMiddleware(logger Logger) cosan.MiddlewareFunc {
	return func(next cosan.HandlerFunc) cosan.HandlerFunc {
		return func(ctx cosan.Context) error {
			logger.Info(fmt.Sprintf("%s %s", ctx.Request().Method, ctx.Request().URL.Path))
			return next(ctx)
		}
	}
}

func RecoveryMiddleware() cosan.MiddlewareFunc {
	return func(next cosan.HandlerFunc) cosan.HandlerFunc {
		return func(ctx cosan.Context) error {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Panic: %v", r)
					ctx.JSON(500, map[string]string{"error": "Internal server error"})
				}
			}()
			return next(ctx)
		}
	}
}

func JSONMiddleware() cosan.MiddlewareFunc {
	return func(next cosan.HandlerFunc) cosan.HandlerFunc {
		return func(ctx cosan.Context) error {
			ctx.Response().Header().Set("Content-Type", "application/json")
			return next(ctx)
		}
	}
}
