package writer

import (
	"fmt"
	"strings"

	"github.com/somaz94/env-output-setter/internal/config"
)

// Error messages for validation
const (
	errMismatchedPairs = "env_key and env_value must have the same number of entries"
	errEmptyValue      = "empty value found for key: %s"
	errDuplicateKey    = "duplicate key found: %s"
)

// Validator handles input validation logic.
type Validator struct {
	cfg *config.Config
}

// NewValidator creates a new Validator instance.
func NewValidator(cfg *config.Config) *Validator {
	return &Validator{cfg: cfg}
}

// ValidatePairs ensures key-value pairs match in count.
func (v *Validator) ValidatePairs(keys, values []string) error {
	if len(keys) != len(values) {
		return fmt.Errorf("%s (keys: %d, values: %d)",
			errMismatchedPairs, len(keys), len(values))
	}
	return nil
}

// ValidateInputs checks for empty values and duplicate keys based on configuration.
func (v *Validator) ValidateInputs(keys, values []string) error {
	seenKeys := make(map[string]bool)

	for i, key := range keys {
		// Apply trimming if configured
		if v.cfg.TrimWhitespace {
			key = strings.TrimSpace(key)
			keys[i] = key
		}

		// Prepare key for duplicate checking
		lookupKey := key
		if !v.cfg.CaseSensitive {
			lookupKey = strings.ToLower(key)
		}

		// Check for empty values if configured to fail
		if v.cfg.FailOnEmpty && !v.cfg.AllowEmpty && (key == "" || values[i] == "") {
			return fmt.Errorf(errEmptyValue, key)
		}

		// Check for duplicate keys if configured
		if v.cfg.ErrorOnDuplicate {
			if seenKeys[lookupKey] {
				return fmt.Errorf(errDuplicateKey, key)
			}
			seenKeys[lookupKey] = true
		}
	}

	return nil
}
