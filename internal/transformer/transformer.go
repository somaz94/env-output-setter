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
	if t.toUpper {
		value = strings.ToUpper(value)
	}
	if t.toLower {
		value = strings.ToLower(value)
	}
	if t.encodeURL {
		value = url.QueryEscape(value)
	}

	if t.escapeNewlines {
		value = strings.ReplaceAll(value, "\n", "\\n")
		value = strings.ReplaceAll(value, "\r", "\\r")
	}

	if t.maxLength > 0 && len(value) > t.maxLength {
		value = value[:t.maxLength]
	}

	return value
}

func (t *Transformer) MaskValue(value string) string {
	if !t.maskSecrets {
		return value
	}

	if t.maskPattern != nil && t.maskPattern.MatchString(value) {
		return "***"
	}

	// 기본 마스킹 로직 (예: 길이가 4 이상인 값)
	if len(value) >= 4 {
		return value[:2] + strings.Repeat("*", len(value)-2)
	}
	return "***"
}
