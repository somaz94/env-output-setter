package writer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/somaz94/env-output-setter/internal/config"
)

func TestNewProcessor(t *testing.T) {
	cfg := &config.Config{
		Delimiter:      ",",
		AllowEmpty:     false,
		TrimWhitespace: true,
	}

	processor := NewProcessor(cfg)

	if processor == nil {
		t.Fatal("NewProcessor() returned nil")
	}

	if processor.cfg != cfg {
		t.Error("NewProcessor() config not set correctly")
	}
}

func TestProcessInputValues(t *testing.T) {
	tests := []struct {
		name           string
		keys           string
		values         string
		delimiter      string
		allowEmpty     bool
		trimWhitespace bool
		jsonSupport    bool
		expectedKeys   []string
		expectedValues []string
	}{
		{
			name:           "Simple comma-separated values",
			keys:           "KEY1,KEY2,KEY3",
			values:         "VALUE1,VALUE2,VALUE3",
			delimiter:      ",",
			allowEmpty:     false,
			trimWhitespace: true,
			jsonSupport:    false,
			expectedKeys:   []string{"KEY1", "KEY2", "KEY3"},
			expectedValues: []string{"VALUE1", "VALUE2", "VALUE3"},
		},
		{
			name:           "Values with whitespace trimming",
			keys:           "  KEY1  ,  KEY2  ",
			values:         "  VALUE1  ,  VALUE2  ",
			delimiter:      ",",
			allowEmpty:     false,
			trimWhitespace: true,
			jsonSupport:    false,
			expectedKeys:   []string{"KEY1", "KEY2"},
			expectedValues: []string{"VALUE1", "VALUE2"},
		},
		{
			name:           "Empty values filtered out",
			keys:           "KEY1,,KEY3",
			values:         "VALUE1,,VALUE3",
			delimiter:      ",",
			allowEmpty:     false,
			trimWhitespace: true,
			jsonSupport:    false,
			expectedKeys:   []string{"KEY1", "KEY3"},
			expectedValues: []string{"VALUE1", "VALUE3"},
		},
		{
			name:           "Empty values allowed",
			keys:           "KEY1,,KEY3",
			values:         "VALUE1,,VALUE3",
			delimiter:      ",",
			allowEmpty:     true,
			trimWhitespace: true,
			jsonSupport:    false,
			expectedKeys:   []string{"KEY1", "", "KEY3"},
			expectedValues: []string{"VALUE1", "", "VALUE3"},
		},
		{
			name:           "Custom delimiter",
			keys:           "KEY1|KEY2|KEY3",
			values:         "VALUE1|VALUE2|VALUE3",
			delimiter:      "|",
			allowEmpty:     false,
			trimWhitespace: true,
			jsonSupport:    false,
			expectedKeys:   []string{"KEY1", "KEY2", "KEY3"},
			expectedValues: []string{"VALUE1", "VALUE2", "VALUE3"},
		},
		{
			name:           "Single key-value pair",
			keys:           "SINGLE_KEY",
			values:         "SINGLE_VALUE",
			delimiter:      ",",
			allowEmpty:     false,
			trimWhitespace: true,
			jsonSupport:    false,
			expectedKeys:   []string{"SINGLE_KEY"},
			expectedValues: []string{"SINGLE_VALUE"},
		},
		{
			name:           "Values with newlines",
			keys:           "KEY1,KEY2",
			values:         "VALUE1\nline2,VALUE2",
			delimiter:      ",",
			allowEmpty:     false,
			trimWhitespace: true,
			jsonSupport:    false,
			expectedKeys:   []string{"KEY1", "KEY2"},
			expectedValues: []string{"VALUE1 line2", "VALUE2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				Delimiter:      tt.delimiter,
				AllowEmpty:     tt.allowEmpty,
				TrimWhitespace: tt.trimWhitespace,
				JsonSupport:    tt.jsonSupport,
			}

			processor := NewProcessor(cfg)
			keyList, valueList, err := processor.ProcessInputValues(tt.keys, tt.values)

			if err != nil {
				t.Errorf("ProcessInputValues() unexpected error: %v", err)
			}

			if len(keyList) != len(tt.expectedKeys) {
				t.Errorf("ProcessInputValues() keyList length = %d, want %d", len(keyList), len(tt.expectedKeys))
			}

			if len(valueList) != len(tt.expectedValues) {
				t.Errorf("ProcessInputValues() valueList length = %d, want %d", len(valueList), len(tt.expectedValues))
			}

			for i, key := range keyList {
				if i < len(tt.expectedKeys) && key != tt.expectedKeys[i] {
					t.Errorf("ProcessInputValues() keyList[%d] = %v, want %v", i, key, tt.expectedKeys[i])
				}
			}

			for i, value := range valueList {
				if i < len(tt.expectedValues) && value != tt.expectedValues[i] {
					t.Errorf("ProcessInputValues() valueList[%d] = %v, want %v", i, value, tt.expectedValues[i])
				}
			}
		})
	}
}

