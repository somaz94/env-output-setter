package writer

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/somaz94/env-output-setter/internal/jsonutil"
	"github.com/somaz94/env-output-setter/internal/printer"
)

// JSONHandler handles JSON value processing and extraction.
type JSONHandler struct{}

// NewJSONHandler creates a new JSONHandler instance.
func NewJSONHandler() *JSONHandler {
	return &JSONHandler{}
}

// ProcessJSONValues extracts nested properties from JSON values.
// It processes JSON objects and arrays, creating flattened key-value pairs.
func (h *JSONHandler) ProcessJSONValues(keyList, valueList []string) ([]string, []string) {
	// Make a copy of the original lists
	originalKeyCount := len(keyList)
	resultKeys := make([]string, len(keyList))
	resultValues := make([]string, len(valueList))
	copy(resultKeys, keyList)
	copy(resultValues, valueList)

	// Process each JSON value in the original list
	for i := 0; i < originalKeyCount; i++ {
		value := valueList[i]
		key := keyList[i]

		// Check if value looks like JSON
		if !jsonutil.IsJSONLike(value) {
			continue
		}

		// Try to parse the JSON value
		var jsonData interface{}
		if err := json.Unmarshal([]byte(value), &jsonData); err != nil {
			printer.PrintWarning(fmt.Sprintf("Warning: Invalid JSON for key '%s': %v", key, err))
			continue
		}

		// Extract nested values based on the JSON type
		switch typedData := jsonData.(type) {
		case map[string]interface{}:
			// Handle JSON object
			nestedKeys, nestedValues := h.extractNestedJSON(key, typedData, "")
			resultKeys = append(resultKeys, nestedKeys...)
			resultValues = append(resultValues, nestedValues...)
		case []interface{}:
			// Handle JSON array
			for idx, item := range typedData {
				arrayKey := fmt.Sprintf("%s_%d", key, idx)
				resultKeys = append(resultKeys, arrayKey)
				resultValues = append(resultValues, fmt.Sprintf("%v", item))

				// Process nested objects in arrays
				if mapItem, ok := item.(map[string]interface{}); ok {
					objKeys, objValues := h.extractNestedJSON(arrayKey, mapItem, "")
					resultKeys = append(resultKeys, objKeys...)
					resultValues = append(resultValues, objValues...)
				}
			}
		}
	}

	return resultKeys, resultValues
}

// extractNestedJSON flattens a nested JSON object into key-value pairs.
// It recursively processes nested objects and arrays, creating concatenated keys.
func (h *JSONHandler) extractNestedJSON(prefix string, jsonObj map[string]interface{}, groupPrefix string) ([]string, []string) {
	var keys []string
	var values []string

	// Prepare the key prefix with group prefix if provided
	keyPrefix := prefix
	if groupPrefix != "" && !strings.HasPrefix(prefix, groupPrefix) {
		keyPrefix = fmt.Sprintf("%s_%s", groupPrefix, prefix)
	}

	// Process each property in the JSON object
	for propKey, propValue := range jsonObj {
		nestedKey := fmt.Sprintf("%s_%s", keyPrefix, propKey)

		// Handle different value types
		switch typedValue := propValue.(type) {
		case map[string]interface{}:
			// Recursively process nested objects
			nestedKeys, nestedValues := h.extractNestedJSON(nestedKey, typedValue, "")
			keys = append(keys, nestedKeys...)
			values = append(values, nestedValues...)
		case []interface{}:
			// Process arrays
			for i, item := range typedValue {
				arrayKey := fmt.Sprintf("%s_%d", nestedKey, i)
				keys = append(keys, arrayKey)
				values = append(values, fmt.Sprintf("%v", item))

				// Process objects within arrays
				if mapItem, ok := item.(map[string]interface{}); ok {
					subKeys, subValues := h.extractNestedJSON(arrayKey, mapItem, "")
					keys = append(keys, subKeys...)
					values = append(values, subValues...)
				}
			}
		default:
			// Handle primitive values
			keys = append(keys, nestedKey)
			values = append(values, fmt.Sprintf("%v", typedValue))
		}
	}

	return keys, values
}
