package writer

import (
	"strings"

	"github.com/somaz94/env-output-setter/internal/config"
	"github.com/somaz94/env-output-setter/internal/filereader"
	"github.com/somaz94/env-output-setter/internal/interpolator"
	"github.com/somaz94/env-output-setter/internal/printer"
)

// Processor handles input processing and transformation.
type Processor struct {
	cfg *config.Config
}

// NewProcessor creates a new Processor instance.
func NewProcessor(cfg *config.Config) *Processor {
	return &Processor{cfg: cfg}
}

// ProcessInputValues processes the input strings into lists with proper formatting.
// It handles splitting, trimming, and processing JSON values if json_support is enabled.
func (p *Processor) ProcessInputValues(keys, values string) ([]string, []string, error) {
	// Split input strings by delimiter (JSON-aware if json_support is enabled)
	keyList := strings.Split(keys, p.cfg.Delimiter)
	var valueList []string
	if p.cfg.JsonSupport {
		valueList = p.splitJSONAware(values, p.cfg.Delimiter)
	} else {
		valueList = strings.Split(values, p.cfg.Delimiter)
	}

	// Process keys and values for whitespace
	keyList = p.processWhitespace(keyList)
	valueList = p.processWhitespace(valueList)

	// Filter out empty entries if not allowed
	keyList = p.removeEmptyEntries(keyList)
	valueList = p.removeEmptyEntries(valueList)

	// Read values from files if any use file:// references
	fileReader := filereader.New(p.cfg.FileEncoding)
	valueList, err := fileReader.ReadValues(valueList)
	if err != nil {
		return nil, nil, err
	}

	// Interpolate variables if enabled
	if p.cfg.EnableInterpolation {
		ip := interpolator.New()
		valueList, err = ip.InterpolateList(valueList)
		if err != nil {
			return nil, nil, err
		}
	}

	// Process JSON values if enabled
	if p.cfg.JsonSupport {
		jsonHandler := NewJSONHandler()
		keyList, valueList = jsonHandler.ProcessJSONValues(keyList, valueList)
	}

	// Prepend the group prefix to every generated key name (including
	// JSON-flattened sub-keys) once the final key list is known.
	if p.cfg.GroupPrefix != "" {
		keyList = p.applyGroupPrefix(keyList)
	}

	return keyList, valueList, nil
}

// applyGroupPrefix prepends the configured group prefix and an underscore
// separator to every non-empty key (e.g. prefix "APP" turns "DATABASE" into
// "APP_DATABASE" and the JSON-flattened "CONFIG_server_host" into
// "APP_CONFIG_server_host"). Empty keys are left untouched so an allow_empty
// placeholder never becomes a bare "APP_" key. Status keys (action_status /
// error_message) are written separately by the Writer and never pass through
// here, so they keep their fixed downstream contract.
func (p *Processor) applyGroupPrefix(keys []string) []string {
	prefix := strings.TrimSpace(p.cfg.GroupPrefix)
	if prefix == "" {
		return keys
	}

	result := make([]string, len(keys))
	for i, key := range keys {
		if key == "" {
			result[i] = key
			continue
		}
		result[i] = prefix + "_" + key
	}
	return result
}

// splitJSONAware splits a string by delimiter while preserving JSON objects and arrays.
// It tracks brace/bracket depth so delimiters inside JSON are not treated as separators.
func (p *Processor) splitJSONAware(s, delimiter string) []string {
	// Fast path: no JSON markers means we can use the plain splitter.
	// Parens make the operator precedence explicit (A || (B && C)).
	if delimiter == "" || (!strings.Contains(s, "{") && !strings.Contains(s, "[")) {
		return strings.Split(s, delimiter)
	}

	var result []string
	depth := 0
	inString := false
	escaped := false
	start := 0
	delimLen := len(delimiter)

	for i := 0; i < len(s); i++ {
		ch := s[i]

		if escaped {
			escaped = false
			continue
		}

		if ch == '\\' && inString {
			escaped = true
			continue
		}

		if ch == '"' {
			inString = !inString
			continue
		}

		if inString {
			continue
		}

		if ch == '{' || ch == '[' {
			depth++
			continue
		}

		if ch == '}' || ch == ']' {
			depth--
			continue
		}

		// Only split on delimiter when at top level (depth == 0)
		if depth == 0 && i+delimLen <= len(s) && s[i:i+delimLen] == delimiter {
			result = append(result, s[start:i])
			start = i + delimLen
			i += delimLen - 1
		}
	}

	result = append(result, s[start:])
	return result
}

// processWhitespace normalizes and trims whitespace from all entries in a list.
func (p *Processor) processWhitespace(entries []string) []string {
	result := make([]string, len(entries))
	for i, entry := range entries {
		// Normalize whitespace (convert newlines to spaces, reduce multiple spaces)
		normalized := p.normalizeWhitespace(entry)
		// Trim any remaining leading/trailing whitespace
		result[i] = strings.TrimSpace(normalized)
	}
	return result
}

// normalizeWhitespace converts all whitespace sequences to a single space.
// It handles newlines, carriage returns, tabs, and multiple consecutive spaces.
func (p *Processor) normalizeWhitespace(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// removeEmptyEntries filters out empty strings from a slice.
// If allowEmpty is configured, all entries are preserved.
func (p *Processor) removeEmptyEntries(entries []string) []string {
	if p.cfg.AllowEmpty {
		return entries
	}

	result := make([]string, 0, len(entries))
	for _, entry := range entries {
		if strings.TrimSpace(entry) != "" {
			result = append(result, entry)
		}
	}
	return result
}

// LogInputValues logs the original input values if debug mode is enabled.
func (p *Processor) LogInputValues(varType, keys, values string) {
	if !p.cfg.DebugMode {
		return
	}

	printer.PrintDebugSection(titleCase(varType))
	printer.PrintDebugInfo("Input Values:\n")
	printer.PrintDebugInfo("  * Keys:      %q\n", keys)
	printer.PrintDebugInfo("  * Values:    %q\n", values)
	printer.PrintDebugInfo("  * Delimiter: %q\n\n", p.cfg.Delimiter)
}

// titleCase upper-cases the first ASCII byte of s and returns the rest unchanged.
// Returns s unchanged when empty to avoid a slice-out-of-range panic.
func titleCase(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// LogProcessedValues logs the processed key-value pairs if debug mode is enabled.
func (p *Processor) LogProcessedValues(keyList, valueList []string) {
	if !p.cfg.DebugMode {
		return
	}

	printer.PrintDebugInfo("Processed Values:\n")
	printer.PrintDebugInfo("  * Keys:   %v\n", keyList)
	printer.PrintDebugInfo("  * Values: %v\n\n", valueList)
}