func TestProcessWhitespace(t *testing.T) {
	tests := []struct {
		name     string
		entries  []string
		expected []string
	}{
		{
			name:     "Trim leading and trailing whitespace",
			entries:  []string{"  value1  ", "  value2  "},
			expected: []string{"value1", "value2"},
		},
		{
			name:     "Normalize newlines to spaces",
			entries:  []string{"value1\nvalue2", "value3\rvalue4"},
			expected: []string{"value1 value2", "value3 value4"},
		},
		{
			name:     "Condense multiple spaces",
			entries:  []string{"value1    value2", "value3  value4"},
			expected: []string{"value1 value2", "value3 value4"},
		},
		{
			name:     "No whitespace",
			entries:  []string{"value1", "value2"},
			expected: []string{"value1", "value2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{}
			processor := NewProcessor(cfg)
			result := processor.processWhitespace(tt.entries)

			if len(result) != len(tt.expected) {
				t.Errorf("processWhitespace() length = %d, want %d", len(result), len(tt.expected))
			}

			for i, value := range result {
				if i < len(tt.expected) && value != tt.expected[i] {
					t.Errorf("processWhitespace()[%d] = %v, want %v", i, value, tt.expected[i])
				}
			}
		})
	}
}

func TestNormalizeWhitespace(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Newlines to spaces",
			input:    "line1\nline2",
			expected: "line1 line2",
		},
		{
			name:     "Carriage returns to spaces",
			input:    "line1\rline2",
			expected: "line1 line2",
		},
		{
			name:     "Multiple spaces condensed",
			input:    "word1    word2",
			expected: "word1 word2",
		},
		{
			name:     "Mixed whitespace",
			input:    "word1\n\n  word2\r\n  word3",
			expected: "word1 word2 word3",
		},
		{
			name:     "No whitespace",
			input:    "simpletext",
			expected: "simpletext",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{}
			processor := NewProcessor(cfg)
			result := processor.normalizeWhitespace(tt.input)

			if result != tt.expected {
				t.Errorf("normalizeWhitespace() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestRemoveEmptyEntries(t *testing.T) {
	tests := []struct {
		name       string
		entries    []string
		allowEmpty bool
		expected   []string
	}{
		{
			name:       "Filter out empty strings",
			entries:    []string{"value1", "", "value2", "  ", "value3"},
			allowEmpty: false,
			expected:   []string{"value1", "value2", "value3"},
		},
		{
			name:       "Keep all entries when allowEmpty is true",
			entries:    []string{"value1", "", "value2"},
			allowEmpty: true,
			expected:   []string{"value1", "", "value2"},
		},
		{
			name:       "No empty entries",
			entries:    []string{"value1", "value2", "value3"},
			allowEmpty: false,
			expected:   []string{"value1", "value2", "value3"},
		},
		{
			name:       "All empty entries",
			entries:    []string{"", "  ", "   "},
			allowEmpty: false,
			expected:   []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				AllowEmpty: tt.allowEmpty,
			}
			processor := NewProcessor(cfg)
			result := processor.removeEmptyEntries(tt.entries)

			if len(result) != len(tt.expected) {
				t.Errorf("removeEmptyEntries() length = %d, want %d", len(result), len(tt.expected))
			}

			for i, value := range result {
				if i < len(tt.expected) && value != tt.expected[i] {
					t.Errorf("removeEmptyEntries()[%d] = %v, want %v", i, value, tt.expected[i])
				}
			}
		})
	}
}

func TestLogInputValues(t *testing.T) {
	tests := []struct {
		name      string
		debugMode bool
		varType   string
		keys      string
		values    string
	}{
		{
			name:      "Debug mode enabled",
			debugMode: true,
			varType:   "env",
			keys:      "KEY1,KEY2",
			values:    "VALUE1,VALUE2",
		},
		{
			name:      "Debug mode disabled",
			debugMode: false,
			varType:   "output",
			keys:      "STATUS",
			values:    "SUCCESS",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				DebugMode: tt.debugMode,
				Delimiter: ",",
			}
			processor := NewProcessor(cfg)

			// This should not panic
			processor.LogInputValues(tt.varType, tt.keys, tt.values)
		})
	}
}

func TestLogProcessedValues(t *testing.T) {
	tests := []struct {
		name      string
		debugMode bool
		keyList   []string
		valueList []string
	}{
		{
			name:      "Debug mode enabled with values",
			debugMode: true,
			keyList:   []string{"KEY1", "KEY2"},
			valueList: []string{"VALUE1", "VALUE2"},
		},
		{
			name:      "Debug mode disabled",
			debugMode: false,
			keyList:   []string{"STATUS"},
			valueList: []string{"SUCCESS"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				DebugMode: tt.debugMode,
			}
			processor := NewProcessor(cfg)

			// This should not panic
			processor.LogProcessedValues(tt.keyList, tt.valueList)
		})
	}
}

