package jwtutil

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	AccessSecret  = []byte(os.Getenv("JWT_ACCESS_SECRET"))
	RefreshSecret = []byte(os.Getenv("JWT_REFRESH_SECRET"))
)

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`

	jwt.RegisteredClaims
}

func GenerateAccessToken(
	userID string,
	email string,
	role string,
) (string, error) {

	claims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(
					15 * time.Minute,
				),
			),

			IssuedAt: jwt.NewNumericDate(
				time.Now(),
			),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString(
		AccessSecret,
	)
}

func GenerateRefreshToken(
	userID string,
) (string, error) {

	claims := jwt.RegisteredClaims{

		Subject: userID,

		ExpiresAt: jwt.NewNumericDate(
			time.Now().Add(
				30 * 24 * time.Hour,
			),
		),

		IssuedAt: jwt.NewNumericDate(
			time.Now(),
		),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString(
		RefreshSecret,
	)
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

	claims, ok := token.Claims.(*Claims)

	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}

func ValidateRefreshToken(
	tokenString string,
) (string, error) {

	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{},
		func(
			token *jwt.Token,
		) (interface{}, error) {

			return RefreshSecret, nil
		},
	)

	if err != nil {
		return "", err
	}

	claims, ok :=
		token.Claims.(*jwt.RegisteredClaims)

	if !ok || !token.Valid {
		return "", jwt.ErrTokenInvalidClaims
	}

	return claims.Subject, nil
}
