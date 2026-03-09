package transformer

import (
	"errors"
	"strings"
	"testing"

	"github.com/somaz94/env-output-setter/internal/jsonutil"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name           string
		maskSecrets    bool
		maskPattern    string
		toUpper        bool
		toLower        bool
		encodeURL      bool
		escapeNewlines bool
		maxLength      int
		wantPattern    bool
	}{
		{
			name:           "Default transformer without masking",
			maskSecrets:    false,
			maskPattern:    "",
			toUpper:        false,
			toLower:        false,
			encodeURL:      false,
			escapeNewlines: false,
			maxLength:      0,
			wantPattern:    false,
		},
		{
			name:           "Transformer with masking enabled",
			maskSecrets:    true,
			maskPattern:    "^secret_",
			toUpper:        false,
			toLower:        false,
			encodeURL:      false,
			escapeNewlines: true,
			maxLength:      0,
			wantPattern:    true,
		},
		{
			name:           "Transformer with all features enabled",
			maskSecrets:    true,
			maskPattern:    "api.*key",
			toUpper:        true,
			toLower:        false,
			encodeURL:      true,
			escapeNewlines: true,
			maxLength:      100,
			wantPattern:    true,
		},
		{
			name:           "Invalid regex pattern",
			maskSecrets:    true,
			maskPattern:    "[invalid(",
			toUpper:        false,
			toLower:        false,
			encodeURL:      false,
			escapeNewlines: false,
			maxLength:      0,
			wantPattern:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := New(tt.maskSecrets, tt.maskPattern, tt.toUpper, tt.toLower, tt.encodeURL, tt.escapeNewlines, tt.maxLength)

			if tr == nil {
				t.Fatal("New() returned nil")
			}

			if tr.maskSecrets != tt.maskSecrets {
				t.Errorf("maskSecrets = %v, want %v", tr.maskSecrets, tt.maskSecrets)
			}

			if tt.wantPattern && tr.maskPattern == nil {
				t.Error("maskPattern should not be nil")
			}

			if !tt.wantPattern && tr.maskPattern != nil {
				t.Error("maskPattern should be nil for invalid patterns")
			}
		})
	}
}

