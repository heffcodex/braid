package response

import (
	"github.com/gofiber/fiber/v2"
	"github.com/heffcodex/braid/status"
)

var _ IResponse = (*Raw)(nil)

type Raw struct {
	s    *status.Status
	data []byte
}

func NewRaw(s *status.Status, data []byte) *Raw {
	return &Raw{s: s, data: data}
}

func (r *Raw) Send(c *fiber.Ctx) error {
	c.Status(r.s.GetHTTP())
	return c.Send(r.data)
}
