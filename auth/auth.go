package auth

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func LoginAndGetJWT(apiEndpoint, username, password string) (string, error) {
	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)

	loginURL := fmt.Sprintf("%s/api/auth/jwt/login", apiEndpoint)
	req, err := http.NewRequest("POST", loginURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("creating login request failed: %w", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("login request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		return "", fmt.Errorf("login failed with status code: %d", resp.StatusCode)
	}

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "session_token" {
			return cookie.Value, nil
		}
	}

	return "", fmt.Errorf("JWT not found in response cookies")
}
