package status

import "github.com/gofiber/fiber/v2"

const (
	CodeGenOK        Code = 200
	CodeGenCreated   Code = 201
	CodeGenNoContent Code = 204

	CodeGenBadRequest     Code = 400
	CodeCSRFTokenMismatch Code = 400001
	CodeValidationFail    Code = 400002
	CodeInvalidQuery      Code = 400003
	CodeInvalidPayload    Code = 400004

	CodeGenUnauthorized Code = 401
	CodeGenForbidden    Code = 403
	CodeGenNotFound     Code = 404

	CodeGenInternalError Code = 500
)

var (
	GenOK        = New(CodeGenOK).SetHTTP(fiber.StatusOK)
	GenCreated   = New(CodeGenCreated).SetHTTP(fiber.StatusCreated)
	GenNoContent = New(CodeGenNoContent).SetHTTP(fiber.StatusNoContent)
)

var (
	GenBadRequest     = New(CodeGenBadRequest, "Bad request").SetHTTP(fiber.StatusBadRequest)
	CSRFTokenMismatch = New(CodeCSRFTokenMismatch, "CSRF token mismatch").AttachTo(GenBadRequest)
	ValidationFail    = New(CodeValidationFail, "Validation failed").AttachTo(GenBadRequest)
	InvalidQuery      = New(CodeInvalidQuery, "Invalid query").AttachTo(GenBadRequest)
	InvalidPayload    = New(CodeInvalidPayload, "Invalid payload").AttachTo(GenBadRequest)
)

var (
	GenUnauthorized = New(CodeGenUnauthorized, "Unauthorized").SetHTTP(fiber.StatusUnauthorized)
	GenForbidden    = New(CodeGenForbidden, "Forbidden").SetHTTP(fiber.StatusForbidden)
	GenNotFound     = New(CodeGenNotFound, "Not found").SetHTTP(fiber.StatusNotFound)
)

var (
	GenInternalError = New(CodeGenInternalError, "Internal server error").SetHTTP(fiber.StatusInternalServerError)
)
