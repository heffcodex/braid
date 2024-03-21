package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	fiber_session "github.com/gofiber/fiber/v2/middleware/session"

	"github.com/heffcodex/braid/session"
	"github.com/heffcodex/braid/vars"
)

func Session(store *fiber_session.Store) fiber.Handler {
	name := strings.SplitN(store.KeyLookup, ":", 2)[1]

	return func(c *fiber.Ctx) error {
		sess, err := store.Get(c)
		if err != nil {
			return fmt.Errorf("get session: %w", err)
		}

		session.S(c, sess)
		defer func() {
			if v, ok := c.Locals(vars.LocalSessionOmit).(bool); ok && v {
				_ = sess.Destroy()
				c.Response().Header.Del(name)
				c.Response().Header.DelCookie(name)
			} else {
				_ = sess.Save()
			}
		}()

		return c.Next()
	}
}
