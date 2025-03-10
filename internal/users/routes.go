package users

import (
	"backend/internal/users/handler"
	"backend/internal/users/handler/dto"
	"backend/pkg/config"
	"backend/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterUserRoutes(api fiber.Router, cfg *config.Config, db *gorm.DB) {
	userHandler := handler.InitUserHandler(db)

	users := api.Group("/users")
	{
		users.Get("/me", userHandler.GetMe)
		users.Put("/me", middleware.ValidateRequest(new(dto.UpdateUserRequest)), userHandler.UpdateMe)
	}
}

func RegisterAuthRoutes(api fiber.Router, cfg *config.Config, db *gorm.DB) {
	authHandler := handler.InitAuthHandler(cfg, db)

	auth := api.Group("/auth")
	{
		auth.Post("/register", middleware.ValidateRequest(new(dto.RegisterRequest)), authHandler.Register)
		auth.Post("/login", middleware.ValidateRequest(new(dto.LoginRequest)), authHandler.Login)
		auth.Get("/oauth/:provider", authHandler.OAuthSignIn)
		auth.Get("/callback/:provider", authHandler.OAuthCallback)
		auth.Post("/logout", authHandler.Logout)
		auth.Get("/session", authHandler.CheckSession)
	}
}
