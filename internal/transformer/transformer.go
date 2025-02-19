package transformer

import (
	"net/url"
	"regexp"
	"strings"
)

type Transformer struct {
	maskSecrets    bool
	maskPattern    *regexp.Regexp
	toUpper        bool
	toLower        bool
	encodeURL      bool
	escapeNewlines bool
	maxLength      int
}

func New(maskSecrets bool, maskPattern string, toUpper, toLower, encodeURL bool, escapeNewlines bool, maxLength int) *Transformer {
	var pattern *regexp.Regexp
	if maskPattern != "" {
		pattern = regexp.MustCompile(maskPattern)
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

func (t *Transformer) TransformValue(value string) string {
	if value == "" {
		return value
	}

	result := value

	if t.toUpper {
		result = strings.ToUpper(result)
	} else if t.toLower {
		result = strings.ToLower(result)
	}

	if t.encodeURL {
		result = url.QueryEscape(result)
	}

	if t.escapeNewlines {
		result = strings.ReplaceAll(result, "\n", "\\n")
		result = strings.ReplaceAll(result, "\r", "\\r")
	}

	if t.maxLength > 0 && len(result) > t.maxLength {
		result = result[:t.maxLength]
	}

	return result
}

func (t *Transformer) MaskValue(value string) string {
	if !t.maskSecrets || value == "" {
		return value
	}

	if t.maskPattern != nil && t.maskPattern.MatchString(value) {
		return "***"
	}

	// Shorten the value to 4 characters
	if len(value) <= 4 {
		return "***"
	}

	visibleChars := 2
	return value[:visibleChars] + strings.Repeat("*", len(value)-visibleChars)
}
