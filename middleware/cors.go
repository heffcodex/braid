package middleware

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"golang.org/x/exp/slices"
)

const (
	CORSAllowAny       = "!*"
	CORSAllowLocalhost = "!localhost"
)

func CORSConfigDefault(allowOrigins ...string) cors.Config {
	return cors.Config{
		AllowOrigins: strings.Join(allowOrigins, ","),
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodHead,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodPatch,
		}, ","),
		AllowHeaders: strings.Join([]string{
			csrf.HeaderName,
			fiber.HeaderAcceptEncoding,
			fiber.HeaderContentLength,
			fiber.HeaderContentType,
			fiber.HeaderOrigin,
			fiber.HeaderXRequestedWith,
		}, ","),
		ExposeHeaders: strings.Join([]string{
			csrf.HeaderName,
			fiber.HeaderXRequestID,
		}, ","),
		AllowCredentials: true,
		MaxAge:           int((24 * time.Hour).Seconds()),
	}
}

func CORS(config cors.Config) fiber.Handler {
	origins := strings.Split(config.AllowOrigins, ",")
	allowFns := make([]func(origin string) bool, 0)

	if slices.Contains(origins, CORSAllowAny) {
		allowFns = append(allowFns, func(origin string) bool {
			return true
		})
	}

	if slices.Contains(origins, CORSAllowLocalhost) {
		allowFns = append(allowFns, func(origin string) bool {
			//goland:noinspection HttpUrlsUsage
			origin = strings.TrimPrefix(origin, "http://")
			origin = strings.TrimPrefix(origin, "https://")

			return origin == "localhost" || strings.HasPrefix(origin, "localhost:")
		})
	}

	if config.AllowOriginsFunc != nil {
		allowFns = append(allowFns, config.AllowOriginsFunc)
	}

	if len(allowFns) > 0 {
		config.AllowOriginsFunc = func(origin string) bool {
			for _, fn := range allowFns {
				if fn(origin) {
					return true
				}
			}

			return false
		}
	}

	return cors.New(config)
}
