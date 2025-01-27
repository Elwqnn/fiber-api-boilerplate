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

		// Update session activity
		sess.Set("last_activity", time.Now().Unix())

		if err := sess.Save(); err != nil {
			return response.Error(c, fiber.StatusInternalServerError, "Failed to save session")
		}

		// Set session data to locals
		c.Locals("store", store)
		c.Locals("last_activity", sess.Get("last_activity"))

		if sess.Get("user_id") != nil {
			c.Locals("user_id", sess.Get("user_id"))
			c.Locals("email", sess.Get("email"))
			c.Locals("role", sess.Get("role"))
			c.Locals("expires_at", sess.Get("expires_at"))
		}
		c.Next()
		return nil
	}
}
