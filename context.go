package braid

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

func Set[T any](c *fiber.Ctx, key string, value T) {
	c.SetUserContext(context.WithValue(c.UserContext(), key, value))
}

func Get[T any](c *fiber.Ctx, key string) (value T, ok bool) {
	v, ok := c.UserContext().Value(key).(T)
	if !ok {
		return *new(T), false
	}

	return v, true
}

func MustGet[T any](c *fiber.Ctx, key string) T {
	v, ok := Get[T](c, key)
	if !ok {
		panic("`" + key + "` not found in context")
	}

	return v
}
