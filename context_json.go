package kira

import (
	"encoding/json"
	"io"
	"net/http"
)

var ErrEmptyBody = E("empty body", http.StatusBadRequest)

// JSON - Send response as json.
func (c *Context) JSON(data any, code ...int) error {
	c.Response().Header().Set("Content-Type", "application/json")

	// Status statusCode
	if len(code) > 0 {
		c.Status(code[0])
	}

	// Encode data
	return json.NewEncoder(c.Response()).Encode(data)
}

// WantsJSON - validate if the request wants a json response.
func (c *Context) WantsJSON() bool {
	return c.Request().Header.Get("Accept") == "application/json"
}

// DecodeJSON - convert json from request body to interface.
func (c *Context) DecodeJSON(dst any) error {
	err := json.NewDecoder(c.Request().Body).Decode(dst)
	if err != nil {
		if err == io.EOF {
			return ErrEmptyBody
		}
		return err
	}

	return nil
}
