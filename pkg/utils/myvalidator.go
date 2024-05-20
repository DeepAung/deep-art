package utils

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

var v = validator.New()

func Validate(s interface{}) error {
	err := v.Struct(s)
	if err == nil {
		return nil
	}

	switch err := err.(type) {
	case validator.ValidationErrors:
		msg := ""
		for _, e := range err {
			msg += fmt.Sprintf("[%s]: '%v' | Needs to implement '%s'\n", e.Field(), e.Value(), e.Tag())
		}
		return errors.New(msg)

	default:
		return err
	}
}
