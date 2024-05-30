package auth

import (
	"encoding/json"
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

func TestGetJWTFromToken(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected 'POST' request, got '%s'", r.Method)
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected 'Content-Type: application/json', got '%s'", r.Header.Get("Content-Type"))
		}

		var body map[string]string
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("Error decoding request body: %v", err)
		}
		expectedToken := "test_token"
		if body["token"] != expectedToken {
			t.Errorf("Expected token '%s', got '%s'", expectedToken, body["token"])
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("mock_jwt_token"))
	}))
	defer mockServer.Close()

	apiEndpoint := mockServer.URL
	fixToken := "test_token"
	expectedJWT := "mock_jwt_token"

	jwt, err := GetJWTFromToken(apiEndpoint, fixToken)
	if err != nil {
		t.Fatalf("Expected no error, got '%v'", err)
	}

	if jwt != expectedJWT {
		t.Errorf("Expected JWT '%s', got '%s'", expectedJWT, jwt)
	}
}
