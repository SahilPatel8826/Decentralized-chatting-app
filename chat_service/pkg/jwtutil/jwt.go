package jwtutil

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var AccessSecret = []byte(
	os.Getenv("JWT_ACCESS_SECRET"),
)

type Claims struct {
	UserID string `json:"user_id"`

	Email string `json:"email"`

	Role string `json:"role"`

	jwt.RegisteredClaims
}

func ValidateAccessToken(
	tokenString string,
) (*Claims, error) {

	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(
			token *jwt.Token,
		) (interface{}, error) {

			return AccessSecret, nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok :=
		token.Claims.(*Claims)

	if !ok || !token.Valid {

		return nil,
			jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
