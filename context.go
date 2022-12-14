package braid

import (
	"context"
	"github.com/gofiber/fiber/v2"
)

func SetContextValue[T any](c *fiber.Ctx, key string, value T) {
	c.SetUserContext(context.WithValue(c.UserContext(), key, value))
}

func GetContextValue[T any](c *fiber.Ctx, key string) (value T, ok bool) {
	v, ok := c.UserContext().Value(key).(T)
	if !ok {
		return *new(T), false
	}

	return v, true
}
