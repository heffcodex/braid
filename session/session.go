package session

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"

	"github.com/heffcodex/braid/vars"
)

func S(c *fiber.Ctx, set ...*session.Session) *session.Session {
	if len(set) > 1 {
		panic("too many arguments")
	} else if len(set) == 1 {
		setSession(c, set[0])
		return set[0]
	}

	return getSession(c)
}

func Omit(c *fiber.Ctx) {
	c.Locals(vars.LocalSessionOmit, true)
}

func getSession(c *fiber.Ctx) *session.Session {
	return c.Locals(vars.LocalSession).(*session.Session)
}

func setSession(c *fiber.Ctx, s *session.Session) {
	c.Locals(vars.LocalSession, s)
}