func TestProcessInputValuesWithJSON(t *testing.T) {
	tests := []struct {
		name          string
		keys          string
		values        string
		jsonSupport   bool
		minKeyCount   int
		minValueCount int
	}{
		{
			name:          "JSON support disabled",
			keys:          "KEY1,KEY2",
			values:        `{"nested":"value"},VALUE2`,
			jsonSupport:   false,
			minKeyCount:   2,
			minValueCount: 2,
		},
		{
			name:          "JSON support enabled",
			keys:          "CONFIG,STATUS",
			values:        `{"server":"localhost","port":8080},active`,
			jsonSupport:   true,
			minKeyCount:   2,
			minValueCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				Delimiter:   ",",
				AllowEmpty:  false,
				JsonSupport: tt.jsonSupport,
			}

			processor := NewProcessor(cfg)
			keyList, valueList, err := processor.ProcessInputValues(tt.keys, tt.values)

			if err != nil {
				t.Errorf("ProcessInputValues() unexpected error: %v", err)
			}

			if len(keyList) < tt.minKeyCount {
				t.Errorf("ProcessInputValues() keyList length = %d, want at least %d", len(keyList), tt.minKeyCount)
			}

			if len(valueList) < tt.minValueCount {
				t.Errorf("ProcessInputValues() valueList length = %d, want at least %d", len(valueList), tt.minValueCount)
			}
		})
	}
}

func TestProcessInputValuesWithInterpolation(t *testing.T) {
	t.Run("Interpolation enabled", func(t *testing.T) {
		t.Setenv("MY_HOST", "localhost")

		cfg := &config.Config{
			Delimiter:           ",",
			AllowEmpty:          false,
			EnableInterpolation: true,
			FileEncoding:        "raw",
		}
		processor := NewProcessor(cfg)
		keyList, valueList, err := processor.ProcessInputValues("SERVER,PORT", "${MY_HOST},${UNSET_PORT:-8080}")
		if err != nil {
			t.Fatalf("ProcessInputValues() error = %v", err)
		}
		if len(valueList) != 2 {
			t.Fatalf("expected 2 values, got %d", len(valueList))
		}
		if valueList[0] != "localhost" {
			t.Errorf("valueList[0] = %q, want %q", valueList[0], "localhost")
		}
		if valueList[1] != "8080" {
			t.Errorf("valueList[1] = %q, want %q", valueList[1], "8080")
		}
		_ = keyList
	})

	t.Run("Interpolation disabled", func(t *testing.T) {
		cfg := &config.Config{
			Delimiter:           ",",
			AllowEmpty:          false,
			EnableInterpolation: false,
			FileEncoding:        "raw",
		}
		processor := NewProcessor(cfg)
		_, valueList, err := processor.ProcessInputValues("KEY", "${SOME_VAR:-default}")
		if err != nil {
			t.Fatalf("ProcessInputValues() error = %v", err)
		}
		if valueList[0] != "${SOME_VAR:-default}" {
			t.Errorf("expected raw value, got %q", valueList[0])
		}
	})

	t.Run("Interpolation error propagated", func(t *testing.T) {
		cfg := &config.Config{
			Delimiter:           ",",
			AllowEmpty:          false,
			EnableInterpolation: true,
			FileEncoding:        "raw",
		}
		processor := NewProcessor(cfg)
		_, _, err := processor.ProcessInputValues("KEY", "${MISSING:?required var}")
		if err == nil {
			t.Error("expected error for missing required variable")
		}
	})
}

func TestProcessInputValuesWithFileReading(t *testing.T) {
	t.Run("Read value from file", func(t *testing.T) {
		tmpFile := filepath.Join(t.TempDir(), "val.txt")
		os.WriteFile(tmpFile, []byte("from_file"), 0644)

		cfg := &config.Config{
			Delimiter:    ",",
			AllowEmpty:   false,
			FileEncoding: "raw",
		}
		processor := NewProcessor(cfg)
		_, valueList, err := processor.ProcessInputValues("KEY", "file://"+tmpFile)
		if err != nil {
			t.Fatalf("ProcessInputValues() error = %v", err)
		}
		if valueList[0] != "from_file" {
			t.Errorf("valueList[0] = %q, want %q", valueList[0], "from_file")
		}
	})

	t.Run("File not found error", func(t *testing.T) {
		cfg := &config.Config{
			Delimiter:    ",",
			AllowEmpty:   false,
			FileEncoding: "raw",
		}
		processor := NewProcessor(cfg)
		_, _, err := processor.ProcessInputValues("KEY", "file:///nonexistent/path.txt")
		if err == nil {
			t.Error("expected error for missing file")
		}
	})
}

