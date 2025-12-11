package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Validator struct {
	v *validator.Validate
}

func NewValidator() *Validator {
	v := validator.New()

	// Use JSON tag names in errors
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Custom UUID validation
	v.RegisterValidation("uuid", validateUUID)

	return &Validator{v: v}
}

func (v *Validator) ValidateStruct(s interface{}) error {
	if err := v.v.Struct(s); err != nil {
		return err
	}
	return nil
}

func validateUUID(fl validator.FieldLevel) bool {
	str := fl.Field().String()
	_, err := uuid.Parse(str)
	return err == nil
}

// ParseValidationErrors converts validator errors into a simple field->message map
func ParseValidationErrors(err error) map[string]string {
	out := make(map[string]string)

	if err == nil {
		return out
	}

	if verrs, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range verrs {
			field := fe.Field() // already uses JSON tag from TagNameFunc
			switch fe.Tag() {
			case "required":
				out[field] = "is required"
			case "min":
				out[field] = fmt.Sprintf("must be at least %s characters", fe.Param())
			case "max":
				out[field] = fmt.Sprintf("must be at most %s characters", fe.Param())
			case "oneof":
				out[field] = fmt.Sprintf("must be one of: %s", fe.Param())
			case "url":
				out[field] = "must be a valid URL"
			case "uuid":
				out[field] = "must be a valid UUID"
			default:
				out[field] = "is invalid"
			}
		}
	} else {
		// Non-typed error
		out["_"] = err.Error()
	}

	return out
}
