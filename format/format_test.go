package format

import (
	"strings"
	"testing"
)

func TestToJSON(t *testing.T) {
	testData := map[string]interface{}{
		"name":  "Test Object",
		"value": 123,
	}
	expectedSubstring := `"name":"Test Object"`

	jsonStr, err := ToJSON(testData)
	if err != nil {
		t.Fatalf("ToJSON returned an error: %v", err)
	}
	if !strings.Contains(jsonStr, expectedSubstring) {
		t.Errorf("Expected JSON string to contain %s, got %s", expectedSubstring, jsonStr)
	}
	if jsonStr[len(jsonStr)-1] != '\n' {
		t.Errorf("Expected JSON string to end with a newline")
	}
}

func TestToYAML(t *testing.T) {
	testData := map[string]interface{}{
		"name":  "Test Object",
		"value": 123,
	}
	expectedSubstring := "name: Test Object"

	yamlStr, err := ToYAML(testData)
	if err != nil {
		t.Fatalf("ToYAML returned an error: %v", err)
	}
	if !strings.Contains(yamlStr, expectedSubstring) {
		t.Errorf("Expected YAML string to contain %s, got %s", expectedSubstring, yamlStr)
	}
}

func TestToCSV(t *testing.T) {
	tests := []struct {
		name    string
		data    map[string]interface{}
		headers []string
		want    string
	}{
		{
			name: "Simple fields",
			data: map[string]interface{}{
				"name": "Example",
				"id":   "123",
			},
			headers: []string{"/name", "/id"},
			want:    "Example,123\n",
		},
		{
			name: "Fields with commas and quotes",
			data: map[string]interface{}{
				"description": `Product "A", the best one`,
				"notes":       "It's, literally, \"awesome\".",
			},
			headers: []string{"/description", "/notes"},
			want:    "\"Product \"\"A\"\", the best one\",\"It's, literally, \"\"awesome\"\".\"\n",
		},
		{
			name: "Nested fields",
			data: map[string]interface{}{
				"reported": map[string]interface{}{
					"location": "Warehouse, 42",
					"status":   "In-stock",
				},
			},
			headers: []string{"/reported.location", "/reported.status"},
			want:    "\"Warehouse, 42\",In-stock\n",
		},
		{
			name: "Missing and empty fields",
			data: map[string]interface{}{
				"reported": map[string]interface{}{
					"location": "Remote",
				},
			},
			headers: []string{"/reported.location", "/reported.quantity"},
			want:    "Remote,\n",
		},
	}

	for _, tt := range tests {
		got, err := ToCSV(tt.data, tt.headers)
		if err != nil {
			t.Errorf("TestToCSV %s failed with error: %v", tt.name, err)
		}
		if got != tt.want {
			t.Errorf("TestToCSV %s expected %q, got %q", tt.name, tt.want, got)
		}
	}
}
