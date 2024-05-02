// file: search/search.go
package search

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type SearchRequest struct {
	Query     string `json:"query"`
	WithEdges bool   `json:"with_edges"`
}

func SearchGraph(apiEndpoint, fixToken, workspaceID, searchStr string, withEdges bool) (<-chan interface{}, <-chan error) {
	results := make(chan interface{})
	errs := make(chan error, 1)

	go func() {
		defer close(results)
		defer close(errs)
		requestBody, err := json.Marshal(SearchRequest{
			Query:     searchStr,
			WithEdges: withEdges,
		})
		if err != nil {
			errs <- fmt.Errorf("error marshaling JSON: %w", err)
			return
		}

		url := fmt.Sprintf("%s/api/workspaces/%s/inventory/search", apiEndpoint, workspaceID)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
		if err != nil {
			errs <- fmt.Errorf("error creating request: %w", err)
			return
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/ndjson")
		req.AddCookie(&http.Cookie{Name: "session_token", Value: fixToken})

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			errs <- fmt.Errorf("error making HTTP request: %w", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			errs <- fmt.Errorf("search request failed with status code: %d", resp.StatusCode)
			return
		}

		scanner := bufio.NewScanner(resp.Body)
		const maxTokenSize = 1024 * 5120
		buf := make([]byte, maxTokenSize)
		scanner.Buffer(buf, maxTokenSize)

		for scanner.Scan() {
			decoder := json.NewDecoder(bytes.NewReader(scanner.Bytes()))
			decoder.UseNumber()

			var result interface{}
			if err := decoder.Decode(&result); err != nil {
				errs <- fmt.Errorf("error unmarshaling JSON: %w", err)
				return
			}
			results <- result
		}

		if err := scanner.Err(); err != nil {
			errs <- fmt.Errorf("error reading response body: %w", err)
			return
		}
	}()

	return results, errs
}
