package braid

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

const LocalSession = "session"

func S(c *fiber.Ctx) *session.Session {
	return getSession(c)
}

func getSession(c *fiber.Ctx) *session.Session {
	return c.Locals(LocalSession).(*session.Session)
}

func setSession(c *fiber.Ctx, s *session.Session) {
	c.Locals(LocalSession, s)
}

func SessionMiddleware(store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, err := store.Get(c)
		if err != nil {
			return fmt.Errorf("get session: %w", err)
		}

		setSession(c, sess)
		defer func() { _ = sess.Save() }()

		return c.Next()
	}
}
