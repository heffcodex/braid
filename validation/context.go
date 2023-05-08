package validation

import (
	"github.com/gofiber/fiber/v2"

	"github.com/heffcodex/braid/vars"
)

type ContextValidator struct {
	c *fiber.Ctx
	v *Validator
}

func V(c *fiber.Ctx, set ...*Validator) *ContextValidator {
	var v *Validator

	if len(set) > 1 {
		panic("too many arguments")
	} else if len(set) == 1 {
		setValidator(c, set[0])
		v = set[0]
	}

	v = getValidator(c)

	return &ContextValidator{c: c, v: v}
}

func getValidator(c *fiber.Ctx) *Validator {
	return c.Locals(vars.LocalValidator).(*Validator)
}

func setValidator(c *fiber.Ctx, v *Validator) {
	c.Locals(vars.LocalValidator, v)
}

func (v *ContextValidator) Form(form any) error {
	return v.v.BindAndValidateForm(v.c, form)
}

func (v *ContextValidator) Query(query any) error {
	return v.v.BindAndValidateQuery(v.c, query)
}
