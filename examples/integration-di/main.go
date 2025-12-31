// Package main demonstrates integration with toutago-nasc-dependency-injector.
package main

import (
	"fmt"
	"log"

	"github.com/toutaio/toutago-cosan-router/pkg/cosan"
)

// Service interfaces
type UserService interface {
	GetUser(id int) (*User, error)
	CreateUser(user *User) error
}

type ProductService interface {
	GetProduct(id int) (*Product, error)
	ListProducts() ([]*Product, error)
}

type Logger interface {
	Info(msg string)
	Error(msg string)
}

// Implementations
type userServiceImpl struct {
	logger Logger
}

func NewUserService(logger Logger) UserService {
	return &userServiceImpl{logger: logger}
}

func (s *userServiceImpl) GetUser(id int) (*User, error) {
	s.logger.Info(fmt.Sprintf("Fetching user %d", id))
	return &User{
		ID:       id,
		Username: fmt.Sprintf("user_%d", id),
		Email:    fmt.Sprintf("user%d@example.com", id),
	}, nil
}

func (s *userServiceImpl) CreateUser(user *User) error {
	s.logger.Info(fmt.Sprintf("Creating user %s", user.Username))
	return nil
}

type productServiceImpl struct {
	logger Logger
}

func NewProductService(logger Logger) ProductService {
	return &productServiceImpl{logger: logger}
}

func (s *productServiceImpl) GetProduct(id int) (*Product, error) {
	s.logger.Info(fmt.Sprintf("Fetching product %d", id))
	return &Product{
		ID:    id,
		Name:  fmt.Sprintf("Product %d", id),
		Price: 99.99,
	}, nil
}

func (s *productServiceImpl) ListProducts() ([]*Product, error) {
	s.logger.Info("Listing all products")
	return []*Product{
		{ID: 1, Name: "Product 1", Price: 99.99},
		{ID: 2, Name: "Product 2", Price: 149.99},
	}, nil
}

type simpleLogger struct{}

func NewLogger() Logger {
	return &simpleLogger{}
}

func (l *simpleLogger) Info(msg string) {
	log.Printf("[INFO] %s", msg)
}

func (l *simpleLogger) Error(msg string) {
	log.Printf("[ERROR] %s", msg)
}

// Models
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// Controllers with dependency injection
type UserController struct {
	userService UserService
}

func NewUserController(userService UserService) *UserController {
	return &UserController{userService: userService}
}

func (c *UserController) GetUser(ctx cosan.Context) error {
	_ = ctx.Param("id") // Path parameter

	// With NASC DI, parameters can be auto-injected:
	// func (c *UserController) GetUser(ctx cosan.Context, id int) error {
	//     user, err := c.userService.GetUser(id)
	//     ...
	// }

	userID := 1 // Simplified conversion
	user, err := c.userService.GetUser(userID)
	if err != nil {
		return ctx.JSON(500, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(200, user)
}

func (c *UserController) CreateUser(ctx cosan.Context) error {
	var user User
	if err := ctx.Bind(&user); err != nil {
		return ctx.JSON(400, map[string]string{"error": "Invalid request"})
	}

	if err := c.userService.CreateUser(&user); err != nil {
		return ctx.JSON(500, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(201, user)
}

type ProductController struct {
	productService ProductService
}

func NewProductController(productService ProductService) *ProductController {
	return &ProductController{productService: productService}
}

func (c *ProductController) GetProduct(ctx cosan.Context) error {
	// With NASC DI and datamapper:
	// func (c *ProductController) GetProduct(ctx cosan.Context, id int) error {
	//     product, err := c.productService.GetProduct(id)
	//     return ctx.JSON(200, product)
	// }

	_ = ctx.Param("id") // Path parameter
	productID := 1      // Simplified
	product, err := c.productService.GetProduct(productID)
	if err != nil {
		return ctx.JSON(500, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(200, product)
}

func (c *ProductController) ListProducts(ctx cosan.Context) error {
	products, err := c.productService.ListProducts()
	if err != nil {
		return ctx.JSON(500, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(200, map[string]interface{}{
		"products": products,
		"total":    len(products),
	})
}

func main() {
	// Manual dependency injection (in production, use NASC container)
	logger := NewLogger()
	userService := NewUserService(logger)
	productService := NewProductService(logger)

	// Create controllers with injected dependencies
	userController := NewUserController(userService)
	productController := NewProductController(productService)

	// With NASC DI container:
	// container := nasc.NewContainer()
	// container.Register(NewLogger)
	// container.Register(NewUserService)
	// container.Register(NewProductService)
	// container.Register(NewUserController)
	// container.Register(NewProductController)
	//
	// router := cosan.New()
	// router.SetContainer(container)
	//
	// Auto-resolve and inject dependencies:
	// router.GET("/users/:id", container.Resolve(UserController).GetUser)

	router := cosan.New()

	// User routes
	router.GET("/users/:id", userController.GetUser)
	router.POST("/users", userController.CreateUser)

	// Product routes
	router.GET("/products/:id", productController.GetProduct)
	router.GET("/products", productController.ListProducts)

	// Demonstrate middleware with DI
	router.Use(LoggingMiddleware(logger))

	log.Println("Server starting on http://localhost:8080")
	log.Println("Examples:")
	log.Println("  GET  http://localhost:8080/users/1")
	log.Println("  POST http://localhost:8080/users")
	log.Println("  GET  http://localhost:8080/products")
	log.Fatal(router.Listen(":8080"))
}

// LoggingMiddleware demonstrates middleware with injected dependencies
func LoggingMiddleware(logger Logger) cosan.MiddlewareFunc {
	return func(next cosan.HandlerFunc) cosan.HandlerFunc {
		return func(ctx cosan.Context) error {
			logger.Info(fmt.Sprintf("%s %s", ctx.Request().Method, ctx.Request().URL.Path))
			return next(ctx)
		}
	}
}
