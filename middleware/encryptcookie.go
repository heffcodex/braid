package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
)

func EncryptCookieConfigDefault(key string) encryptcookie.Config {
	return encryptcookie.Config{
		Key: key,
	}
}

func EncryptCookie(config encryptcookie.Config) fiber.Handler {
	return encryptcookie.New(config)
}
