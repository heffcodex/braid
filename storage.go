package braid

import (
	"context"
)

type ContextProvider func() context.Context

var ctxProviderDefault = func() context.Context { return context.Background() }
