package middleware

import (
	"fmt"

	"github.com/heffcodex/braid/session"

	"github.com/gofiber/fiber/v2"
	fiber_session "github.com/gofiber/fiber/v2/middleware/session"
)

func Session(store *fiber_session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, err := store.Get(c)
		if err != nil {
			return fmt.Errorf("get session: %w", err)
		}

		session.S(c, sess)
		defer func() { _ = sess.Save() }()

		return c.Next()
	}
}
