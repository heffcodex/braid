package braid

import (
	"context"
	"github.com/gofiber/fiber/v2"
)

const CtxKeyFiber = "__fiber"

func SetContextValue[T any](c *fiber.Ctx, key string, value T) {
	c.SetUserContext(context.WithValue(c.UserContext(), key, value))
}

func GetContextValue[T any](c *fiber.Ctx, key string) (value T, ok bool) {
	v, ok := c.UserContext().Value(key).(T)
	if !ok || v == nil {
		return *new(T), false
	}

	return v, true
}

func FiberToInnerContext(c *fiber.Ctx) {
	SetContextValue(c, CtxKeyFiber, c)
}

func FiberFromStdContext(ctx context.Context) (c *fiber.Ctx, ok bool) {
	c, ok = ctx.Value(CtxKeyFiber).(*fiber.Ctx)
	return
}
