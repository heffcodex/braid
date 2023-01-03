package response

import (
	"fmt"
	gojson "github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/heffcodex/braid/status"
)

var (
	_ Sender = (*JSONError)(nil)
	_ error  = (*JSONError)(nil)
)

type JSONErrorFormat struct {
	Code    status.Code `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    any         `json:"data,omitempty"`
}

type JSONError struct {
	s           *status.Status
	internalErr error
	data        any
}

func NewJSONError(s *status.Status, internal error, data ...any) *JSONError {
	var _data any

	if len(data) > 1 {
		panic("too many arguments")
	} else if len(data) == 1 {
		_data = data[0]
	}

	return &JSONError{s: s, internalErr: internal, data: _data}
}

func (r *JSONError) Send(c *fiber.Ctx) error {
	raw, err := gojson.MarshalContext(c.Context(), JSONFormat{
		Error: &JSONErrorFormat{
			Code:    r.s.Code(),
			Message: r.s.Message(),
			Data:    r.data,
		},
	})
	if err != nil {
		return err
	}

	c.Status(r.s.GetHTTP())
	c.Response().SetBodyRaw(raw)
	c.Response().Header.SetContentType(fiber.MIMEApplicationJSON)

	return nil
}

func (r *JSONError) IsInternal() bool {
	return r.internalErr != nil
}

func (r *JSONError) InternalError() error {
	return r.internalErr
}

func (r *JSONError) Error() string {
	str := fmt.Sprintf("[%d] %s", r.s.Code(), r.s.Message())
	if r.internalErr != nil {
		str = fmt.Sprintf("%s: %s", str, r.internalErr.Error())
	}

	return str
}
