package braid

import (
	"context"
	"github.com/heffcodex/braid/response"
	"github.com/heffcodex/braid/status"
	"reflect"
	"strings"
	"unicode"

	"github.com/go-playground/mold/v4"
	"github.com/go-playground/mold/v4/modifiers"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

const LocalValidator = "validator"

type validationErrorData struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Param string `json:"param"`
}

type Validator struct {
	tf *mold.Transformer
	v  *validator.Validate
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
		return response.EBadRequest(status.InvalidPayload)
	}

	vErr, err := v.validate(c.Context(), form)
	if err != nil {
		return response.EInternal(errors.Wrap(err, "cannot validate form"))
	} else if len(vErr) > 0 {
		return response.EBadRequest(status.ValidationFail, vErr)
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
		nsArrDirty := strings.Split(err.Namespace(), ".")
		if len(nsArrDirty) > 1 { // skip struct name
			nsArrDirty = nsArrDirty[1:]
		}

		nsArr := make([]string, 0, len(nsArrDirty))
		for _, part := range nsArrDirty {
			if len(part) == 0 { // skip empty parts
				continue
			}

			if unicode.IsUpper([]rune(part)[0]) { // skip upper case parts as they are possible embedded structs
				continue
			}

			nsArr = append(nsArr, part)
		}

		ns := strings.Join(nsArr, ".")
		if ns == "" {
			ns = err.Field()
		}

		m[ns] = validationErrorData{
			Field: err.Field(),
			Tag:   err.Tag(),
			Param: err.Param(),
		}
	}

	return m
}
