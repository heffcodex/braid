package braid

import "github.com/gofiber/fiber/v2"

type ValidatorContext struct {
	v *Validator
	c *fiber.Ctx
}

func V(c *fiber.Ctx) *ValidatorContext {
	return &ValidatorContext{
		v: getValidator(c),
		c: c,
	}
}

func (v *ValidatorContext) BindAndValidate(form any) error {
	return v.v.BindAndValidate(v.c, form)
}
