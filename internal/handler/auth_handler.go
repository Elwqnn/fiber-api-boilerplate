package handler

import (
	"fiber-api-boilerplate/internal/handler/dto"
	"fiber-api-boilerplate/internal/handler/response"
	"fiber-api-boilerplate/internal/model"
	"fiber-api-boilerplate/internal/service"
	"fiber-api-boilerplate/internal/repository"
	"fiber-api-boilerplate/internal/config"
	"time"

	"gorm.io/gorm"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type AuthHandler struct {
	authService service.AuthService
	userService service.UserService
}

func NewAuthHandler(authService service.AuthService, userService service.UserService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userService: userService,
	}
}

func InitAuthHandler(cfg *config.Config, db *gorm.DB) *AuthHandler {
	userRepo := repository.NewUserRepository(db)
	accountRepo := repository.NewAccountRepository(db)
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(
		userRepo,
		accountRepo,
		cfg.JWTSecret,
	)

	return NewAuthHandler(authService, userService)
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	req := c.Locals("payload").(*dto.RegisterRequest)

	user := &model.User{
		Name:  req.Name,
		Email: req.Email,
	}

	if err := h.authService.Register(c.Context(), user, req.Password); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	token, err := h.authService.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, dto.AuthResponse{
		Token: token,
		User:  *user,
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	req := c.Locals("payload").(*dto.LoginRequest)

	token, err := h.authService.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "Invalid credentials")
	}

	user, err := h.userService.GetByEmail(c.Context(), req.Email)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	// Get session from context
    sess, err := c.Locals("store").(*session.Store).Get(c)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to retreive session from locals")
	}

	// Set session data
	sess.Set("user_id", user.ID.String())
	sess.Set("email", user.Email)
	sess.Set("role", user.Role)
	sess.Set("last_activity", time.Now().Unix())
	sess.Set("expires_at", time.Now().Add(time.Hour*24).Unix())

    if err := sess.Save(); err != nil {
        return response.Error(c, fiber.StatusInternalServerError, "Failed to save session")
    }

	return response.Success(c, dto.AuthResponse{
		Token: token,
		User:  *user,
	})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	sess, err := c.Locals("store").(*session.Store).Get(c)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to retreive session from locals")
	}

	if err := sess.Destroy(); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to destroy session")
	}

	// Clear session cookie
	c.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		SameSite: "Lax",
		Path:     "/",
	})

	return response.Success(c, nil)
}

func (h *AuthHandler) OAuthSignIn(c *fiber.Ctx) error {
	provider := c.Params("provider")
	redirectURL, err := h.authService.GetOAuthRedirectURL(provider)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}
	return c.Redirect(redirectURL)
}

func (h *AuthHandler) OAuthCallback(c *fiber.Ctx) error {
	provider := c.Params("provider")
	code := c.Query("code")
	state := c.Query("state")

	token, user, err := h.authService.HandleOAuthCallback(c.Context(), provider, code, state)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, dto.AuthResponse{
		Token: token,
		User:  *user,
	})
}

func (h *AuthHandler) CheckSession(c *fiber.Ctx) error {
	return response.Success(c, fiber.Map{
		"user_id":       c.Locals("user_id"),
		"email":         c.Locals("email"),
		"role":          c.Locals("role"),
		"last_activity": c.Locals("last_activity"),
		"expires_at":    c.Locals("expires_at"),
	})
}
