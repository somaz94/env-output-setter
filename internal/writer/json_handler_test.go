package writer

import (
	"strings"
	"testing"

	"github.com/somaz94/env-output-setter/internal/jsonutil"
)

func TestNewJSONHandler(t *testing.T) {
	handler := NewJSONHandler()

	if handler == nil {
		t.Fatal("NewJSONHandler() returned nil")
	}
}

func TestIsJSONLike(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{
			name:     "Valid JSON object",
			value:    `{"key":"value"}`,
			expected: true,
		},
		{
			name:     "Valid JSON array",
			value:    `["item1","item2"]`,
			expected: true,
		},
		{
			name:     "JSON object with whitespace",
			value:    `  {"key":"value"}  `,
			expected: true,
		},
		{
			name:     "Plain text",
			value:    "not json",
			expected: false,
		},
		{
			name:     "Incomplete JSON",
			value:    `{"key":`,
			expected: false,
		},
		{
			name:     "Empty string",
			value:    "",
			expected: false,
		},
		{
			name:     "Just opening brace",
			value:    "{",
			expected: false,
		},
		{
			name:     "Just closing brace",
			value:    "}",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := jsonutil.IsJSONLike(tt.value)

			if result != tt.expected {
				t.Errorf("IsJSONLike() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestProcessJSONValues(t *testing.T) {
	tests := []struct {
		name              string
		keyList           []string
		valueList         []string
		minExpectedKeys   int
		minExpectedValues int
	}{
		{
			name:              "No JSON values",
			keyList:           []string{"KEY1", "KEY2"},
			valueList:         []string{"VALUE1", "VALUE2"},
			minExpectedKeys:   2,
			minExpectedValues: 2,
		},
		{
			name:              "Single JSON object",
			keyList:           []string{"CONFIG"},
			valueList:         []string{`{"server":"localhost","port":8080}`},
			minExpectedKeys:   1,
			minExpectedValues: 1,
		},
		{
			name:              "JSON array",
			keyList:           []string{"ITEMS"},
			valueList:         []string{`["item1","item2","item3"]`},
			minExpectedKeys:   1,
			minExpectedValues: 1,
		},
		{
			name:              "Mixed JSON and plain values",
			keyList:           []string{"CONFIG", "STATUS"},
			valueList:         []string{`{"host":"localhost"}`, "active"},
			minExpectedKeys:   2,
			minExpectedValues: 2,
		},
		{
			name:              "Nested JSON object",
			keyList:           []string{"APP"},
			valueList:         []string{`{"server":{"host":"localhost","port":8080}}`},
			minExpectedKeys:   1,
			minExpectedValues: 1,
		},
		{
			name:              "Invalid JSON",
			keyList:           []string{"BAD"},
			valueList:         []string{`{invalid json}`},
			minExpectedKeys:   1,
			minExpectedValues: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewJSONHandler()
			resultKeys, resultValues := handler.ProcessJSONValues(tt.keyList, tt.valueList)

			if len(resultKeys) < tt.minExpectedKeys {
				t.Errorf("ProcessJSONValues() resultKeys length = %d, want at least %d",
					len(resultKeys), tt.minExpectedKeys)
			}

			if len(resultValues) < tt.minExpectedValues {
				t.Errorf("ProcessJSONValues() resultValues length = %d, want at least %d",
					len(resultValues), tt.minExpectedValues)
			}

			// Verify keys and values have same length
			if len(resultKeys) != len(resultValues) {
				t.Errorf("ProcessJSONValues() keys length %d != values length %d",
					len(resultKeys), len(resultValues))
			}
		})
	}
}

func TestExtractNestedJSON(t *testing.T) {
	tests := []struct {
		name        string
		prefix      string
		jsonObj     map[string]interface{}
		groupPrefix string
		minKeys     int
		minValues   int
	}{
		{
			name:   "Simple flat object",
			prefix: "CONFIG",
			jsonObj: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
			groupPrefix: "",
			minKeys:     2,
			minValues:   2,
		},
		{
			name:   "Nested object",
			prefix: "APP",
			jsonObj: map[string]interface{}{
				"server": map[string]interface{}{
					"host": "localhost",
					"port": 8080,
				},
			},
			groupPrefix: "",
			minKeys:     2,
			minValues:   2,
		},
		{
			name:   "Object with array",
			prefix: "DATA",
			jsonObj: map[string]interface{}{
				"items": []interface{}{"item1", "item2"},
			},
			groupPrefix: "",
			minKeys:     2,
			minValues:   2,
		},
		{
			name:   "Mixed types",
			prefix: "CONFIG",
			jsonObj: map[string]interface{}{
				"name":    "test",
				"enabled": true,
				"count":   42,
			},
			groupPrefix: "",
			minKeys:     3,
			minValues:   3,
		},
		{
			name:        "Empty object",
			prefix:      "EMPTY",
			jsonObj:     map[string]interface{}{},
			groupPrefix: "",
			minKeys:     0,
			minValues:   0,
		},
		{
			name:   "With group prefix",
			prefix: "CONFIG",
			jsonObj: map[string]interface{}{
				"setting": "value",
			},
			groupPrefix: "APP",
			minKeys:     1,
			minValues:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewJSONHandler()
			keys, values := handler.extractNestedJSON(tt.prefix, tt.jsonObj, tt.groupPrefix)

			if len(keys) < tt.minKeys {
				t.Errorf("extractNestedJSON() keys length = %d, want at least %d",
					len(keys), tt.minKeys)
			}

			if len(values) < tt.minValues {
				t.Errorf("extractNestedJSON() values length = %d, want at least %d",
					len(values), tt.minValues)
			}

			// Verify keys and values have same length
			if len(keys) != len(values) {
				t.Errorf("extractNestedJSON() keys length %d != values length %d",
					len(keys), len(values))
			}
		})
	}
}

func TestProcessJSONValuesWithComplexStructures(t *testing.T) {
	tests := []struct {
		name      string
		keyList   []string
		valueList []string
		checkKey  string
		wantKey   bool
	}{
		{
			name:      "Extract nested properties",
			keyList:   []string{"CONFIG"},
			valueList: []string{`{"database":{"host":"localhost","port":5432}}`},
			checkKey:  "CONFIG_database_host",
			wantKey:   true,
		},
		{
			name:      "Array with objects",
			keyList:   []string{"SERVERS"},
			valueList: []string{`[{"name":"server1"},{"name":"server2"}]`},
			checkKey:  "SERVERS_0",
			wantKey:   true,
		},
		{
			name:      "Deep nesting",
			keyList:   []string{"APP"},
			valueList: []string{`{"level1":{"level2":{"level3":"value"}}}`},
			checkKey:  "APP_level1_level2_level3",
			wantKey:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewJSONHandler()
			resultKeys, _ := handler.ProcessJSONValues(tt.keyList, tt.valueList)

			found := false
			for _, key := range resultKeys {
				if key == tt.checkKey {
					found = true
					break
				}
			}

			if tt.wantKey && !found {
				t.Errorf("ProcessJSONValues() expected key %q not found in %v",
					tt.checkKey, resultKeys)
			}
		})
	}
}

