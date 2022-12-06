package braid

import (
	"context"
	"github.com/gofiber/fiber/v2"
)

const (
	UserContextFiberCtxKey = "__fiberCtx"
)

func CtxToUserContext(c *fiber.Ctx) {
	c.SetUserContext(context.WithValue(c.UserContext(), UserContextFiberCtxKey, c))
}

func CtxFromUserContext(ctx context.Context) *fiber.Ctx {
	if ctx == nil {
		return nil
	}

	if c, ok := ctx.Value(UserContextFiberCtxKey).(*fiber.Ctx); ok {
		return c
	}

	return nil
}

func CtxInjectionMiddleware(c *fiber.Ctx) error {
	CtxToUserContext(c)
	return c.Next()
}
