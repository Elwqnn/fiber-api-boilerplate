package middleware

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

func ValidateRequest(payload interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := c.BodyParser(payload); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
		}

		if err := validate.Struct(payload); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		c.Locals("payload", payload)
		return c.Next()
	}
}
