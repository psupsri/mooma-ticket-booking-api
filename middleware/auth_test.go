package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const testSecret = "test-secret"

func sign(t *testing.T, method jwt.SigningMethod, key any, claims jwt.MapClaims) string {
	t.Helper()
	token, err := jwt.NewWithClaims(method, claims).SignedString(key)
	if err != nil {
		t.Fatal(err)
	}
	return token
}

func validClaims() jwt.MapClaims {
	return jwt.MapClaims{"customer_id": 42, "exp": time.Now().Add(time.Hour).Unix()}
}

// serve runs AuthMiddleware in front of a handler that echoes the customer ID
// it was given, so a test sees both the status and what reached the handler.
func serve(t *testing.T, authorization string) (int, map[string]any) {
	t.Helper()
	t.Setenv("JWT_SECRET", testSecret)
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.GET("/protected", AuthMiddleware, func(ctx *gin.Context) {
		id, _ := ctx.Get("customer_id")
		ctx.JSON(http.StatusOK, gin.H{"customer_id": id})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	if authorization != "" {
		req.Header.Set("Authorization", authorization)
	}
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	var body map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &body)
	return rec.Code, body
}

func TestAuthMiddlewareAcceptsValidToken(t *testing.T) {
	token := sign(t, jwt.SigningMethodHS256, []byte(testSecret), validClaims())

	status, body := serve(t, "Bearer "+token)

	if status != http.StatusOK {
		t.Fatalf("status = %d, want 200 (body %v)", status, body)
	}
	if body["customer_id"] != float64(42) {
		t.Fatalf("handler got customer_id %v, want 42", body["customer_id"])
	}
}

func TestAuthMiddlewareRejects(t *testing.T) {
	expired := validClaims()
	expired["exp"] = time.Now().Add(-time.Minute).Unix()

	noCustomer := validClaims()
	delete(noCustomer, "customer_id")

	stringCustomer := validClaims()
	stringCustomer["customer_id"] = "42"

	tests := []struct {
		name          string
		authorization string
	}{
		{"no header", ""},
		{"not a JWT", "Bearer not-a-token"},
		{"missing Bearer token", "Bearer "},
		{"expired", "Bearer " + sign(t, jwt.SigningMethodHS256, []byte(testSecret), expired)},
		// The production bug this guards: tokens signed with an empty or
		// different secret must never be accepted.
		{"signed with an empty secret", "Bearer " + sign(t, jwt.SigningMethodHS256, []byte(""), validClaims())},
		{"signed with another secret", "Bearer " + sign(t, jwt.SigningMethodHS256, []byte("attacker"), validClaims())},
		{"other HMAC algorithm", "Bearer " + sign(t, jwt.SigningMethodHS512, []byte(testSecret), validClaims())},
		{"unsigned (alg none)", "Bearer " + sign(t, jwt.SigningMethodNone, jwt.UnsafeAllowNoneSignatureType, validClaims())},
		{"no customer_id claim", "Bearer " + sign(t, jwt.SigningMethodHS256, []byte(testSecret), noCustomer)},
		{"customer_id is not a number", "Bearer " + sign(t, jwt.SigningMethodHS256, []byte(testSecret), stringCustomer)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, body := serve(t, tt.authorization)
			if status != http.StatusUnauthorized {
				t.Fatalf("status = %d, want 401 (body %v)", status, body)
			}
			if _, reached := body["customer_id"]; reached {
				t.Fatalf("request reached the handler: %v", body)
			}
		})
	}
}
