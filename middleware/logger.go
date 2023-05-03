package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.uber.org/zap"
)

func LoggerConfigDefault(log *zap.Logger) logger.Config {
	return logger.Config{
		Format: fmt.Sprintf(
			"${%s} - ${%s} ${%s} ${%s}\n",
			logger.TagStatus,
			logger.TagLatency,
			logger.TagMethod,
			logger.TagPath,
		),
		Output: zap.NewStdLog(log).Writer(),
	}
}

func Logger(config logger.Config) fiber.Handler {
	return logger.New(config)
}
