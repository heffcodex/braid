package middleware

import (
	"github.com/gofiber/fiber/v2"

	"github.com/heffcodex/braid/validator"
)

func Validator(v *validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		validator.V(c, v)
		return c.Next()
	}
}
