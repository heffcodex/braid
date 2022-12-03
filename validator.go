package braid

import (
	"context"
	"reflect"
	"strings"

	"github.com/go-playground/mold/v4"
	"github.com/go-playground/mold/v4/modifiers"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

const LocalValidator = "validator"

type validationErrorData struct {
	Tag     string `json:"tag"`
	Param   string `json:"param"`
	Message string `json:"message"`
}

type Validator struct {
	tf *mold.Transformer
	v  *validator.Validate
}

func V(c *fiber.Ctx) *Validator {
	return getValidator(c)
}

func getValidator(c *fiber.Ctx) *Validator {
	return c.Locals(LocalValidator).(*Validator)
}

func setValidator(c *fiber.Ctx, v *Validator) {
	c.Locals(LocalValidator, v)
}

func NewValidator() *Validator {
	v := validator.New()

	v.RegisterTagNameFunc(
		func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		},
	)

	return &Validator{
		tf: modifiers.New(),
		v:  v,
	}
}

func (v *Validator) RegisterValidation(tag string, fn validator.Func, callValidationEvenIfNull ...bool) error {
	return v.v.RegisterValidation(tag, fn, callValidationEvenIfNull...)
}

func (v *Validator) RegisterModifier(tag string, fn mold.Func) {
	v.tf.Register(tag, fn)
}

func (v *Validator) BindAndValidate(c *fiber.Ctx, form any) error {
	if err := c.BodyParser(form); err != nil {
		return EResponseBadRequest(c, ErrorCodeInvalidPayload)
	}

	vErr, err := v.validate(c.Context(), form)
	if err != nil {
		return EResponseInternalError(c, errors.Wrap(err, "cannot validate form"))
	} else if len(vErr) > 0 {
		return EResponseBadRequest(c, ErrorCodeValidationFail, vErr)
	}

	return nil
}

func (v *Validator) Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		setValidator(c, v)
		return c.Next()
	}
}

func (v *Validator) validate(ctx context.Context, form any) (map[string]validationErrorData, error) {
	err := v.tf.Struct(ctx, form)
	if err != nil {
		return nil, err
	}

	err = v.v.Struct(form)
	if vErr, ok := err.(validator.ValidationErrors); ok {
		return v.transformValidationErrors(vErr), nil
	}

	return nil, err
}

func (v *Validator) transformValidationErrors(errors validator.ValidationErrors) map[string]validationErrorData {
	m := make(map[string]validationErrorData)

	for _, err := range errors {
		m[err.Field()] = validationErrorData{
			Tag:     err.Tag(),
			Param:   err.Param(),
			Message: err.Error(),
		}
	}

	return m
}
