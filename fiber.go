package braid

import (
	gojson "github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

type OnInternalError func(c *fiber.Ctx, err error)

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

func ErrorHandlerDefault(onInternalError OnInternalError) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		var fe *fiber.Error

		if errors.As(err, &fe) { // catch fiber error and wrap it into braid.Response, hiding some details from client
			switch fe.Code {
			case fiber.StatusBadRequest:
				err = EResponseBadRequest(c, NewErrorCode(ECICustom, fe.Message))
			case fiber.StatusUnauthorized:
				err = EResponseUnauthorized(c)
			case fiber.StatusForbidden:
				err = EResponseForbidden(c)
			case fiber.StatusNotFound:
				err = EResponseNotFound(c)
			default:
				if fe.Code >= 500 && fe.Code <= 599 {
					err = EResponseInternalError(c, err)
				} else {
					err = NewResponse(c).SetError(fe.Code, NewErrorCode(ECICustom, fe.Message)).Error
				}
			}
		}

		re := new(ResponseError)

		if !errors.As(err, re) {
			re = EResponseInternalError(c, err).(*ResponseError)
		}

		if re.IsInternal() && onInternalError != nil {
			onInternalError(c, re.InternalError())
		}

		return re.Response().SendJSON()
	}
}
