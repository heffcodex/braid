package validation

import "github.com/go-playground/validator/v10"

type Rule struct {
	Tag            string
	Fn             validator.Func
	CallEvenIfNull bool
}
