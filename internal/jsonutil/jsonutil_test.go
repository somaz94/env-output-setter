package jsonutil

import "testing"

func TestIsJSONLike(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{name: "Valid JSON object", value: `{"key":"value"}`, expected: true},
		{name: "Valid JSON array", value: `["item1","item2"]`, expected: true},
		{name: "JSON with whitespace", value: `  {"key":"value"}  `, expected: true},
		{name: "Nested JSON object", value: `{"a":{"b":"c"}}`, expected: true},
		{name: "Empty JSON object", value: `{}`, expected: true},
		{name: "Empty JSON array", value: `[]`, expected: true},
		{name: "Plain text", value: "not json", expected: false},
		{name: "Incomplete JSON", value: `{"key":`, expected: false},
		{name: "Empty string", value: "", expected: false},
		{name: "Just opening brace", value: "{", expected: false},
		{name: "Just closing brace", value: "}", expected: false},
		{name: "Mismatched braces", value: "{]", expected: false},
		{name: "Mismatched brackets", value: "[}", expected: false},
		{name: "Number", value: "42", expected: false},
		{name: "Boolean", value: "true", expected: false},
		{name: "Whitespace only", value: "   ", expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsJSONLike(tt.value)
			if result != tt.expected {
				t.Errorf("IsJSONLike(%q) = %v, want %v", tt.value, result, tt.expected)
			}
		})
	}
}

func BenchmarkIsJSONLike(b *testing.B) {
	value := `{"key":"value","nested":{"inner":"data"}}`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IsJSONLike(value)
	}
}
