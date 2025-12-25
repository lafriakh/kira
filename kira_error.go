package kira

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/go-playground/validator/v10"
)

type StatusCode int
type Path string
type FieldError struct {
	// Field is the name of the struct field that failed validation.
	Field string `json:"field"`
	// Tag is the name of the validation rule that failed (e.g., "required", "min").
	Tag string `json:"tag"`
	// Value is the actual invalid value that was provided.
	Value any `json:"value"`
	// Param is the validation rule's parameter, if any (e.g., "2" for "min=2").
	Param string `json:"param,omitempty"`
}

type Error struct {
	Status      StatusCode      `json:"status,omitempty"`
	Path        Path            `json:"path,omitempty"`
	Message     string          `json:"message"`
	FieldErrors []FieldError    `json:"errors,omitempty"`
	Timestamp   time.Time       `json:"timestamp,omitzero"`
	Frames      []runtime.Frame `json:"frames,omitempty"`

	Err error `json:"-"`
}

func AsError(err error) (*Error, bool) {
	var target *Error

	if errors.As(err, &target) {
		return target, true
	}

	return nil, false
}

func E(args ...any) error {
	if len(args) == 0 {
		panic("call to kira.E with no arguments")
	}

	e := &Error{
		Timestamp: time.Now().UTC(),
		Status:    500,
	}

	for _, arg := range args {
		switch arg := arg.(type) {
		case string:
			e.Message = arg
			e.Err = errors.New(arg)
		case StatusCode:
			e.Status = arg
		case Path:
			e.Path = arg
		case *Error:
			copy := *arg
			e.Err = &copy
		case validator.ValidationErrors:
			if e.FieldErrors == nil {
				var fieldErrors []FieldError
				for _, fe := range arg {
					err := FieldError{
						Field: fe.Field(),
						Tag:   fe.Tag(),
						Value: fe.Value(),
						Param: fe.Param(),
					}
					fieldErrors = append(fieldErrors, err)
				}
				e.FieldErrors = fieldErrors
			}
		case error:
			e.Err = arg
		default:
			return fmt.Errorf("unknown type %T, value %v in error call", arg, arg)
		}
	}

	return e
}

func (e *Error) Error() string {
	b := new(bytes.Buffer)
	
	if e.Status != 0 {
		pad(b, ": ")
		fmt.Fprint(b, http.StatusText(int(e.Status)))
	}
	if e.Path != "" {
		pad(b, ": ")
		b.WriteString(string(e.Path))
	}

	if e.Err != nil {
		pad(b, ": ")
		b.WriteString(string(e.Err.Error()))
	}

	if b.Len() == 0 {
		return "no error"
	}

	return b.String()
}

// pad appends str to the buffer if the buffer already has some data.
func pad(b *bytes.Buffer, str string) {
	if b.Len() == 0 {
		return
	}
	b.WriteString(str)
}
