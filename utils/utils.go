package utils

import (
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
)

func GetEnvOrDefault(envKey, defaultValue string) string {
	if value, exists := os.LookupEnv(envKey); exists {
		return value
	}
	return defaultValue
}

func SanitizeAPIEndpoint(endpoint string) (string, error) {
	if endpoint == "" {
		return "", fmt.Errorf("API endpoint is empty")
	}

	u, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	}

	isLocal := strings.HasPrefix(u.Hostname(), "localhost") || strings.HasPrefix(u.Hostname(), "127.") || u.Hostname() == "::1"
	isSecureDomain := strings.HasSuffix(u.Hostname(), "fix.security") || strings.HasSuffix(u.Hostname(), "fixcloud.io")

	switch {
	case isSecureDomain && u.Scheme != "https":
		return "", fmt.Errorf("API endpoint must use https scheme")
	case !isLocal && !isSecureDomain:
		return "", fmt.Errorf("invalid API endpoint")
	}

	endpoint = strings.TrimSuffix(endpoint, "/")
	return endpoint, nil
}

func SanitizeCredentials(username, password string) (string, string, error) {
	if strings.Contains(username, " ") || strings.Contains(password, " ") || len(username) > 128 || len(password) > 128 {
		return "", "", fmt.Errorf("username or password contains spaces or is too long")
	}
	return username, password, nil
}

func SanitizeSearchString(search string) (string, error) {
	if search == "" {
		return "", fmt.Errorf("search string is empty")
	}
	if len(search) > 4096 {
		return "", fmt.Errorf("search string is too long")
	}
	return search, nil
}

func SanitizeToken(token string) (string, error) {
	if len(token) > 4096 {
		return "", fmt.Errorf("token is too long")
	}
	return token, nil
}

func SanitizeWorkspaceId(workspaceId string) (string, error) {
	guidRegex := regexp.MustCompile(`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[1-5][a-fA-F0-9]{3}-[89abAB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$`)

	if !guidRegex.MatchString(workspaceId) {
		return "", fmt.Errorf("workspace ID is not a valid GUID")
	}

	return workspaceId, nil
}

func SanitizeCSVHeaders(headers string) ([]string, error) {
	if headers == "" {
		return nil, fmt.Errorf("headers cannot be empty")
	}

	rawHeaders := strings.Split(headers, ",")
	if len(rawHeaders) == 0 {
		return nil, fmt.Errorf("at least one header must be specified")
	}

	csvHeaders := make([]string, len(rawHeaders))
	for i, header := range rawHeaders {
		trimmedHeader := strings.TrimSpace(header)
		if trimmedHeader == "" {
			return nil, fmt.Errorf("empty CSV header found")
		}

		if !strings.HasPrefix(trimmedHeader, "/") {
			trimmedHeader = "/reported." + trimmedHeader
		}
		csvHeaders[i] = trimmedHeader
	}
	return csvHeaders, nil
}

func SanitizeOutputFormat(format string) (string, error) {
	switch format {
	case "json", "yaml", "csv":
		return format, nil
	default:
		return "", fmt.Errorf("unsupported output format")
	}
}
