package braid

import (
	gojson "github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

const (
	HeaderContentDescription = "Content-Description"
)

type ResponseError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Data    any       `json:"data,omitempty"`

	response *Response
	internal error
}

func (re *ResponseError) Error() string {
	return re.Message
}

func (re *ResponseError) IsInternal() bool {
	return re.Code == ErrorCodeInternal
}

func (re *ResponseError) InternalError() error {
	return re.internal
}

func (re *ResponseError) Response() *Response {
	return re.response
}

type Response struct {
	Meta  any            `json:"meta,omitempty"`
	Data  any            `json:"data,omitempty"`
	Error *ResponseError `json:"err,omitempty"`

	c          *fiber.Ctx
	statusCode int
}

func NewResponse(c *fiber.Ctx) *Response {
	return &Response{c: c, statusCode: c.Response().StatusCode()}
}

func (r *Response) SendEmpty() error {
	r.c.Status(r.statusCode)
	return nil
}

func (r *Response) SendRawData() error {
	r.c.Status(r.statusCode)
	return r.c.Send(r.Data.([]byte))
}

func (r *Response) SendJSON() error {
	raw, err := gojson.MarshalContext(r.c.UserContext(), r)
	if err != nil {
		return err
	}

	r.c.Status(r.statusCode)
	r.c.Response().SetBodyRaw(raw)
	r.c.Response().Header.SetContentType(fiber.MIMEApplicationJSON)

	return nil
}

func (r *Response) SetInternalError(err error) *Response {
	return r.setError(fiber.StatusInternalServerError, ErrorCodeInternal, err, nil)
}

func (r *Response) SetError(statusCode int, eCode ErrorCode, data ...any) *Response {
	if len(data) > 1 {
		panic("too many arguments")
	}

	_data := any(nil)
	if len(data) > 0 {
		_data = data[0]
	}

	return r.setError(statusCode, eCode, nil, _data)
}

func (r *Response) setError(statusCode int, eCode ErrorCode, internal error, data any) *Response {
	r.statusCode = statusCode
	r.Meta = nil
	r.Data = nil
	r.Error = &ResponseError{
		Code:     eCode,
		Message:  eCode.GetMessage(),
		Data:     data,
		response: r,
		internal: internal,
	}

	return r
}

func (r *Response) SetData(httpCode int, data any, meta ...any) *Response {
	if len(meta) > 1 {
		panic("too many arguments")
	}

	_meta := any(nil)
	if len(meta) > 0 {
		_meta = meta[0]
	}

	r.statusCode = httpCode
	r.Meta = _meta
	r.Data = data
	r.Error = nil

	return r
}

// 2xx

func ResponseOK(c *fiber.Ctx, data any, meta ...any) error {
	return NewResponse(c).SetData(fiber.StatusOK, data, meta...).SendJSON()
}

func ResponseFileTransfer(c *fiber.Ctx, filename string, data []byte) error {
	c.Set(HeaderContentDescription, "File Transfer")
	c.Attachment(filename)

	return NewResponse(c).SetData(fiber.StatusOK, data).SendRawData()
}

func ResponseCreated(c *fiber.Ctx, data any, meta ...any) error {
	return NewResponse(c).SetData(fiber.StatusCreated, data, meta...).SendJSON()
}

func ResponseNoContent(c *fiber.Ctx) error {
	return NewResponse(c).SetData(fiber.StatusNoContent, nil).SendEmpty()
}

// All functions below with E prefix return not the error of sending the response to the client,
// but the error-wrapped ResponseError object to be handled by the server-wide error handler:

// 4xx

func EResponseBadRequest(c *fiber.Ctx, eCode ErrorCode, data ...any) error {
	return NewResponse(c).SetError(fiber.StatusBadRequest, eCode, data...).Error
}

func EResponseUnauthorized(c *fiber.Ctx, data ...any) error {
	return NewResponse(c).SetError(fiber.StatusUnauthorized, ErrorCodeNone, data...).Error
}

func EResponseForbidden(c *fiber.Ctx, data ...any) error {
	return NewResponse(c).SetError(fiber.StatusForbidden, ErrorCodeNone, data...).Error
}

func EResponseNotFound(c *fiber.Ctx, data ...any) error {
	return NewResponse(c).SetError(fiber.StatusNotFound, ErrorCodeNone, data...).Error
}

// 5xx

func EResponseInternalError(c *fiber.Ctx, internalErr error) error {
	return NewResponse(c).SetInternalError(internalErr).Error
}
