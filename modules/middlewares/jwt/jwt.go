package jwt

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lafriakh/kira"
)

var (
	ErrInvalidToken = kira.E("Invalid token", kira.StatusCode(http.StatusUnauthorized))
)

type JWT struct {
	secret string
}

func New(secret string) *JWT {
	return &JWT{secret: secret}
}

// CreateToken generate JWT token.
func CreateToken(secret string, claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func ValidateToken(secret, s string) bool {
	token, err := jwt.Parse(s, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
	if err != nil {
		return false
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return true
	}

	return false
}

// Middleware handler
func (j *JWT) Middleware(ctx *kira.Context, next kira.HandlerFunc) error {
	// The token string that will be validated.
	var tokenString string
	var err error

	// From where we should grap the token.
	lookup := strings.Split(ctx.Config().GetString("jwt.lookup", "header:Authorization"), ":")

	switch lookup[0] {
	case "header": // From header
		tokenString, err = j.fromHeader(ctx, lookup[1])
		if err != nil {
			return err
		}
	case "cookie": // From cookie
		tokenString, err = j.fromCookie(ctx, lookup[1])
		if err != nil {
			return err
		}
	}

	// Validate if the request has a valide JWT Token.
	if ValidateToken(j.secret, tokenString) {
		return next(ctx)
	} else {
		return ErrInvalidToken
	}
}

// Extract the token from the request header.
func (j *JWT) fromHeader(ctx *kira.Context, key string) (string, error) {
	authHeader := ctx.Request().Header.Get(key)
	if authHeader == "" {
		return "", nil // No error, just no token
	}

	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", fmt.Errorf("Authorization header format is not valid")
	}

	return authHeaderParts[1], nil
}

// Extract the token from the request cookie.
func (j *JWT) fromCookie(ctx *kira.Context, key string) (string, error) {
	c, err := ctx.Request().Cookie(key)
	if err != nil {
		return "", err
	}

	return c.Value, nil
}
