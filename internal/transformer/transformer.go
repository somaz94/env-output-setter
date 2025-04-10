package transformer

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// TransformationError represents an error during transformation
type TransformationError struct {
	Message string
	Cause   error
}

// Error implements the error interface
func (e *TransformationError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// Transformer handles value transformations
type Transformer struct {
	// Masking Related Settings
	maskSecrets bool
	maskPattern *regexp.Regexp

	// Case Conversion Settings
	toUpper bool
	toLower bool

	// Encoding Settings
	encodeURL      bool
	escapeNewlines bool

	// Length Limitation Settings
	maxLength int
}

// New creates a new Transformer with the specified options
func New(maskSecrets bool, maskPattern string, toUpper, toLower, encodeURL bool, escapeNewlines bool, maxLength int) *Transformer {
	var pattern *regexp.Regexp
	if maskPattern != "" {
		var err error
		pattern, err = regexp.Compile(maskPattern)
		if err != nil {
			// Log only if there is an error in the regular expression compilation
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

// TransformValue applies all configured transformations to a value
// Transformation Order:
// 1. Case Conversion (toUpper/toLower)
// 2. URL Encoding (encodeURL)
// 3. Escape Newlines (escapeNewlines)
// 4. Length Limitation (maxLength)
func (t *Transformer) TransformValue(value string, supportJson bool) string {
	if value == "" {
		return value
	}

	result := value

	// JSON인 경우 특별한 처리를 건너뛰기
	if supportJson && (strings.HasPrefix(result, "{") || strings.HasPrefix(result, "[")) {
		// JSON 형식 검증
		var jsonObj interface{}
		if err := json.Unmarshal([]byte(result), &jsonObj); err == nil {
			// 유효한 JSON이면 그대로 반환
			return result
		}
		// 유효하지 않은 JSON이지만 JSON 형식을 가장하는 경우
		// 특별한 처리 없이 계속 진행 (프로세스 중에 실패 처리)
	}

	// 1. Case Conversion (mutually exclusive)
	if t.toUpper {
		result = strings.ToUpper(result)
	} else if t.toLower {
		result = strings.ToLower(result)
	}

	// 2. URL Encoding
	if t.encodeURL {
		result = url.QueryEscape(result)
	}

	// 3. Escape Newlines
	if t.escapeNewlines {
		result = strings.ReplaceAll(result, "\n", "\\n")
		result = strings.ReplaceAll(result, "\r", "\\r")
	}

	// 4. Length Limitation
	if t.maxLength > 0 && len(result) > t.maxLength {
		result = result[:t.maxLength]
	}

	return result
}

// MaskValue applies masking to sensitive values
func (t *Transformer) MaskValue(value string) string {
	if !t.maskSecrets || value == "" {
		return value
	}

	// If there is a regular expression pattern and it matches, mask the entire value
	if t.maskPattern != nil && t.maskPattern.MatchString(value) {
		return "***"
	}

	// Short values are fully masked
	if len(value) <= 4 {
		return "***"
	}

	// Default masking: show the first 2 characters and mask the rest
	visibleChars := 2
	return value[:visibleChars] + strings.Repeat("*", len(value)-visibleChars)
}

// CustomMask applies a custom masking pattern
func (t *Transformer) CustomMask(value string, visiblePrefix, visibleSuffix int) string {
	if value == "" || len(value) <= (visiblePrefix+visibleSuffix) {
		return "***"
	}

	prefix := ""
	if visiblePrefix > 0 {
		prefix = value[:visiblePrefix]
	}

	suffix := ""
	if visibleSuffix > 0 {
		suffix = value[len(value)-visibleSuffix:]
	}

	maskedLength := len(value) - visiblePrefix - visibleSuffix
	return prefix + strings.Repeat("*", maskedLength) + suffix
}

// TransformJSON converts a JSON string to a compact JSON string
func (t *Transformer) TransformJSON(value string) (string, error) {
	var data interface{}
	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return "", &TransformationError{
			Message: "Invalid JSON format",
			Cause:   err,
		}
	}

	// JSON 압축 (공백 제거)
	compact, err := json.Marshal(data)
	if err != nil {
		return "", &TransformationError{
			Message: "Failed to compact JSON",
			Cause:   err,
		}
	}

	return string(compact), nil
}
