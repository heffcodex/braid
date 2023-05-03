package validator

import (
	"context"
	"testing"

	"github.com/go-playground/mold/v4"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestValidator_WithRules(t *testing.T) {
	v := New(WithRules(
		Rule{
			Tag: "testTrue",
			Fn:  func(fl validator.FieldLevel) bool { return true },
		},
		Rule{
			Tag: "testFalse",
			Fn:  func(fl validator.FieldLevel) bool { return false },
		},
	))

	type E struct {
		Inner struct {
			FieldFalse string `json:"field_false" validate:"testFalse"`
		} `json:"inner"`
	}

	type form struct {
		E
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

func TestValidator_WithMods(t *testing.T) {
	v := New(WithMods(
		Mod{
			Tag: "test",
			Fn: func(ctx context.Context, fl mold.FieldLevel) error {
				fl.Field().SetString("test")
				return nil
			},
		},
	))

	type form struct {
		Field string `json:"field" mod:"test"`
	}

	f := new(form)

	vErr, err := v.validate(context.Background(), f)
	require.NoError(t, err)
	require.Empty(t, vErr)
	require.Equal(t, "test", f.Field)
}
