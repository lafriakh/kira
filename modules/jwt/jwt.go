package jwt

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lafriakh/kira"
)

var (
	ErrInvalidToken  = kira.E("Invalid token", http.StatusUnauthorized)
	ErrInvalidClaims = kira.E("Invalid claims", http.StatusUnauthorized)
)

const ContextClaimsKey = "jwt_claims"

type JWT[T jwt.Claims] struct {
	secret string
}

func New[T jwt.Claims](secret string) *JWT[T] {
	return &JWT[T]{secret: secret}
}

func CreateToken(secret string, claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func ValidateToken[T jwt.Claims](secret, tokenStr string) (*jwt.Token, bool) {
	claims := new(T)

	dst, ok := any(claims).(jwt.Claims)
	if !ok {
		return nil, false
	}

	token, err := jwt.ParseWithClaims(tokenStr, dst, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return nil, false
	}

	return token, true
}

func ExtractClaims[T jwt.Claims](ctx *kira.Context) *T {
	data := ctx.GetData(ContextClaimsKey)

	if claims, ok := data.(T); ok {
		return &claims
	}

	if claimsPtr, ok := data.(*T); ok {
		if claimsPtr != nil {
			return claimsPtr
		}
	}

	return nil
}

func (j *JWT[T]) Middleware(ctx *kira.Context, next kira.HandlerFunc) error {
	var tokenString string
	var err error

	lookup := strings.Split(ctx.Config().GetString("jwt.lookup", "header:Authorization"), ":")

	switch lookup[0] {
	case "header":
		tokenString, err = j.fromHeader(ctx, lookup[1])
		if err != nil {
			return err
		}
	case "cookie":
		tokenString, err = j.fromCookie(ctx, lookup[1])
		if err != nil {
			return err
		}
	}

	if token, ok := ValidateToken[T](j.secret, tokenString); ok {
		ctx.SetData(ContextClaimsKey, token.Claims)

		return next(ctx)
	} else {
		return ErrInvalidToken
	}
}
