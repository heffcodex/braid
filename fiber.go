package braid

import (
	gojson "github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/heffcodex/braid/response"
	"github.com/heffcodex/braid/status"
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
			if fe.Code >= 500 && fe.Code <= 599 {
				err = response.EInternal(err)
			} else {
				err = response.NewJSONError(status.New(status.Code(fe.Code), fe.Message).WithHTTP(fe.Code), nil)
			}
		}

		var je *response.JSONError

		if !errors.As(err, &je) {
			je = response.NewJSONError(status.GenInternalError, err)
		}

		if je.IsInternal() && onInternalError != nil {
			onInternalError(c, je.InternalError())
		}

		return je.Send(c)
	}
}
