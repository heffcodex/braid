package response

import (
	"github.com/gofiber/fiber/v2"
	"github.com/heffcodex/braid/status"
)

var _ IResponse = (*Empty)(nil)

type Empty struct {
	s *status.Status
}

func NewEmpty(s *status.Status) *Empty {
	return &Empty{s: s}
}

func (r *Empty) Send(c *fiber.Ctx) error {
	c.Status(r.s.GetHTTP())
	return nil
}
