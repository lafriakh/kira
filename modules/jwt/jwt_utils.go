package jwt

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/lafriakh/kira"
)

// Extract the token from the request header.
func (j *JWT[T]) fromHeader(ctx *kira.Context, key string) (string, error) {
	authHeader := ctx.Request().Header.Get(key)
	if authHeader == "" {
		return "", ErrInvalidToken // No token present
	}

	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", fmt.Errorf("Authorization header format is not valid")
	}

	return authHeaderParts[1], nil
}

// Extract the token from the request cookie.
func (j *JWT[T]) fromCookie(ctx *kira.Context, key string) (string, error) {
	c, err := ctx.Request().Cookie(key)
	if err != nil {
		if err == http.ErrNoCookie {
			return "", ErrInvalidToken
		}
		return "", err
	}

	return c.Value, nil
}
