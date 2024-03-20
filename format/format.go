package format

import (
	"encoding/json"

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
