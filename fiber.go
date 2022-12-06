package braid

import (
	gojson "github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

type OnErrorHandler func(c *fiber.Ctx, err error)

func ConfigDefault() fiber.Config {
	return fiber.Config{
		StrictRouting:                true,
		CaseSensitive:                true,
		ErrorHandler:                 ErrorHandlerDefault(nil),
		StreamRequestBody:            true,
		DisablePreParseMultipartForm: true,
		JSONEncoder:                  gojson.Marshal,
		JSONDecoder:                  gojson.Unmarshal,
		EnablePrintRoutes:            true,
	}
}

func ErrorHandlerDefault(onErrorHandler OnErrorHandler) fiber.ErrorHandler {
	if onErrorHandler == nil {
		onErrorHandler = func(*fiber.Ctx, error) {}
	}

	return func(c *fiber.Ctx, err error) error {
		var fe *fiber.Error

		if errors.As(err, &fe) { // catch fiber error and wrap it into braid.Response, hiding 5xx error details from client
			message := fe.Message

			if fe.Code >= 500 && fe.Code <= 599 {
				message = ErrorCodeInternal.GetMessage()
				onErrorHandler(c, errors.Wrapf(err, "fiber.Error [%d]", fe.Code))
			}

			return NewResponse(c).SetError(fe.Code, NewErrorCode(ECICustom, message)).JSON()
		}

		var re *ResponseError

		if errors.As(err, &re) {
			if re.IsInternal() {
				onErrorHandler(c, errors.Wrap(re.InternalError(), "braid.ResponseError"))
			}

			return re.Response().JSON()
		}

		onErrorHandler(c, errors.Wrap(err, "unknown error"))
		c.Status(fiber.StatusInternalServerError).Response().ResetBody()

		return nil
	}
}
