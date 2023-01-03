package response

import (
	"github.com/gofiber/fiber/v2"
)

type Sender interface {
	Send(c *fiber.Ctx) error
}