func TestSplitJSONAware(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		delimiter string
		expected  []string
	}{
		{
			name:      "Simple comma-separated",
			input:     "a,b,c",
			delimiter: ",",
			expected:  []string{"a", "b", "c"},
		},
		{
			name:      "JSON object preserved",
			input:     `{"api_url":"https://api.example.com","timeout":30},value2`,
			delimiter: ",",
			expected:  []string{`{"api_url":"https://api.example.com","timeout":30}`, "value2"},
		},
		{
			name:      "Multiple JSON objects",
			input:     `{"a":1,"b":2},{"c":3,"d":4}`,
			delimiter: ",",
			expected:  []string{`{"a":1,"b":2}`, `{"c":3,"d":4}`},
		},
		{
			name:      "Nested JSON",
			input:     `{"outer":{"inner":"val","x":1}},simple`,
			delimiter: ",",
			expected:  []string{`{"outer":{"inner":"val","x":1}}`, "simple"},
		},
		{
			name:      "JSON array preserved",
			input:     `[1,2,3],value`,
			delimiter: ",",
			expected:  []string{`[1,2,3]`, "value"},
		},
		{
			name:      "No JSON - plain split",
			input:     "a,b,c",
			delimiter: ",",
			expected:  []string{"a", "b", "c"},
		},
		{
			name:      "Single value",
			input:     `{"key":"val"}`,
			delimiter: ",",
			expected:  []string{`{"key":"val"}`},
		},
		{
			name:      "Escaped quotes in JSON",
			input:     `{"msg":"hello, \"world\""},other`,
			delimiter: ",",
			expected:  []string{`{"msg":"hello, \"world\""}`, "other"},
		},
		{
			name:      "Pipe delimiter with JSON",
			input:     `{"a":1}|value`,
			delimiter: "|",
			expected:  []string{`{"a":1}`, "value"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{Delimiter: tt.delimiter}
			processor := NewProcessor(cfg)
			result := processor.splitJSONAware(tt.input, tt.delimiter)

			if len(result) != len(tt.expected) {
				t.Fatalf("splitJSONAware() length = %d, want %d\ngot:  %v\nwant: %v", len(result), len(tt.expected), result, tt.expected)
			}
			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("splitJSONAware()[%d] = %q, want %q", i, v, tt.expected[i])
				}
			}
		})
	}
}

func TestProcessInputValuesWithJSONComma(t *testing.T) {
	t.Run("JSON value with commas preserved", func(t *testing.T) {
		cfg := &config.Config{
			Delimiter:    ",",
			AllowEmpty:   false,
			JsonSupport:  true,
			FileEncoding: "raw",
		}
		processor := NewProcessor(cfg)
		keyList, valueList, err := processor.ProcessInputValues(
			"CONFIG_JSON",
			`{"api_url":"https://api.example.com","timeout":30}`,
		)
		if err != nil {
			t.Fatalf("ProcessInputValues() error = %v", err)
		}
		// Original key + expanded nested keys
		if len(keyList) < 1 {
			t.Fatalf("expected at least 1 key, got %d", len(keyList))
		}
		if keyList[0] != "CONFIG_JSON" {
			t.Errorf("keyList[0] = %q, want %q", keyList[0], "CONFIG_JSON")
		}
		// The original JSON value should be intact
		if !strings.Contains(valueList[0], "api.example.com") {
			t.Errorf("valueList[0] should contain api.example.com, got %q", valueList[0])
		}
		_ = valueList
	})
}

func BenchmarkProcessInputValues(b *testing.B) {
	cfg := &config.Config{
		Delimiter:      ",",
		AllowEmpty:     false,
		TrimWhitespace: true,
		JsonSupport:    false,
	}

	processor := NewProcessor(cfg)
	keys := "KEY1,KEY2,KEY3,KEY4,KEY5"
	values := "VALUE1,VALUE2,VALUE3,VALUE4,VALUE5"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		processor.ProcessInputValues(keys, values)
	}
}

func BenchmarkNormalizeWhitespace(b *testing.B) {
	cfg := &config.Config{}
	processor := NewProcessor(cfg)
	input := "This is a test\nwith multiple\r\n  whitespace   issues"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		processor.normalizeWhitespace(input)
	}
}

// Test helper to clean up environment
func cleanupEnv() {
	os.Unsetenv("INPUT_ENV_KEY")
	os.Unsetenv("INPUT_ENV_VALUE")
	os.Unsetenv("INPUT_OUTPUT_KEY")
	os.Unsetenv("INPUT_OUTPUT_VALUE")
	os.Unsetenv("GITHUB_ENV")
	os.Unsetenv("GITHUB_OUTPUT")
}
