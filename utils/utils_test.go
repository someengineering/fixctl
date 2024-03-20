package utils

import (
	"os"
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
