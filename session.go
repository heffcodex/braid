package braid

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

const LocalSession = "session"

func S(c *fiber.Ctx) *session.Store {
	return getSession(c)
}

func getSession(c *fiber.Ctx) *session.Store {
	return c.Locals(LocalSession).(*session.Store)
}

func setSession(c *fiber.Ctx, s *session.Store) {
	c.Locals(LocalSession, s)
}

func SessionMiddleware(s *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		setSession(c, s)
		return c.Next()
	}
}
