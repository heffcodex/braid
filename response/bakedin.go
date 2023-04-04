package response

import (
	"github.com/gofiber/fiber/v2"

	"github.com/heffcodex/braid/status"
)

const (
	HeaderContentDescription = "Content-Description"
)

// 2xx

func OK(c *fiber.Ctx, data any, meta ...any) error {
	return NewJSON(status.GenOK, data, meta...).Send(c)
}

func FileTransfer(c *fiber.Ctx, filename string, data []byte) error {
	c.Set(HeaderContentDescription, "File Transfer")
	c.Attachment(filename)

	return NewRaw(status.GenOK, data).Send(c)
}

func Created(c *fiber.Ctx, data any, meta ...any) error {
	return NewJSON(status.GenCreated, data, meta...).Send(c)
}

func NoContent(c *fiber.Ctx) error {
	return NewEmpty(status.GenNoContent).Send(c)
}

// All functions below with E prefix return not the error of sending the response to the client,
// but the error-wrapped ResponseError object to be handled by the server-wide error handler:

func eJSON(base, sub *status.Status, internal error, data ...any) error {
	if sub == nil {
		sub = base
	} else if !sub.Is(base) {
		panic("status must be a sub-status of `" + base.String() + "`")
	}

	return NewJSONError(sub, internal, data...)
}

// 4xx

func EBadRequest(sub *status.Status, data ...any) error {
	return eJSON(status.GenBadRequest, sub, nil, data...)
}

func EUnauthorized(data ...any) error {
	return eJSON(status.GenUnauthorized, nil, nil, data...)
}

func EForbidden(data ...any) error {
	return eJSON(status.GenForbidden, nil, nil, data...)
}

func ENotFound(data ...any) error {
	return eJSON(status.GenNotFound, nil, nil, data...)
}

func EConflict(data ...any) error {
	return eJSON(status.GenConflict, nil, nil, data...)
}

// 5xx

func EInternal(err error) error {
	return eJSON(status.GenInternalError, nil, err)
}
