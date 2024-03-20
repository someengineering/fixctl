// file: search/search.go
package search

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// SearchRequest defines the structure of the search request body
type SearchRequest struct {
	Query string      `json:"query"`
	Skip  int         `json:"skip"`
	Limit int         `json:"limit"`
	Count bool        `json:"count"`
	Sort  []SortField `json:"sort"`
}

// SortField defines the structure of the sort field in the search request
type SortField struct {
	Path      string `json:"path"`
	Direction string `json:"direction"`
}

// SearchTable performs the HTTP POST request to search in the inventory
func SearchTable(apiEndpoint, fixToken, workspaceID, searchStr string) (string, error) {
	requestBody, err := json.Marshal(SearchRequest{
		Query: searchStr,
		Skip:  0,
		Limit: 50,
		Count: false,
		Sort: []SortField{
			{Path: "string", Direction: "asc"},
		},
	})
	if err != nil {
		return "", fmt.Errorf("error marshaling JSON: %w", err)
	}

	url := fmt.Sprintf("%s/api/workspaces/%s/inventory/search/table", apiEndpoint, workspaceID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.AddCookie(&http.Cookie{Name: "session_token", Value: fixToken})

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making HTTP request: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("search request failed with status code: %d, response: %s", resp.StatusCode, responseBody)
	}
	return string(responseBody), nil
}
