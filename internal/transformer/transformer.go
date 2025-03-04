package transformer

import (
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
	// 마스킹 관련 설정
	maskSecrets bool
	maskPattern *regexp.Regexp

	// 대소문자 변환 설정
	toUpper bool
	toLower bool

	// 인코딩 설정
	encodeURL      bool
	escapeNewlines bool

	// 길이 제한 설정
	maxLength int
}

// New creates a new Transformer with the specified options
func New(maskSecrets bool, maskPattern string, toUpper, toLower, encodeURL bool, escapeNewlines bool, maxLength int) *Transformer {
	var pattern *regexp.Regexp
	if maskPattern != "" {
		var err error
		pattern, err = regexp.Compile(maskPattern)
		if err != nil {
			// 정규식 컴파일 오류 시 로그만 출력하고 nil 패턴 사용
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
// 변환 순서:
// 1. 대소문자 변환 (toUpper/toLower)
// 2. URL 인코딩 (encodeURL)
// 3. 줄바꿈 이스케이프 (escapeNewlines)
// 4. 길이 제한 (maxLength)
func (t *Transformer) TransformValue(value string) string {
	if value == "" {
		return value
	}

	result := value

	// 1. 대소문자 변환 (상호 배타적)
	if t.toUpper {
		result = strings.ToUpper(result)
	} else if t.toLower {
		result = strings.ToLower(result)
	}

	// 2. URL 인코딩
	if t.encodeURL {
		result = url.QueryEscape(result)
	}

	// 3. 줄바꿈 이스케이프
	if t.escapeNewlines {
		result = strings.ReplaceAll(result, "\n", "\\n")
		result = strings.ReplaceAll(result, "\r", "\\r")
	}

	// 4. 길이 제한
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

	// 정규식 패턴이 있고 매치되면 전체 마스킹
	if t.maskPattern != nil && t.maskPattern.MatchString(value) {
		return "***"
	}

	// 짧은 값은 전체 마스킹
	if len(value) <= 4 {
		return "***"
	}

	// 기본 마스킹: 앞 2자리만 표시하고 나머지 마스킹
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
