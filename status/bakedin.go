package status

import "github.com/gofiber/fiber/v2"

const (
	CodeGenOK        Code = fiber.StatusOK
	CodeGenCreated   Code = fiber.StatusCreated
	CodeGenNoContent Code = fiber.StatusNoContent

	CodeGenBadRequest     Code = fiber.StatusBadRequest
	CodeCSRFTokenMismatch Code = fiber.StatusBadRequest*1000 + 1
	CodeValidationFail    Code = fiber.StatusBadRequest*1000 + 2
	CodeInvalidQuery      Code = fiber.StatusBadRequest*1000 + 3
	CodeInvalidPayload    Code = fiber.StatusBadRequest*1000 + 4

	CodeGenUnauthorized Code = fiber.StatusUnauthorized
	CodeGenForbidden    Code = fiber.StatusForbidden
	CodeGenNotFound     Code = fiber.StatusNotFound
	CodeGenConflict     Code = fiber.StatusConflict

	CodeGenInternalError Code = fiber.StatusInternalServerError
)

var (
	GenOK        = New(CodeGenOK).WithHTTP(fiber.StatusOK)
	GenCreated   = New(CodeGenCreated).WithHTTP(fiber.StatusCreated)
	GenNoContent = New(CodeGenNoContent).WithHTTP(fiber.StatusNoContent)
)

var (
	GenBadRequest     = New(CodeGenBadRequest, "Bad Request").WithHTTP(fiber.StatusBadRequest)
	CSRFTokenMismatch = New(CodeCSRFTokenMismatch, "CSRF Token Mismatch").AttachTo(GenBadRequest)
	ValidationFail    = New(CodeValidationFail, "Validation Failed").AttachTo(GenBadRequest)
	InvalidQuery      = New(CodeInvalidQuery, "Invalid Query").AttachTo(GenBadRequest)
	InvalidPayload    = New(CodeInvalidPayload, "Invalid Payload").AttachTo(GenBadRequest)
)

var (
	GenUnauthorized = New(CodeGenUnauthorized, "Unauthorized").WithHTTP(fiber.StatusUnauthorized)
	GenForbidden    = New(CodeGenForbidden, "Forbidden").WithHTTP(fiber.StatusForbidden)
	GenNotFound     = New(CodeGenNotFound, "Not Found").WithHTTP(fiber.StatusNotFound)
	GenConflict     = New(CodeGenConflict, "Conflict").WithHTTP(fiber.StatusConflict)
)

var (
	GenInternalError = New(CodeGenInternalError, "Internal server error").WithHTTP(fiber.StatusInternalServerError)
)
