package middleware

import (
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"
)

func RecoverConfigDefault(log *zap.Logger) recover.Config {
	return recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e any) {
			log.Error("panic recovered", zap.Any("error", e), zap.ByteString("stack", debug.Stack()))
		},
	}
}

func Recover(config recover.Config) fiber.Handler {
	return recover.New(config)
}
