package writer

import (
	"testing"

	"github.com/somaz94/env-output-setter/internal/config"
)

func TestNewValidator(t *testing.T) {
	cfg := &config.Config{
		FailOnEmpty:      true,
		TrimWhitespace:   true,
		CaseSensitive:    true,
		ErrorOnDuplicate: true,
	}

	validator := NewValidator(cfg)

	if validator == nil {
		t.Fatal("NewValidator() returned nil")
	}

	if validator.cfg != cfg {
		t.Error("NewValidator() config not set correctly")
	}
}

func TestValidatePairs(t *testing.T) {
	tests := []struct {
		name      string
		keys      []string
		values    []string
		wantError bool
	}{
		{
			name:      "Equal length pairs",
			keys:      []string{"KEY1", "KEY2", "KEY3"},
			values:    []string{"VALUE1", "VALUE2", "VALUE3"},
			wantError: false,
		},
		{
			name:      "Empty pairs",
			keys:      []string{},
			values:    []string{},
			wantError: false,
		},
		{
			name:      "Single pair",
			keys:      []string{"KEY"},
			values:    []string{"VALUE"},
			wantError: false,
		},
		{
			name:      "Mismatched - more keys",
			keys:      []string{"KEY1", "KEY2", "KEY3"},
			values:    []string{"VALUE1", "VALUE2"},
			wantError: true,
		},
		{
			name:      "Mismatched - more values",
			keys:      []string{"KEY1", "KEY2"},
			values:    []string{"VALUE1", "VALUE2", "VALUE3"},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{}
			validator := NewValidator(cfg)
			err := validator.ValidatePairs(tt.keys, tt.values)

			if tt.wantError {
				if err == nil {
					t.Error("ValidatePairs() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("ValidatePairs() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestValidateInputs(t *testing.T) {
	tests := []struct {
		name             string
		keys             []string
		values           []string
		failOnEmpty      bool
		allowEmpty       bool
		trimWhitespace   bool
		caseSensitive    bool
		errorOnDuplicate bool
		wantError        bool
		errorContains    string
	}{
		{
			name:             "Valid inputs no duplicates",
			keys:             []string{"KEY1", "KEY2", "KEY3"},
			values:           []string{"VALUE1", "VALUE2", "VALUE3"},
			failOnEmpty:      true,
			allowEmpty:       false,
			trimWhitespace:   true,
			caseSensitive:    true,
			errorOnDuplicate: true,
			wantError:        false,
		},
		{
			name:             "Empty value with failOnEmpty",
			keys:             []string{"KEY1", "KEY2"},
			values:           []string{"VALUE1", ""},
			failOnEmpty:      true,
			allowEmpty:       false,
			trimWhitespace:   true,
			caseSensitive:    true,
			errorOnDuplicate: true,
			wantError:        true,
			errorContains:    "empty value",
		},
		{
			name:             "Empty value allowed",
			keys:             []string{"KEY1", "KEY2"},
			values:           []string{"VALUE1", ""},
			failOnEmpty:      true,
			allowEmpty:       true,
			trimWhitespace:   true,
			caseSensitive:    true,
			errorOnDuplicate: true,
			wantError:        false,
		},
		{
			name:             "Duplicate keys with errorOnDuplicate",
			keys:             []string{"KEY1", "KEY2", "KEY1"},
			values:           []string{"VALUE1", "VALUE2", "VALUE3"},
			failOnEmpty:      true,
			allowEmpty:       false,
			trimWhitespace:   true,
			caseSensitive:    true,
			errorOnDuplicate: true,
			wantError:        true,
			errorContains:    "duplicate key",
		},
		{
			name:             "Duplicate keys case insensitive",
			keys:             []string{"key1", "KEY2", "KEY1"},
			values:           []string{"VALUE1", "VALUE2", "VALUE3"},
			failOnEmpty:      true,
			allowEmpty:       false,
			trimWhitespace:   true,
			caseSensitive:    false,
			errorOnDuplicate: true,
			wantError:        true,
			errorContains:    "duplicate key",
		},
		{
			name:             "Duplicate keys case sensitive - no error",
			keys:             []string{"key1", "KEY2", "KEY1"},
			values:           []string{"VALUE1", "VALUE2", "VALUE3"},
			failOnEmpty:      true,
			allowEmpty:       false,
			trimWhitespace:   true,
			caseSensitive:    true,
			errorOnDuplicate: true,
			wantError:        false,
		},
		{
			name:             "Duplicate keys allowed",
			keys:             []string{"KEY1", "KEY2", "KEY1"},
			values:           []string{"VALUE1", "VALUE2", "VALUE3"},
			failOnEmpty:      true,
			allowEmpty:       false,
			trimWhitespace:   true,
			caseSensitive:    true,
			errorOnDuplicate: false,
			wantError:        false,
		},
		{
			name:             "Whitespace trimming applied",
			keys:             []string{"  KEY1  ", "KEY2"},
			values:           []string{"VALUE1", "VALUE2"},
			failOnEmpty:      true,
			allowEmpty:       false,
			trimWhitespace:   true,
			caseSensitive:    true,
			errorOnDuplicate: true,
			wantError:        false,
		},
		{
			name:             "Empty key with failOnEmpty",
			keys:             []string{"", "KEY2"},
			values:           []string{"VALUE1", "VALUE2"},
			failOnEmpty:      true,
			allowEmpty:       false,
			trimWhitespace:   true,
			caseSensitive:    true,
			errorOnDuplicate: true,
			wantError:        true,
			errorContains:    "empty value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				FailOnEmpty:      tt.failOnEmpty,
				AllowEmpty:       tt.allowEmpty,
				TrimWhitespace:   tt.trimWhitespace,
				CaseSensitive:    tt.caseSensitive,
				ErrorOnDuplicate: tt.errorOnDuplicate,
			}

			validator := NewValidator(cfg)
			err := validator.ValidateInputs(tt.keys, tt.values)

			if tt.wantError {
				if err == nil {
					t.Error("ValidateInputs() expected error, got nil")
				} else if tt.errorContains != "" && !contains(err.Error(), tt.errorContains) {
					t.Errorf("ValidateInputs() error = %v, want to contain %v", err, tt.errorContains)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateInputs() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestValidateInputsEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		keys      []string
		values    []string
		cfg       *config.Config
		wantError bool
	}{
		{
			name:   "Empty arrays",
			keys:   []string{},
			values: []string{},
			cfg: &config.Config{
				FailOnEmpty:      true,
				ErrorOnDuplicate: true,
			},
			wantError: false,
		},
		{
			name:   "Multiple empty keys with allowEmpty",
			keys:   []string{"", "", "KEY1"},
			values: []string{"VALUE1", "VALUE2", "VALUE3"},
			cfg: &config.Config{
				FailOnEmpty:      false,
				AllowEmpty:       true,
				ErrorOnDuplicate: false, // Changed to false since we have duplicate empty keys
			},
			wantError: false,
		},
		{
			name:   "Whitespace-only keys",
			keys:   []string{"   ", "KEY1"},
			values: []string{"VALUE1", "VALUE2"},
			cfg: &config.Config{
				FailOnEmpty:    true,
				AllowEmpty:     false,
				TrimWhitespace: true,
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := NewValidator(tt.cfg)
			err := validator.ValidateInputs(tt.keys, tt.values)

			if tt.wantError {
				if err == nil {
					t.Error("ValidateInputs() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("ValidateInputs() unexpected error: %v", err)
				}
			}
		})
	}
}

func BenchmarkValidatePairs(b *testing.B) {
	cfg := &config.Config{}
	validator := NewValidator(cfg)
	keys := []string{"KEY1", "KEY2", "KEY3", "KEY4", "KEY5"}
	values := []string{"VALUE1", "VALUE2", "VALUE3", "VALUE4", "VALUE5"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator.ValidatePairs(keys, values)
	}
}

func BenchmarkValidateInputs(b *testing.B) {
	cfg := &config.Config{
		FailOnEmpty:      true,
		TrimWhitespace:   true,
		CaseSensitive:    true,
		ErrorOnDuplicate: true,
	}
	validator := NewValidator(cfg)
	keys := []string{"KEY1", "KEY2", "KEY3", "KEY4", "KEY5"}
	values := []string{"VALUE1", "VALUE2", "VALUE3", "VALUE4", "VALUE5"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator.ValidateInputs(keys, values)
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && hasSubstring(s, substr)))
}

func hasSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
