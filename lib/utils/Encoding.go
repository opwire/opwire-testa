package utils

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
)

func Unmarshal(format string, source []byte, target interface{}) error {
	if format == BODY_FORMAT_JSON {
		return json.Unmarshal(source, target)
	}
	if format == BODY_FORMAT_YAML {
		return yaml.Unmarshal(source, target)
	}
	return fmt.Errorf("Invalid body format: %s", format)
}
