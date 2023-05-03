package braid

import (
	"github.com/gofiber/fiber/v2"
	fibersession "github.com/gofiber/fiber/v2/middleware/session"

	"github.com/heffcodex/braid/session"
	"github.com/heffcodex/braid/validation"
)

func V(c *fiber.Ctx) *validation.ContextValidator {
	return validation.V(c)
}

func S(c *fiber.Ctx) *fibersession.Session {
	return session.S(c)
}
