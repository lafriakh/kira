package kira

import (
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func Validate[T any](c *Context) (*T, error) {
	result := new(T)
	if err := c.DecodeJSON(result); err != nil {
		return nil, E("cannot decode json", err, StatusCode(http.StatusBadRequest))
	}

	if err := validate.Struct(result); err != nil {
		return nil, E("validation failed", err, StatusCode(http.StatusBadRequest))
	}

	return result, nil
}