func TestTransformValue(t *testing.T) {
	tests := []struct {
		name           string
		value          string
		supportJSON    bool
		toUpper        bool
		toLower        bool
		encodeURL      bool
		escapeNewlines bool
		maxLength      int
		expected       string
	}{
		{
			name:           "Empty value",
			value:          "",
			supportJSON:    false,
			toUpper:        false,
			toLower:        false,
			encodeURL:      false,
			escapeNewlines: false,
			maxLength:      0,
			expected:       "",
		},
		{
			name:           "Simple value no transformation",
			value:          "hello world",
			supportJSON:    false,
			toUpper:        false,
			toLower:        false,
			encodeURL:      false,
			escapeNewlines: false,
			maxLength:      0,
			expected:       "hello world",
		},
		{
			name:           "Convert to uppercase",
			value:          "hello world",
			supportJSON:    false,
			toUpper:        true,
			toLower:        false,
			encodeURL:      false,
			escapeNewlines: false,
			maxLength:      0,
			expected:       "HELLO WORLD",
		},
		{
			name:           "Convert to lowercase",
			value:          "HELLO WORLD",
			supportJSON:    false,
			toUpper:        false,
			toLower:        true,
			encodeURL:      false,
			escapeNewlines: false,
			maxLength:      0,
			expected:       "hello world",
		},
		{
			name:           "URL encoding",
			value:          "hello world@example.com",
			supportJSON:    false,
			toUpper:        false,
			toLower:        false,
			encodeURL:      true,
			escapeNewlines: false,
			maxLength:      0,
			expected:       "hello+world%40example.com",
		},
		{
			name:           "Escape newlines",
			value:          "line1\nline2\rline3",
			supportJSON:    false,
			toUpper:        false,
			toLower:        false,
			encodeURL:      false,
			escapeNewlines: true,
			maxLength:      0,
			expected:       "line1\\nline2\\rline3",
		},
		{
			name:           "Max length truncation",
			value:          "this is a very long string",
			supportJSON:    false,
			toUpper:        false,
			toLower:        false,
			encodeURL:      false,
			escapeNewlines: false,
			maxLength:      10,
			expected:       "this is a ",
		},
		{
			name:           "Combined transformations",
			value:          "hello\nworld",
			supportJSON:    false,
			toUpper:        true,
			toLower:        false,
			encodeURL:      false,
			escapeNewlines: true,
			maxLength:      0,
			expected:       "HELLO\\nWORLD",
		},
		{
			name:           "JSON object preserved",
			value:          `{"key":"value"}`,
			supportJSON:    true,
			toUpper:        false,
			toLower:        false,
			encodeURL:      false,
			escapeNewlines: false,
			maxLength:      0,
			expected:       `{"key":"value"}`,
		},
		{
			name:           "JSON array preserved",
			value:          `["item1","item2"]`,
			supportJSON:    true,
			toUpper:        false,
			toLower:        false,
			encodeURL:      false,
			escapeNewlines: false,
			maxLength:      0,
			expected:       `["item1","item2"]`,
		},
		{
			name:           "Invalid JSON with transformations",
			value:          `{invalid}`,
			supportJSON:    true,
			toUpper:        true,
			toLower:        false,
			encodeURL:      false,
			escapeNewlines: false,
			maxLength:      0,
			expected:       `{INVALID}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := New(false, "", tt.toUpper, tt.toLower, tt.encodeURL, tt.escapeNewlines, tt.maxLength)
			result := tr.TransformValue(tt.value, tt.supportJSON)

			if result != tt.expected {
				t.Errorf("TransformValue() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMaskValue(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		maskSecrets bool
		maskPattern string
		expected    string
	}{
		{
			name:        "Masking disabled",
			value:       "secret123",
			maskSecrets: false,
			maskPattern: "",
			expected:    "secret123",
		},
		{
			name:        "Empty value",
			value:       "",
			maskSecrets: true,
			maskPattern: "",
			expected:    "",
		},
		{
			name:        "Short value (4 chars or less)",
			value:       "abc",
			maskSecrets: true,
			maskPattern: "",
			expected:    "***",
		},
		{
			name:        "Short value exactly 4 chars",
			value:       "abcd",
			maskSecrets: true,
			maskPattern: "",
			expected:    "***",
		},
		{
			name:        "Default masking for longer value",
			value:       "secret123",
			maskSecrets: true,
			maskPattern: "",
			expected:    "se*******",
		},
		{
			name:        "Pattern matching",
			value:       "secret_key_value",
			maskSecrets: true,
			maskPattern: "^secret_",
			expected:    "***",
		},
		{
			name:        "Pattern not matching",
			value:       "public_key_value",
			maskSecrets: true,
			maskPattern: "^secret_",
			expected:    "pu**************",
		},
		{
			name:        "API key masking",
			value:       "api_key_1234567890",
			maskSecrets: true,
			maskPattern: "api.*key",
			expected:    "***",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := New(tt.maskSecrets, tt.maskPattern, false, false, false, false, 0)
			result := tr.MaskValue(tt.value)

			if result != tt.expected {
				t.Errorf("MaskValue() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCustomMask(t *testing.T) {
	tests := []struct {
		name          string
		value         string
		visiblePrefix int
		visibleSuffix int
		expected      string
	}{
		{
			name:          "Empty value",
			value:         "",
			visiblePrefix: 2,
			visibleSuffix: 2,
			expected:      "***",
		},
		{
			name:          "Value too short",
			value:         "abc",
			visiblePrefix: 2,
			visibleSuffix: 2,
			expected:      "***",
		},
		{
			name:          "Show prefix only",
			value:         "1234567890",
			visiblePrefix: 3,
			visibleSuffix: 0,
			expected:      "123*******",
		},
		{
			name:          "Show suffix only",
			value:         "1234567890",
			visiblePrefix: 0,
			visibleSuffix: 3,
			expected:      "*******890",
		},
		{
			name:          "Show both prefix and suffix",
			value:         "1234567890",
			visiblePrefix: 2,
			visibleSuffix: 2,
			expected:      "12******90",
		},
		{
			name:          "Credit card style masking",
			value:         "1234-5678-9012-3456",
			visiblePrefix: 0,
			visibleSuffix: 4,
			expected:      "***************3456",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := New(false, "", false, false, false, false, 0)
			result := tr.CustomMask(tt.value, tt.visiblePrefix, tt.visibleSuffix)

			if result != tt.expected {
				t.Errorf("CustomMask() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTransformJSON(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		wantErr   bool
		wantEmpty bool
	}{
		{
			name:      "Valid JSON object",
			value:     `{"key": "value", "number": 123}`,
			wantErr:   false,
			wantEmpty: false,
		},
		{
			name:      "Valid JSON array",
			value:     `["item1", "item2", "item3"]`,
			wantErr:   false,
			wantEmpty: false,
		},
		{
			name:      "Valid nested JSON",
			value:     `{"outer": {"inner": "value"}}`,
			wantErr:   false,
			wantEmpty: false,
		},
		{
			name:      "Invalid JSON",
			value:     `{invalid json}`,
			wantErr:   true,
			wantEmpty: false,
		},
		{
			name:      "Empty JSON object",
			value:     `{}`,
			wantErr:   false,
			wantEmpty: false,
		},
		{
			name:      "Empty JSON array",
			value:     `[]`,
			wantErr:   false,
			wantEmpty: false,
		},
		{
			name:      "JSON with whitespace",
			value:     `  {  "key"  :  "value"  }  `,
			wantErr:   false,
			wantEmpty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := New(false, "", false, false, false, false, 0)
			result, err := tr.TransformJSON(tt.value)

			if tt.wantErr {
				if err == nil {
					t.Error("TransformJSON() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("TransformJSON() unexpected error: %v", err)
				}
				if !tt.wantEmpty && result == "" {
					t.Error("TransformJSON() returned empty string unexpectedly")
				}
			}
		})
	}
}

func TestTransformationError(t *testing.T) {
	tests := []struct {
		name     string
		err      *TransformationError
		expected string
	}{
		{
			name: "Error without cause",
			err: &TransformationError{
				Message: "Invalid JSON format",
				Cause:   nil,
			},
			expected: "Invalid JSON format",
		},
		{
			name: "Error with cause",
			err: &TransformationError{
				Message: "Failed to parse",
				Cause:   errors.New("underlying error"),
			},
			expected: "Failed to parse:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.Error()
			if !strings.Contains(result, tt.expected) {
				t.Errorf("Error() = %v, want to contain %v", result, tt.expected)
			}
		})
	}
}

func TestApplyCaseConversion(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		toUpper  bool
		toLower  bool
		expected string
	}{
		{
			name:     "No conversion",
			value:    "Hello World",
			toUpper:  false,
			toLower:  false,
			expected: "Hello World",
		},
		{
			name:     "To uppercase",
			value:    "hello world",
			toUpper:  true,
			toLower:  false,
			expected: "HELLO WORLD",
		},
		{
			name:     "To lowercase",
			value:    "HELLO WORLD",
			toUpper:  false,
			toLower:  true,
			expected: "hello world",
		},
		{
			name:     "Upper takes precedence",
			value:    "Hello World",
			toUpper:  true,
			toLower:  true,
			expected: "HELLO WORLD",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := New(false, "", tt.toUpper, tt.toLower, false, false, 0)
			result := tr.applyCaseConversion(tt.value)

			if result != tt.expected {
				t.Errorf("applyCaseConversion() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestEscapeNewlineCharacters(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{
			name:     "No newlines",
			value:    "simple text",
			expected: "simple text",
		},
		{
			name:     "Unix newline",
			value:    "line1\nline2",
			expected: "line1\\nline2",
		},
		{
			name:     "Windows newline",
			value:    "line1\r\nline2",
			expected: "line1\\r\\nline2",
		},
		{
			name:     "Multiple newlines",
			value:    "line1\n\nline2\nline3",
			expected: "line1\\n\\nline2\\nline3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := New(false, "", false, false, false, false, 0)
			result := tr.escapeNewlineCharacters(tt.value)

			if result != tt.expected {
				t.Errorf("escapeNewlineCharacters() = %v, want %v", result, tt.expected)
			}
		})
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
			value:    `["item"]`,
			expected: true,
		},
		{
			name:     "JSON object with whitespace",
			value:    `  {"key":"value"}  `,
			expected: true,
		},
		{
			name:     "Not JSON",
			value:    "plain text",
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

func TestHandleJSONValueInvalidWithAllTransforms(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		toUpper  bool
		encode   bool
		escape   bool
		maxLen   int
		expected string
	}{
		{
			name:     "Invalid JSON with URL encoding",
			value:    `{invalid json}`,
			encode:   true,
			expected: "%7Binvalid+json%7D",
		},
		{
			name:     "Invalid JSON with newline escaping",
			value:    "{invalid\njson}",
			escape:   true,
			expected: "{invalid\\njson}",
		},
		{
			name:     "Invalid JSON with max length",
			value:    `{invalid json value here}`,
			maxLen:   10,
			expected: `{invalid j`,
		},
		{
			name:     "Invalid JSON with all transforms",
			value:    `{bad}`,
			toUpper:  true,
			encode:   true,
			escape:   true,
			maxLen:   20,
			expected: "%7BBAD%7D",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := New(false, "", tt.toUpper, false, tt.encode, tt.escape, tt.maxLen)
			result := tr.TransformValue(tt.value, true)
			if result != tt.expected {
				t.Errorf("TransformValue() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func BenchmarkTransformValue(b *testing.B) {
	tr := New(false, "", false, false, false, true, 0)
	value := "This is a test value\nwith newlines\rand other content"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr.TransformValue(value, false)
	}
}

func BenchmarkMaskValue(b *testing.B) {
	tr := New(true, "^secret_", false, false, false, false, 0)
	value := "secret_api_key_1234567890"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr.MaskValue(value)
	}
}

func BenchmarkTransformJSON(b *testing.B) {
	tr := New(false, "", false, false, false, false, 0)
	value := `{"key": "value", "nested": {"inner": "data"}}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr.TransformJSON(value)
	}
}
