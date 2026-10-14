package jwtutil

import (
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// Read at call time (not package init) so values loaded
// from .env in main() are visible.
func accessSecret() []byte {
	return []byte(os.Getenv("JWT_ACCESS_SECRET"))
}

type Claims struct {
	UserID string `json:"user_id"`

	Email string `json:"email"`

	Role string `json:"role"`

	jwt.RegisteredClaims
}

func ValidateAccessToken(
	tokenString string,
) (*Claims, error) {

	secret := accessSecret()
	if len(secret) == 0 {
		return nil, errors.New("jwt secret not configured")
	}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		},
		jwt.WithValidMethods([]string{"HS256"}),
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)

	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
