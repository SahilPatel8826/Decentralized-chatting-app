package jwtutil

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const testSecret = "test-secret"

func sign(t *testing.T, method jwt.SigningMethod, secret string, exp time.Time) string {
	t.Helper()
	claims := &Claims{
		UserID: "42",
		Email:  "a@test.com",
		Role:   "user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}
	s, err := jwt.NewWithClaims(method, claims).SignedString([]byte(secret))
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func TestValidToken(t *testing.T) {
	t.Setenv("JWT_ACCESS_SECRET", testSecret)
	tok := sign(t, jwt.SigningMethodHS256, testSecret, time.Now().Add(time.Hour))

	c, err := ValidateAccessToken(tok)
	if err != nil {
		t.Fatalf("expected valid token, got %v", err)
	}
	if c.UserID != "42" || c.Email != "a@test.com" || c.Role != "user" {
		t.Fatalf("unexpected claims: %+v", c)
	}
}

func TestSecretReadAtCallTime(t *testing.T) {
	// Regression: secret used to be captured at package init.
	tok := sign(t, jwt.SigningMethodHS256, "late-secret", time.Now().Add(time.Hour))
	t.Setenv("JWT_ACCESS_SECRET", "late-secret")

	if _, err := ValidateAccessToken(tok); err != nil {
		t.Fatalf("secret set after init should work: %v", err)
	}
}

func TestWrongSecretRejected(t *testing.T) {
	t.Setenv("JWT_ACCESS_SECRET", testSecret)
	tok := sign(t, jwt.SigningMethodHS256, "other-secret", time.Now().Add(time.Hour))

	if _, err := ValidateAccessToken(tok); err == nil {
		t.Fatal("token signed with another secret must be rejected")
	}
}

func TestExpiredRejected(t *testing.T) {
	t.Setenv("JWT_ACCESS_SECRET", testSecret)
	tok := sign(t, jwt.SigningMethodHS256, testSecret, time.Now().Add(-time.Minute))

	if _, err := ValidateAccessToken(tok); err == nil {
		t.Fatal("expired token must be rejected")
	}
}

func TestOtherAlgorithmRejected(t *testing.T) {
	t.Setenv("JWT_ACCESS_SECRET", testSecret)
	tok := sign(t, jwt.SigningMethodHS512, testSecret, time.Now().Add(time.Hour))

	if _, err := ValidateAccessToken(tok); err == nil {
		t.Fatal("only HS256 should be accepted")
	}
}

func TestEmptySecretRejected(t *testing.T) {
	t.Setenv("JWT_ACCESS_SECRET", "")
	tok := sign(t, jwt.SigningMethodHS256, "", time.Now().Add(time.Hour))

	if _, err := ValidateAccessToken(tok); err == nil {
		t.Fatal("empty secret must never validate tokens")
	}
}

func TestGarbageRejected(t *testing.T) {
	t.Setenv("JWT_ACCESS_SECRET", testSecret)
	if _, err := ValidateAccessToken("not-a-jwt"); err == nil {
		t.Fatal("garbage must be rejected")
	}
}
