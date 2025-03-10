package middleware

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

type ErrorResponse struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func getErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("Minimum length is %s", fe.Param())
	case "max":
		return fmt.Sprintf("Maximum length is %s", fe.Param())
	}
	return fe.Error() // default error
}

func ValidateRequest(payload interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := c.BodyParser(payload); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
		}

		if err := validate.Struct(payload); err != nil {
			var errorMsgs []string
			for _, err := range err.(validator.ValidationErrors) {
				field := strings.ToLower(err.Field())
				errorMsgs = append(errorMsgs, fmt.Sprintf("%s: %s", field, getErrorMessage(err)))
			}
			return fiber.NewError(fiber.StatusBadRequest, strings.Join(errorMsgs, "; "))
		}

		c.Locals("payload", payload)
		return c.Next()
	}
}
