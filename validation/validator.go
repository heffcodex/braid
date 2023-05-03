package validation

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/go-playground/mold/v4"
	"github.com/go-playground/mold/v4/modifiers"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/heffcodex/braid/response"
	"github.com/heffcodex/braid/status"
)

type validationErrorData struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Param string `json:"param"`
}

type Validator struct {
	tf *mold.Transformer
	v  *validator.Validate
}

func New(opts ...Option) *Validator {
	vld := validator.New()

	vld.RegisterTagNameFunc(
		func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		},
	)

	v := &Validator{
		tf: modifiers.New(),
		v:  vld,
	}

	for _, opt := range opts {
		opt(v)
	}

	return v
}

func (v *Validator) BindAndValidate(c *fiber.Ctx, form any) error {
	if err := c.BodyParser(form); err != nil {
		return response.EBadRequest(status.InvalidPayload)
	}

	vErr, err := v.validate(c.Context(), form)
	if err != nil {
		return response.EInternal(fmt.Errorf("validate form"))
	} else if len(vErr) > 0 {
		return response.EBadRequest(status.ValidationFail, vErr)
	}

	return nil
}

func (v *Validator) registerRules(rules ...Rule) error {
	for _, rule := range rules {
		if err := v.v.RegisterValidation(rule.Tag, rule.Fn, rule.CallEvenIfNull); err != nil {
			return fmt.Errorf("%s: %w", rule.Tag, err)
		}
	}

	return nil
}

func (v *Validator) registerModifiers(mods ...Mod) {
	for _, mod := range mods {
		v.tf.Register(mod.Tag, mod.Fn)
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
