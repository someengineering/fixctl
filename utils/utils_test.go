package utils

import (
	"os"
	"strings"
	"testing"
)

func TestGetEnvOrDefault(t *testing.T) {
	tests := []struct {
		envKey       string
		setValue     string
		defaultValue string
		want         string
	}{
		{"TEST_ENV_SET", "actualValue", "defaultValue", "actualValue"},
		{"TEST_ENV_NOT_SET", "", "defaultValue", "defaultValue"},
	}

	for _, tt := range tests {
		if tt.setValue != "" {
			os.Setenv(tt.envKey, tt.setValue)
		} else {
			os.Unsetenv(tt.envKey)
		}

		got := GetEnvOrDefault(tt.envKey, tt.defaultValue)

		if got != tt.want {
			t.Errorf("GetEnvOrDefault(%q, %q) = %q, want %q", tt.envKey, tt.defaultValue, got, tt.want)
		}

		os.Unsetenv(tt.envKey)
	}
}

func TestSanitizeAPIEndpoint(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"Empty Endpoint", "", "", true},
		{"Valid HTTPS fix.security", "https://api.fix.security", "https://api.fix.security", false},
		{"Valid HTTP Localhost", "http://localhost:8080", "http://localhost:8080", false},
		{"Invalid Scheme fix.security", "http://api.fix.security", "", true},
		{"Invalid Domain", "https://api.example.com", "", true},
		{"Trailing Slash HTTPS fix.security", "https://api.fix.security/", "https://api.fix.security", false},
		{"Trailing Slash HTTP Localhost", "http://localhost:8080/", "http://localhost:8080", false},
		{"Trailing Slash Invalid Scheme", "http://api.fix.security/", "", true},
	}

	for _, tt := range tests {
		got, err := SanitizeAPIEndpoint(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. SanitizeAPIEndpoint() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("%q. SanitizeAPIEndpoint() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestSanitizeCredentials(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
		wantErr  bool
	}{
		{"Valid Credentials", "user", "pass", false},
		{"Username With Space", "user name", "pass", true},
		{"Password With Space", "user", "pa ss", true},
		{"Long Username", strings.Repeat("a", 129), "pass", true},
		{"Long Password", "user", strings.Repeat("a", 129), true},
	}

	for _, tt := range tests {
		_, _, err := SanitizeCredentials(tt.username, tt.password)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. SanitizeCredentials() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestSanitizeSearchString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"Valid Search", "query", false},
		{"Empty Search", "", true},
		{"Long Search", strings.Repeat("a", 4097), true},
	}

	for _, tt := range tests {
		_, err := SanitizeSearchString(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. SanitizeSearchString() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestSanitizeToken(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{"Valid Token", "token", false},
		{"Long Token", strings.Repeat("a", 4097), true},
	}

	for _, tt := range tests {
		_, err := SanitizeToken(tt.token)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. SanitizeToken() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestSanitizeWorkspaceId(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{"Valid GUID", "123e4567-e89b-12d3-a456-426614174000", false},
		{"Invalid GUID", "123", true},
		{"Empty GUID", "", true},
	}

	for _, tt := range tests {
		_, err := SanitizeWorkspaceId(tt.id)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. SanitizeWorkspaceId() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}
