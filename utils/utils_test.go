package utils

import (
	"os"
	"reflect"
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

func TestSanitizeOutputFormat(t *testing.T) {
	tests := []struct {
		name      string
		format    string
		want      string
		wantError bool
	}{
		{
			name:      "Valid format json",
			format:    "json",
			want:      "json",
			wantError: false,
		},
		{
			name:      "Valid format yaml",
			format:    "yaml",
			want:      "yaml",
			wantError: false,
		},
		{
			name:      "Valid format csv",
			format:    "csv",
			want:      "csv",
			wantError: false,
		},
		{
			name:      "Unsupported format",
			format:    "xml",
			want:      "",
			wantError: true,
		},
		{
			name:      "Empty format",
			format:    "",
			want:      "",
			wantError: true,
		},
		{
			name:      "Whitespace format",
			format:    " ",
			want:      "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		got, err := SanitizeOutputFormat(tt.format)
		if (err != nil) != tt.wantError {
			t.Errorf("%s: SanitizeOutputFormat(%s) expected error: %v, got: %v", tt.name, tt.format, tt.wantError, err)
		}
		if got != tt.want {
			t.Errorf("%s: SanitizeOutputFormat(%s) = %v, want %v", tt.name, tt.format, got, tt.want)
		}
	}
}

func TestSanitizeCSVHeaders(t *testing.T) {
	tests := []struct {
		name      string
		headers   string
		want      []string
		wantError bool
	}{
		{
			name:      "Non-empty headers without leading slash",
			headers:   "id,name,kind",
			want:      []string{"/reported.id", "/reported.name", "/reported.kind"},
			wantError: false,
		},
		{
			name:      "Headers with leading slash",
			headers:   "/metadata.expires,/metadata.cleaned",
			want:      []string{"/metadata.expires", "/metadata.cleaned"},
			wantError: false,
		},
		{
			name:      "Mixed headers",
			headers:   "name,/metadata.expires,kind",
			want:      []string{"/reported.name", "/metadata.expires", "/reported.kind"},
			wantError: false,
		},
		{
			name:      "Empty headers string",
			headers:   "",
			want:      nil,
			wantError: true,
		},
		{
			name:      "Only whitespace",
			headers:   "  ",
			want:      nil,
			wantError: true,
		},
		{
			name:      "Header with only commas",
			headers:   ",,,",
			want:      nil,
			wantError: true,
		},
		{
			name:      "Empty header among valid headers",
			headers:   "id,,name",
			want:      nil,
			wantError: true,
		},
	}

	for _, tt := range tests {
		got, err := SanitizeCSVHeaders(tt.headers)
		if (err != nil) != tt.wantError {
			t.Errorf("%s: SanitizeCSVHeaders() error = %v, wantError %v", tt.name, err, tt.wantError)
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%s: SanitizeCSVHeaders() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
