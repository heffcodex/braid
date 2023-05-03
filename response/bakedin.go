package response

import (
	"github.com/gofiber/fiber/v2"

	"github.com/heffcodex/braid/status"
)

const (
	HeaderContentDescription       = "Content-Description"
	ContentDescriptionFileTransfer = "File Transfer"
)

// 2xx

func OK(c *fiber.Ctx, data any, meta ...any) error {
	return NewJSON(status.OK, data, meta...).Send(c)
}

func FileTransfer(c *fiber.Ctx, filename string, data []byte) error {
	c.Set(HeaderContentDescription, ContentDescriptionFileTransfer)
	c.Attachment(filename)

	return NewRaw(status.OK, data).Send(c)
}

func Created(c *fiber.Ctx, data any, meta ...any) error {
	return NewJSON(status.Created, data, meta...).Send(c)
}

func NoContent(c *fiber.Ctx) error {
	return NewEmpty(status.NoContent).Send(c)
}

// All functions below with E prefix return not an error of sending the response to the client,
// but the error-like ResponseError object to be handled by the server-wide error handler:

// 4xx

func EBadRequest(sub *status.Status, data ...any) error {
	return NewJSONError(sub, nil, data...)
}

func EUnauthorized(data ...any) error {
	return NewJSONError(status.Unauthorized, nil, data...)
}

func EForbidden(data ...any) error {
	return NewJSONError(status.Forbidden, nil, data...)
}

func ENotFound(data ...any) error {
	return NewJSONError(status.NotFound, nil, data...)
}

func EConflict(data ...any) error {
	return NewJSONError(status.Conflict, nil, data...)
}

// 5xx

func EInternal(err error) error {
	return NewJSONError(status.InternalError, err)
}
