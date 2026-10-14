package jwtutil

import "testing"

func TestAccessTokenRoundTrip(t *testing.T) {
	t.Setenv("JWT_ACCESS_SECRET", "access-secret")
	t.Setenv("JWT_REFRESH_SECRET", "refresh-secret")

	tok, err := GenerateAccessToken("5", "a@test.com", "user")
	if err != nil {
		t.Fatal(err)
	}
	c, err := ValidateAccessToken(tok)
	if err != nil {
		t.Fatalf("valid access token rejected: %v", err)
	}
	if c.UserID != "5" || c.Email != "a@test.com" || c.Role != "user" {
		t.Fatalf("unexpected claims: %+v", c)
	}
}

func TestRefreshTokenRoundTrip(t *testing.T) {
	t.Setenv("JWT_ACCESS_SECRET", "access-secret")
	t.Setenv("JWT_REFRESH_SECRET", "refresh-secret")

	tok, err := GenerateRefreshToken("5")
	if err != nil {
		t.Fatal(err)
	}
	id, err := ValidateRefreshToken(tok)
	if err != nil || id != "5" {
		t.Fatalf("want 5, got %q err=%v", id, err)
	}
}

func TestTokensAreNotInterchangeable(t *testing.T) {
	t.Setenv("JWT_ACCESS_SECRET", "access-secret")
	t.Setenv("JWT_REFRESH_SECRET", "refresh-secret")

	access, _ := GenerateAccessToken("5", "a@test.com", "user")
	refresh, _ := GenerateRefreshToken("5")

	if _, err := ValidateRefreshToken(access); err == nil {
		t.Fatal("access token must not validate as refresh token")
	}
	if _, err := ValidateAccessToken(refresh); err == nil {
		t.Fatal("refresh token must not validate as access token")
	}
}

func TestSecretReadAtCallTime(t *testing.T) {
	// Regression: secrets used to be captured at package init,
	// before main() loaded .env.
	t.Setenv("JWT_ACCESS_SECRET", "set-after-init")
	tok, err := GenerateAccessToken("1", "x@test.com", "user")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ValidateAccessToken(tok); err != nil {
		t.Fatalf("round trip with late secret failed: %v", err)
	}
}

func TestWrongSecretRejected(t *testing.T) {
	t.Setenv("JWT_ACCESS_SECRET", "one")
	tok, _ := GenerateAccessToken("1", "x@test.com", "user")

	t.Setenv("JWT_ACCESS_SECRET", "two")
	if _, err := ValidateAccessToken(tok); err == nil {
		t.Fatal("token signed with a different secret must be rejected")
	}
}
