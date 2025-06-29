package kira

import (
	"net/http"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	KiraError

	Errors map[string]any `json:"errors"`
}

var validate = validator.New(validator.WithRequiredStructEnabled())

func Validate[T any](c *Context) T {
	result := new(T)
	if err := c.DecodeJSON(result); err != nil {
		c.Error(E("can't decode json", err, StatusCode(http.StatusBadRequest)))
	}

	if err := validate.Struct(result); err != nil {
		c.Error(E("validation failed", err, StatusCode(http.StatusBadRequest)))
	}
	
	return *result
}
