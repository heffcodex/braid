package storage

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

type Storage interface {
	fiber.Storage
}

type ContextProvider func() context.Context

var ctxProviderDefault = func() context.Context { return context.Background() }
