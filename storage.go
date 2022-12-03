package braid

import (
	"context"
)

type ContextProvider func() context.Context
