package braid

import (
	"github.com/gofiber/fiber/v2"
)

func CtxValue[T any](c *fiber.Ctx, key string, value ...T) (v T, ok bool) {
	if len(value) > 1 {
		panic("too many arguments")
	} else if len(value) == 1 {
		return c.Locals(key, value[0]).(T), true
	}

	v, ok = c.Locals(key).(T)
	return
}

func MustCtxValue[T any](c *fiber.Ctx, key string, value ...T) T {
	v, ok := CtxValue[T](c, key, value...)
	if !ok {
		panic("`" + key + "` not found in context")
	}

	return v
}