func TestExtractNestedJSONWithArrayOfObjects(t *testing.T) {
	handler := NewJSONHandler()
	jsonObj := map[string]interface{}{
		"servers": []interface{}{
			map[string]interface{}{"name": "s1", "port": 80},
			map[string]interface{}{"name": "s2", "port": 443},
		},
	}

	keys, values := handler.extractNestedJSON("APP", jsonObj, "")

	if len(keys) < 4 {
		t.Errorf("extractNestedJSON() expected at least 4 keys for array of objects, got %d: %v", len(keys), keys)
	}

	if len(keys) != len(values) {
		t.Errorf("extractNestedJSON() keys/values length mismatch: %d vs %d", len(keys), len(values))
	}
}

func TestExtractNestedJSONWithGroupPrefixAlreadyPresent(t *testing.T) {
	handler := NewJSONHandler()
	jsonObj := map[string]interface{}{
		"key": "value",
	}

	// prefix already starts with groupPrefix
	keys, values := handler.extractNestedJSON("GRP_CONFIG", jsonObj, "GRP")

	if len(keys) != 1 {
		t.Errorf("extractNestedJSON() expected 1 key, got %d: %v", len(keys), keys)
	}

	// Should not double-prefix
	if len(keys) > 0 && strings.HasPrefix(keys[0], "GRP_GRP_") {
		t.Errorf("extractNestedJSON() double-prefixed key: %s", keys[0])
	}

	_ = values
}

func BenchmarkProcessJSONValues(b *testing.B) {
	handler := NewJSONHandler()
	keyList := []string{"CONFIG", "STATUS"}
	valueList := []string{
		`{"server":"localhost","port":8080,"database":{"host":"db.local","port":5432}}`,
		"active",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler.ProcessJSONValues(keyList, valueList)
	}
}

func BenchmarkExtractNestedJSON(b *testing.B) {
	handler := NewJSONHandler()
	jsonObj := map[string]interface{}{
		"server": map[string]interface{}{
			"host": "localhost",
			"port": 8080,
		},
		"database": map[string]interface{}{
			"host": "db.local",
			"port": 5432,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler.extractNestedJSON("CONFIG", jsonObj, "")
	}
}

func TestProcessJSONValuesPreservesOriginal(t *testing.T) {
	handler := NewJSONHandler()
	originalKeys := []string{"KEY1", "KEY2"}
	originalValues := []string{"VALUE1", "VALUE2"}

	resultKeys, resultValues := handler.ProcessJSONValues(originalKeys, originalValues)

	// Original keys should still be present
	if len(resultKeys) < len(originalKeys) {
		t.Error("ProcessJSONValues() should preserve original keys")
	}

	if len(resultValues) < len(originalValues) {
		t.Error("ProcessJSONValues() should preserve original values")
	}

	// Check first elements match
	if resultKeys[0] != originalKeys[0] {
		t.Errorf("ProcessJSONValues() first key = %v, want %v", resultKeys[0], originalKeys[0])
	}

	if resultValues[0] != originalValues[0] {
		t.Errorf("ProcessJSONValues() first value = %v, want %v", resultValues[0], originalValues[0])
	}
}
