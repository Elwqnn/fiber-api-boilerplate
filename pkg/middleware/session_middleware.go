package middleware

import (
	"fiber-api-boilerplate/internal/handler/response"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func HandleSession(store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get or create session
		sess, err := store.Get(c)
		if err != nil {
			return response.Error(c, fiber.StatusInternalServerError, "Session error")
		}

		sess.Set("last_activity", time.Now().Unix())

		if err := sess.Save(); err != nil {
			return response.Error(c, fiber.StatusInternalServerError, "Failed to save session")
		}

		c.Locals("store", store)
		c.Next()
		return nil
	}
}
