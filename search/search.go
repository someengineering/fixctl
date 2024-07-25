// file: search/search.go
package search

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

type SearchRequest struct {
	Query     string `json:"query"`
	WithEdges bool   `json:"with_edges"`
}

func escapeSingleQuotes(s string) string {
	return strings.ReplaceAll(s, "'", "'\\''")
}

func SearchGraph(apiEndpoint, fixJWT, workspaceID, searchStr string, withEdges bool) (<-chan interface{}, <-chan error) {
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
		req.AddCookie(&http.Cookie{
			Name:     "session_token",
			Value:    fixJWT,
			HttpOnly: true,
			Secure:   true,
		})

		escapedRequestBody := escapeSingleQuotes(string(requestBody))
		curlCommand := fmt.Sprintf("curl -X POST -H 'Content-Type: application/json' -H 'Accept: application/ndjson' -H 'Cookie: session_token=%s' -d '%s' %s", fixJWT, escapedRequestBody, url)
		logrus.Debugln("Equivalent curl command:", curlCommand)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			errs <- fmt.Errorf("error making HTTP request: %w", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			bodyBytes, readErr := io.ReadAll(resp.Body)
			if readErr != nil {
				errs <- fmt.Errorf("search request failed with status code: %d, and error reading response body: %w", resp.StatusCode, readErr)
				return
			}
			errs <- fmt.Errorf("search request failed with status code: %d, error: %s", resp.StatusCode, string(bodyBytes))
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
				errs <- fmt.Errorf("error unmarshalling JSON: %w", err)
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
