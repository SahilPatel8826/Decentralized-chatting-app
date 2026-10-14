package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"chat_service/pkg/jwtutil"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const secret = "mw-test-secret"

func token(t *testing.T) string {
	t.Helper()
	claims := &jwtutil.Claims{
		UserID: "7",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	s, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func setup() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := func(c *gin.Context) { c.String(200, c.GetString("userID")) }
	r.GET("/ws", AuthMiddleware(), h)
	r.GET("/other", AuthMiddleware(), h)
	return r
}

func do(r *gin.Engine, path, auth string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	if auth != "" {
		req.Header.Set("Authorization", auth)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestMissingTokenRejected(t *testing.T) {
	t.Setenv("JWT_ACCESS_SECRET", secret)
	if w := do(setup(), "/other", ""); w.Code != 401 {
		t.Fatalf("want 401, got %d", w.Code)
	}
}

func TestHeaderTokenAccepted(t *testing.T) {
	t.Setenv("JWT_ACCESS_SECRET", secret)
	w := do(setup(), "/other", "Bearer "+token(t))
	if w.Code != 200 || w.Body.String() != "7" {
		t.Fatalf("want 200/7, got %d/%s", w.Code, w.Body.String())
	}
}

func TestWSQueryTokenAccepted(t *testing.T) {
	t.Setenv("JWT_ACCESS_SECRET", secret)
	w := do(setup(), "/ws?token="+token(t), "")
	if w.Code != 200 || w.Body.String() != "7" {
		t.Fatalf("want 200/7, got %d/%s", w.Code, w.Body.String())
	}
}

func TestQueryTokenRejectedOnOtherRoutes(t *testing.T) {
	t.Setenv("JWT_ACCESS_SECRET", secret)
	if w := do(setup(), "/other?token="+token(t), ""); w.Code != 401 {
		t.Fatalf("query token must only work on /ws, got %d", w.Code)
	}
}

func TestBadTokenRejected(t *testing.T) {
	t.Setenv("JWT_ACCESS_SECRET", secret)
	if w := do(setup(), "/ws?token=garbage", ""); w.Code != 401 {
		t.Fatalf("want 401, got %d", w.Code)
	}
}
