package format

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"
)

func ToJSON(data interface{}) (string, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	bytes = append(bytes, '\n')
	return string(bytes), nil
}

func ToYAML(data interface{}) (string, error) {
	bytes, err := yaml.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func ToCSV(data interface{}, headers []string) (string, error) {
	jsonObj, ok := data.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("data is not a JSON object")
	}

	var csvBuffer bytes.Buffer
	writer := csv.NewWriter(&csvBuffer)

	record := make([]string, len(headers))
	for i, header := range headers {
		header = strings.TrimPrefix(header, "/")
		var value interface{} = jsonObj

		for _, key := range strings.Split(header, ".") {
			if tempMap, ok := value.(map[string]interface{}); ok {
				value, ok = tempMap[key]
				if !ok {
					value = ""
					break
				}
			} else {
				break
			}
		}

		record[i] = fmt.Sprintf("%v", value)
	}

	if err := writer.Write(record); err != nil {
		return "", fmt.Errorf("writing record to CSV failed: %w", err)
	}

	writer.Flush()

	if err := writer.Error(); err != nil {
		return "", fmt.Errorf("CSV writing failed: %w", err)
	}

	return csvBuffer.String(), nil
}
