package payload

import (
	"encoding/json"
	"fmt"
	"os"
)

func LoadPayloadFile(path string) (map[string]any, error) {
	if path == "" {
		return map[string]any{}, nil
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("input file not found: %s", path)
		}
		return nil, fmt.Errorf("cannot read input file %s: %w", path, err)
	}

	var data map[string]any
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, fmt.Errorf("invalid JSON in %s", path)
	}

	return data, nil
}

func MergePayload(base map[string]any, overrides map[string]any) map[string]any {
	merged := make(map[string]any)
	for k, v := range base {
		merged[k] = v
	}
	for k, v := range overrides {
		if v != nil {
			merged[k] = v
		}
	}
	result := make(map[string]any)
	for k, v := range merged {
		if v != nil {
			result[k] = v
		}
	}
	return result
}
