package handler

import (
	"log"
	"backend/internal/users/handler/dto"
	"backend/internal/users/repository"
	"backend/internal/users/service"
	"backend/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func InitUserHandler(db *gorm.DB) *UserHandler {
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	return NewUserHandler(userService)
}

func (h *UserHandler) GetMe(c *fiber.Ctx) error {
	userID := c.Locals("user_id")

	userIDToString, err := uuid.Parse(userID.(string))
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Invalid user ID format")
	}

	user, err := h.userService.GetByID(c.Context(), userIDToString)
	if err != nil {
		log.Println("Error: ", err)
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, user)
}

func (h *UserHandler) UpdateMe(c *fiber.Ctx) error {
	userID := c.Locals("user_id")

	userIDToString, err := uuid.Parse(userID.(string))
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Invalid user ID format")
	}

	req := c.Locals("payload").(*dto.UpdateUserRequest)
	user, err := h.userService.GetByID(c.Context(), userIDToString)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Image != "" {
		user.Image = req.Image
	}

	if err := h.userService.Update(c.Context(), user); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, user)
}
