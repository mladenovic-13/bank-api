package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validation struct {
	*validator.Validate
}

func NewCustomValidator(v *validator.Validate) *Validation {
	return &Validation{v}
}

var validate = validator.New()

var Validate = NewCustomValidator(validate)

func (v *Validation) CustomRequest(r interface{}) map[string]string {
	errs := v.Validate.Struct(r)

	if errs != nil {
		msgs := make(map[string]string)

		for _, err := range errs.(validator.ValidationErrors) {
			msg := fmt.Sprintf("%s %s", err.Tag(), err.Param())
			msgs[strings.ToLower(err.Field())] = msg
		}

		return msgs
	}

	return nil
}
