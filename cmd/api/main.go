package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"
	"backend/internal/users"
	"backend/pkg/config"
	"backend/pkg/database"
	"backend/pkg/middleware"
	"backend/pkg/response"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/gorm"
)

func main() {
	// Load config from .env file
	cfg := config.LoadConfig()

	// Connect to database
	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Fiber instance
	app := fiber.New(fiber.Config{
		ErrorHandler: customErrorHandler,
		JSONEncoder:  json.Marshal,   // optimized JSON serialization
		JSONDecoder:  json.Unmarshal, // optimized JSON deserialization
	})

	// Middlewares
	setupMiddlewares(app)

	// Routes
	setupRoutes(app, cfg, db)

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		<-c
		log.Println("Gracefully shutting down...")
		_ = app.Shutdown()
	}()

	// Start server
	log.Fatal(app.Listen(":" + cfg.Port))
}

// setupRoutes initializes all routes for the application
func setupRoutes(app *fiber.App, cfg *config.Config, db *gorm.DB) {
	api := app.Group(fmt.Sprintf("/api/%s", strings.ToLower(cfg.Env)))

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	users.RegisterAuthRoutes(api, cfg, db)
	users.RegisterUserRoutes(api, cfg, db)
}

// setupMiddlewares initializes all mandatory middlewares for the application
func setupMiddlewares(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
		ExposeHeaders:    "Set-Cookie",
		MaxAge:           300,
	}))
	app.Use(favicon.New())
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(limiter.New(limiter.Config{
		Max:               20,
		Expiration:        30 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
	}))

	store := config.SetupSessionStore()
	if store == nil {
		log.Fatal("Failed to setup session store")
	}
	app.Use(middleware.HandleSession(store))
}

// customErrorHandler allows for a standardized error response
func customErrorHandler(c *fiber.Ctx, err error) error {
	// Get status code
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	// Log error for internal server errors
	if code == fiber.StatusInternalServerError {
		log.Printf("Internal Server Error: %v", err)
	}

	// Return standardized error response
	return c.Status(code).JSON(response.Response{
		Success: false,
		Message: err.Error(),
	})
}
