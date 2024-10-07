package config

import (
	"testing"
)

func TestVersion(t *testing.T) {
	originalVersion := Version

	newVersion := "v1.2.3"
	Version = newVersion
	if Version != newVersion {
		t.Errorf("Expected version %s, but got %s", newVersion, Version)
	}

	Version = originalVersion
	if Version != originalVersion {
		t.Errorf("Expected version %s, but got %s", originalVersion, Version)
	}
}

func TestGetUserAgent(t *testing.T) {
	originalVersion := Version

	expectedUserAgent := "fixctl-dev"
	if userAgent := GetUserAgent(); userAgent != expectedUserAgent {
		t.Errorf("Expected User-Agent %s, but got %s", expectedUserAgent, userAgent)
	}

	Version = "v1.2.3"
	expectedUserAgent = "fixctl-v1.2.3"
	if userAgent := GetUserAgent(); userAgent != expectedUserAgent {
		t.Errorf("Expected User-Agent %s, but got %s", expectedUserAgent, userAgent)
	}

	Version = originalVersion
}
