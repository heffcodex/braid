package middleware

import (
	"github.com/gofiber/fiber/v2"

	"github.com/heffcodex/braid/validation"
)

func Validator(v *validation.Validator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		validation.V(c, v)
		return c.Next()
	}
}
