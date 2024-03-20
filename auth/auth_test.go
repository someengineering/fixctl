package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoginAndGetJWT(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/api/auth/jwt/login" {
			t.Errorf("Expected POST /api/auth/jwt/login, got %s %s", r.Method, r.URL.Path)
		}

		w.Header().Set("Set-Cookie", "session_token=mocked_jwt_token; Path=/; HttpOnly")
		w.WriteHeader(204)
	}))
	defer mockServer.Close()

	jwt, err := LoginAndGetJWT(mockServer.URL, "user", "pass")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedJWT := "mocked_jwt_token"
	if jwt != expectedJWT {
		t.Errorf("Expected JWT %s, got %s", expectedJWT, jwt)
	}
}
