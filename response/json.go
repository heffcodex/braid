package response

import (
	gojson "github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"

	"github.com/heffcodex/braid/status"
)

var _ Sender = (*JSON)(nil)

type JSONFormat struct {
	Meta  any              `json:"meta,omitempty"`
	Data  any              `json:"data,omitempty"`
	Error *JSONErrorFormat `json:"error,omitempty"`
}

type JSON struct {
	s          *status.Status
	meta, data any
}

func NewJSON(s *status.Status, data any, meta ...any) *JSON {
	var _meta any

	if len(meta) > 1 {
		panic("too many arguments")
	} else if len(meta) == 1 {
		_meta = meta[0]
	}

	return &JSON{
		s:    s,
		meta: _meta,
		data: data,
	}
}

func (r *JSON) Send(c *fiber.Ctx) error {
	raw, err := gojson.MarshalContext(c.Context(), JSONFormat{
		Meta: r.meta,
		Data: r.data,
	})
	if err != nil {
		return err
	}

	c.Status(r.s.Code().HTTP())
	c.Response().SetBodyRaw(raw)
	c.Response().Header.SetContentType(fiber.MIMEApplicationJSON)

	return nil
}
