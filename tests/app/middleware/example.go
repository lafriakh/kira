package middleware

import "github.com/lafriakh/kira"

// Example
type Example struct{}

func New() *Example {
	return &Example{}
}
func (e *Example) Middleware(ctx *kira.Context, next kira.HandlerFunc) error {
	ctx.WriteString("Before")
	err := next(ctx)
	ctx.WriteString("After")

	return err
}

// Example2
type Example2 struct{}

func New2() *Example2 {
	return &Example2{}
}
func (e *Example2) Middleware(ctx *kira.Context, next kira.HandlerFunc) error {
	ctx.WriteString("Before2")
	err := next(ctx)
	ctx.WriteString("After2")
	return err
}

// ContextData
type ContextData struct{}

func NewContextData() *ContextData {
	return &ContextData{}
}
func (e *ContextData) Middleware(ctx *kira.Context, next kira.HandlerFunc) error {
	// We should see this inside a normal handler that use this middleware.
	ctx.SetData("foo", "bar")

	return next(ctx)
}
