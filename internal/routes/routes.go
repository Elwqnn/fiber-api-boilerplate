package routes

import (
	"fiber-api-boilerplate/internal/handler"
	"fiber-api-boilerplate/internal/handler/dto"
	"fiber-api-boilerplate/internal/config"
	"fiber-api-boilerplate/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, cfg *config.Config, db *gorm.DB) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	v1 := app.Group("/api/v1")

	// Auth routes
	authHandler := handler.InitAuthHandler(cfg, db)
	auth := v1.Group("/auth")
	auth.Post("/register", middleware.ValidateRequest(new(dto.RegisterRequest)), authHandler.Register)
	auth.Post("/login", middleware.ValidateRequest(new(dto.LoginRequest)), authHandler.Login)
	auth.Get("/oauth/:provider", authHandler.OAuthSignIn)
	auth.Get("/callback/:provider", authHandler.OAuthCallback)
	auth.Post("/logout", middleware.RequireJWT(cfg.JWTSecret), authHandler.Logout)
	auth.Get("/session", middleware.RequireJWT(cfg.JWTSecret), authHandler.CheckSession)

	// Protected routes
	v1.Use(middleware.RequireJWT(cfg.JWTSecret))
	v1.Use(middleware.RequireRole("user", "admin"))

	// User routes
	userHandler := handler.InitUserHandler(db)
	users := v1.Group("/users")
	users.Get("/me", userHandler.GetMe)
	users.Put("/me", middleware.ValidateRequest(new(dto.UpdateUserRequest)), userHandler.UpdateMe)
}
