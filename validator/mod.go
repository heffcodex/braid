package validator

import "github.com/go-playground/mold/v4"

type Mod struct {
	Tag string
	Fn  mold.Func
}
