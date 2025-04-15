package transformer

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// TransformationError represents an error that occurs during value transformation.
// It includes both an error message and an optional underlying cause.
type TransformationError struct {
	Message string
	Cause   error
}

// Error implements the error interface for TransformationError.
// It returns a formatted error string that includes the underlying cause if present.
func (e *TransformationError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// Transformer handles various value transformations including casing,
// encoding, masking and length limitations.
type Transformer struct {
	// Masking settings
	maskSecrets bool
	maskPattern *regexp.Regexp

	// Case conversion settings
	toUpper bool
	toLower bool

	// Encoding settings
	encodeURL      bool
	escapeNewlines bool

	// Length limitation settings
	maxLength int
}

// New creates a new Transformer with the specified configuration options.
// It handles the compilation of regular expression patterns for masking.
func New(
	maskSecrets bool,
	maskPattern string,
	toUpper bool,
	toLower bool,
	encodeURL bool,
	escapeNewlines bool,
	maxLength int,
) *Transformer {
	var pattern *regexp.Regexp
	if maskPattern != "" {
		var err error
		pattern, err = regexp.Compile(maskPattern)
		if err != nil {
			fmt.Printf("Warning: Invalid mask pattern '%s': %v\n", maskPattern, err)
		}
	}

	return &Transformer{
		maskSecrets:    maskSecrets,
		maskPattern:    pattern,
		toUpper:        toUpper,
		toLower:        toLower,
		encodeURL:      encodeURL,
		escapeNewlines: escapeNewlines,
		maxLength:      maxLength,
	}
}

// TransformValue applies all configured transformations to a value in the following order:
// 1. Case conversion (upper/lower)
// 2. URL encoding
// 3. Newline escaping
// 4. Length limitation
//
// If JSON support is enabled and the value looks like JSON, it preserves the JSON format.
func (t *Transformer) TransformValue(value string, supportJSON bool) string {
	// Handle empty values early
	if value == "" {
		return value
	}

	// Check if value is JSON and JSON support is enabled
	if supportJSON && isJSONValue(value) {
		return t.handleJSONValue(value)
	}

	// Apply transformations in sequence
	result := value

	// 1. Apply case conversion (mutually exclusive)
	result = t.applyCaseConversion(result)

	// 2. Apply URL encoding if enabled
	if t.encodeURL {
		result = url.QueryEscape(result)
	}

	// 3. Escape newlines if enabled
	if t.escapeNewlines {
		result = t.escapeNewlineCharacters(result)
	}

	// 4. Apply length limitation if configured
	if t.maxLength > 0 && len(result) > t.maxLength {
		result = result[:t.maxLength]
	}

	return result
}

// isJSONValue checks if a string appears to be JSON (object or array).
func isJSONValue(value string) bool {
	value = strings.TrimSpace(value)
	return (strings.HasPrefix(value, "{") && strings.HasSuffix(value, "}")) ||
		(strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]"))
}

// handleJSONValue processes a value that appears to be JSON.
// It validates the JSON and returns it unchanged if valid.
func (t *Transformer) handleJSONValue(value string) string {
	// Verify the JSON is valid
	var jsonObj interface{}
	if err := json.Unmarshal([]byte(value), &jsonObj); err == nil {
		// Valid JSON is returned as-is
		return value
	}

	// If JSON is invalid, proceed with normal transformations
	result := value
	result = t.applyCaseConversion(result)

	if t.encodeURL {
		result = url.QueryEscape(result)
	}

	if t.escapeNewlines {
		result = t.escapeNewlineCharacters(result)
	}

	if t.maxLength > 0 && len(result) > t.maxLength {
		result = result[:t.maxLength]
	}

	return result
}

// applyCaseConversion applies upper or lower case conversion if enabled.
// Upper case takes precedence over lower case if both are somehow enabled.
func (t *Transformer) applyCaseConversion(value string) string {
	if t.toUpper {
		return strings.ToUpper(value)
	} else if t.toLower {
		return strings.ToLower(value)
	}
	return value
}

// escapeNewlineCharacters replaces newline characters with their escaped equivalents.
func (t *Transformer) escapeNewlineCharacters(value string) string {
	result := strings.ReplaceAll(value, "\n", "\\n")
	result = strings.ReplaceAll(result, "\r", "\\r")
	return result
}

// MaskValue applies masking to sensitive values to hide their content.
// It uses different masking strategies based on the configuration and value length.
func (t *Transformer) MaskValue(value string) string {
	// Skip masking if disabled or value is empty
	if !t.maskSecrets || value == "" {
		return value
	}

	// Apply regex pattern masking if configured and matched
	if t.maskPattern != nil && t.maskPattern.MatchString(value) {
		return "***"
	}

	// Apply full masking for short values
	if len(value) <= 4 {
		return "***"
	}

	// Default masking: show first 2 characters and mask the rest
	const visiblePrefix = 2
	return value[:visiblePrefix] + strings.Repeat("*", len(value)-visiblePrefix)
}

// CustomMask applies a custom masking pattern with configurable visible
// prefix and suffix lengths. This allows for more precise control over
// what parts of a value remain visible.
func (t *Transformer) CustomMask(value string, visiblePrefix, visibleSuffix int) string {
	// Handle empty strings or values too short for the requested visibility
	if value == "" || len(value) <= (visiblePrefix+visibleSuffix) {
		return "***"
	}

	// Extract visible prefix if requested
	prefix := ""
	if visiblePrefix > 0 {
		prefix = value[:visiblePrefix]
	}

	// Extract visible suffix if requested
	suffix := ""
	if visibleSuffix > 0 {
		suffix = value[len(value)-visibleSuffix:]
	}

	// Create mask string for the middle portion
	maskedLength := len(value) - visiblePrefix - visibleSuffix
	maskString := strings.Repeat("*", maskedLength)

	// Combine the parts
	return prefix + maskString + suffix
}

// TransformJSON converts a JSON string to a compact form by removing
// unnecessary whitespace. It also validates that the input is valid JSON.
func (t *Transformer) TransformJSON(value string) (string, error) {
	// Parse the JSON to validate and normalize it
	var data interface{}
	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return "", &TransformationError{
			Message: "Invalid JSON format",
			Cause:   err,
		}
	}

	// Re-marshal to get a compact representation
	compact, err := json.Marshal(data)
	if err != nil {
		return "", &TransformationError{
			Message: "Failed to compact JSON",
			Cause:   err,
		}
	}

	return string(compact), nil
}
