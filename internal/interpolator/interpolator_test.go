package interpolator

import (
	"testing"
)

func TestInterpolate(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		envVars   map[string]string
		expected  string
		wantError bool
		errMsg    string
	}{
		{
			name:     "No interpolation needed",
			value:    "simple_value",
			expected: "simple_value",
		},
		{
			name:     "Simple variable substitution",
			value:    "${MY_VAR}",
			envVars:  map[string]string{"MY_VAR": "hello"},
			expected: "hello",
		},
		{
			name:     "Variable not set returns empty",
			value:    "${UNSET_VAR}",
			expected: "",
		},
		{
			name:     "Default value when not set",
			value:    "${UNSET_VAR:-default_val}",
			expected: "default_val",
		},
		{
			name:     "Default value not used when set",
			value:    "${MY_VAR:-default_val}",
			envVars:  map[string]string{"MY_VAR": "actual"},
			expected: "actual",
		},
		{
			name:      "Error when required variable not set",
			value:     "${REQUIRED_VAR:?variable is required}",
			wantError: true,
			errMsg:    "variable is required",
		},
		{
			name:     "Error syntax with variable set",
			value:    "${MY_VAR:?must be set}",
			envVars:  map[string]string{"MY_VAR": "present"},
			expected: "present",
		},
		{
			name:      "Error with default message when not set",
			value:     "${MISSING:?}",
			wantError: true,
			errMsg:    "variable MISSING is not set",
		},
		{
			name:     "Multiple variables in one string",
			value:    "${HOST}:${PORT}",
			envVars:  map[string]string{"HOST": "localhost", "PORT": "8080"},
			expected: "localhost:8080",
		},
		{
			name:     "Mixed text and variables",
			value:    "Hello ${NAME:-World}! Port is ${PORT:-3000}",
			envVars:  map[string]string{"NAME": "Go"},
			expected: "Hello Go! Port is 3000",
		},
		{
			name:     "Empty string",
			value:    "",
			expected: "",
		},
		{
			name:     "No dollar sign",
			value:    "just a regular value",
			expected: "just a regular value",
		},
		{
			name:     "Partial syntax not matched",
			value:    "${incomplete",
			expected: "${incomplete",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.envVars {
				t.Setenv(k, v)
			}

			ip := New()
			result, err := ip.Interpolate(tt.value)

			if tt.wantError {
				if err == nil {
					t.Error("Interpolate() expected error, got nil")
				} else if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
					t.Errorf("Interpolate() error = %v, want to contain %q", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Interpolate() unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("Interpolate() = %q, want %q", result, tt.expected)
				}
			}
		})
	}
}

func TestInterpolateList(t *testing.T) {
	tests := []struct {
		name      string
		values    []string
		envVars   map[string]string
		expected  []string
		wantError bool
	}{
		{
			name:     "Multiple values with interpolation",
			values:   []string{"${HOST}", "${PORT:-8080}", "static"},
			envVars:  map[string]string{"HOST": "localhost"},
			expected: []string{"localhost", "8080", "static"},
		},
		{
			name:      "Error propagated from single value",
			values:    []string{"ok", "${MISSING:?required}"},
			wantError: true,
		},
		{
			name:     "Empty list",
			values:   []string{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.envVars {
				t.Setenv(k, v)
			}

			ip := New()
			result, err := ip.InterpolateList(tt.values)

			if tt.wantError {
				if err == nil {
					t.Error("InterpolateList() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("InterpolateList() unexpected error: %v", err)
				}
				if len(result) != len(tt.expected) {
					t.Fatalf("InterpolateList() length = %d, want %d", len(result), len(tt.expected))
				}
				for i, v := range result {
					if v != tt.expected[i] {
						t.Errorf("InterpolateList()[%d] = %q, want %q", i, v, tt.expected[i])
					}
				}
			}
		})
	}
}

func BenchmarkInterpolate(b *testing.B) {
	b.Setenv("HOST", "localhost")
	b.Setenv("PORT", "8080")

	ip := New()
	value := "http://${HOST}:${PORT}/api/${VERSION:-v1}"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ip.Interpolate(value)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
