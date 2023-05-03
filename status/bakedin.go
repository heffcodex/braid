package status

import "github.com/gofiber/fiber/v2"

const (
	CodeOK        Code = fiber.StatusOK
	CodeCreated   Code = fiber.StatusCreated
	CodeNoContent Code = fiber.StatusNoContent
)

const (
	_ Code = iota + fiber.StatusBadRequest*1000
	CodeCSRFTokenMismatch
	CodeValidationFail
	CodeInvalidQuery
	CodeInvalidPayload
)

var (
	OK        = New(CodeOK)
	Created   = New(CodeCreated)
	NoContent = New(CodeNoContent)
)

var (
	BadRequest        = FromFiber(fiber.ErrBadRequest)
	CSRFTokenMismatch = New(CodeCSRFTokenMismatch, "CSRF Token Mismatch")
	ValidationFail    = New(CodeValidationFail, "Validation Failed")
	InvalidQuery      = New(CodeInvalidQuery, "Invalid Query")
	InvalidPayload    = New(CodeInvalidPayload, "Invalid Payload")
)

var (
	Unauthorized = FromFiber(fiber.ErrUnauthorized)
	Forbidden    = FromFiber(fiber.ErrForbidden)
	NotFound     = FromFiber(fiber.ErrNotFound)
	Conflict     = FromFiber(fiber.ErrConflict)
)

var (
	InternalError = FromFiber(fiber.ErrInternalServerError)
)
