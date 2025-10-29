package writer

import (
	"strings"

	"github.com/somaz94/env-output-setter/internal/config"
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
	// Split input strings by delimiter
	keyList := strings.Split(keys, p.cfg.Delimiter)
	valueList := strings.Split(values, p.cfg.Delimiter)

	// Process keys and values for whitespace
	keyList = p.processWhitespace(keyList)
	valueList = p.processWhitespace(valueList)

	// Filter out empty entries if not allowed
	keyList = p.removeEmptyEntries(keyList)
	valueList = p.removeEmptyEntries(valueList)

	// Process JSON values if enabled
	if p.cfg.JsonSupport {
		jsonHandler := NewJSONHandler()
		keyList, valueList = jsonHandler.ProcessJSONValues(keyList, valueList)
	}

	return keyList, valueList, nil
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
// It handles newlines, carriage returns, and multiple consecutive spaces.
func (p *Processor) normalizeWhitespace(s string) string {
	// Convert line breaks to spaces
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")

	// Condense multiple spaces to a single space
	for strings.Contains(s, "  ") {
		s = strings.ReplaceAll(s, "  ", " ")
	}

	return s
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

	printer.PrintDebugSection(strings.Title(varType))
	printer.PrintDebugInfo("📥 Input Values:\n")
	printer.PrintDebugInfo("  • Keys:      %q\n", keys)
	printer.PrintDebugInfo("  • Values:    %q\n", values)
	printer.PrintDebugInfo("  • Delimiter: %q\n\n", p.cfg.Delimiter)
}

// LogProcessedValues logs the processed key-value pairs if debug mode is enabled.
func (p *Processor) LogProcessedValues(keyList, valueList []string) {
	if !p.cfg.DebugMode {
		return
	}

	printer.PrintDebugInfo("📋 Processed Values:\n")
	printer.PrintDebugInfo("  • Keys:   %v\n", keyList)
	printer.PrintDebugInfo("  • Values: %v\n\n", valueList)
}
