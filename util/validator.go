package util

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"strings"
)

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

func init() {
	validate = validator.New()
}

type Validator func(input string) error

func UrlValidator(input string) error {
	return validate.Var(input, "url")
}

func RequiredValidator(input string) error {
	if strings.TrimSpace(input) != "" {
		return validate.Var(input, "required")
	}
	return errors.New("required")
}

func YesNoValidator(input string) error {
	return validate.Var(input, "oneof=Y y N n")
}
