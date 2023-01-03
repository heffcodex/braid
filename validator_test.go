package braid

import (
	"context"
	"github.com/go-playground/mold/v4"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidator_RegisterValidation(t *testing.T) {
	v := NewValidator()

	require.NoError(t, v.RegisterValidation("testTrue", func(fl validator.FieldLevel) bool {
		return true
	}))

	require.NoError(t, v.RegisterValidation("testFalse", func(fl validator.FieldLevel) bool {
		return false
	}))

	type form struct {
		Inner struct {
			FieldFalse string `json:"field_false" validate:"testFalse"`
		} `json:"inner"`
		FieldTrue string `json:"field_true" validate:"testTrue"`
	}

	vErr, err := v.validate(context.Background(), new(form))
	require.NoError(t, err)
	require.Equal(t, map[string]validationErrorData{
		"inner.field_false": {
			Field: "field_false",
			Tag:   "testFalse",
			Param: "",
		},
	}, vErr)
}

func TestValidator_RegisterModifier(t *testing.T) {
	v := NewValidator()

	v.RegisterModifier("test", func(ctx context.Context, fl mold.FieldLevel) error {
		fl.Field().SetString("test")
		return nil
	})

	type form struct {
		Field string `json:"field" mod:"test"`
	}

	f := new(form)

	vErr, err := v.validate(context.Background(), f)
	require.NoError(t, err)
	require.Empty(t, vErr)
	require.Equal(t, "test", f.Field)
}
