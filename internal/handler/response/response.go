package response

import "github.com/gofiber/fiber/v2"

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(c *fiber.Ctx, data interface{}) error {
	return c.JSON(Response{
		Success: true,
		Data:    data,
	})
}

func Error(c *fiber.Ctx, code int, message string) error {
	return c.Status(code).JSON(Response{
		Success: false,
		Message: message,
	})
}
