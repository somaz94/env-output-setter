package writer

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/somaz94/env-output-setter/internal/config"
)

// Error messages for validation
const (
	errMismatchedPairs  = "env_key and env_value must have the same number of entries"
	errEmptyValue       = "empty value found for key: %s"
	errDuplicateKey     = "duplicate key found: %s"
	errValidationFailed = "validation failed for key %q: %s"
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

// ValidationRule defines a validation rule for a specific key.
type ValidationRule struct {
	Pattern       string   `json:"pattern"`        // Regex pattern the value must match
	AllowedValues []string `json:"allowed_values"` // List of allowed values
	Message       string   `json:"message"`        // Custom error message
}

// ParseValidationRules parses a JSON string into a map of validation rules.
func ParseValidationRules(rulesJSON string) (map[string]ValidationRule, error) {
	if rulesJSON == "" {
		return nil, nil
	}

	var rules map[string]ValidationRule
	if err := json.Unmarshal([]byte(rulesJSON), &rules); err != nil {
		return nil, fmt.Errorf("failed to parse validation rules: %w", err)
	}
	return rules, nil
}

// ValidateOutputs validates key-value pairs against the configured validation rules.
func (v *Validator) ValidateOutputs(keys, values []string) error {
	if v.cfg.ValidationRules == "" {
		return nil
	}

	rules, err := ParseValidationRules(v.cfg.ValidationRules)
	if err != nil {
		return err
	}

	for i, key := range keys {
		rule, exists := rules[key]
		if !exists {
			continue
		}

		value := values[i]

		// Check regex pattern
		if rule.Pattern != "" {
			matched, err := regexp.MatchString(rule.Pattern, value)
			if err != nil {
				return fmt.Errorf("invalid regex pattern for key %q: %w", key, err)
			}
			if !matched {
				msg := rule.Message
				if msg == "" {
					msg = fmt.Sprintf("value %q does not match pattern %q", value, rule.Pattern)
				}
				return fmt.Errorf(errValidationFailed, key, msg)
			}
		}

		// Check allowed values
		if len(rule.AllowedValues) > 0 {
			found := false
			for _, allowed := range rule.AllowedValues {
				if value == allowed {
					found = true
					break
				}
			}
			if !found {
				msg := rule.Message
				if msg == "" {
					msg = fmt.Sprintf("value %q is not in allowed values %v", value, rule.AllowedValues)
				}
				return fmt.Errorf(errValidationFailed, key, msg)
			}
		}
	}

	return nil
}
