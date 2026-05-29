// Package payload provides helpers for loading and merging JSON request
// payloads supplied to commands via the --input/-i flag.
package payload

import (
	"encoding/json"
	"fmt"
	"os"
)

// LoadPayloadFile reads a JSON file from path and returns it as a generic
// map. An empty path yields an empty map (used to mean "no file supplied").
// Missing files produce a user-facing "input file not found" error.
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

// MergePayload combines a base map with overrides, returning a new map.
// Non-nil values in overrides take precedence over base, and any keys whose
// final value is nil are dropped from the result.
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
