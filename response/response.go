package response

import (
	"github.com/gofiber/fiber/v2"
)

type IResponse interface {
	Send(c *fiber.Ctx) error
}
