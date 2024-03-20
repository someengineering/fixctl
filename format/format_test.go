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
