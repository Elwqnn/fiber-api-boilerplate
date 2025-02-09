package main

import (
	"fiber-api-boilerplate/internal/config"
	"fiber-api-boilerplate/internal/database"
	"fiber-api-boilerplate/internal/handler/response"
	"fiber-api-boilerplate/internal/routes"
	"fiber-api-boilerplate/pkg/middleware"
	"log"
	"os"
	"os/signal"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	cfg := config.LoadConfig()

	store := config.SetupSessionStore()
	if store == nil {
		log.Fatal("Failed to setup session store")
	}

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
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true, // Required for cookies
		ExposeHeaders:    "Set-Cookie",
		MaxAge:           300,
	}))
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(middleware.HandleSession(store))

	// Routes
	routes.SetupRoutes(app, cfg, db)

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
