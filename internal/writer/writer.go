package writer

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/somaz94/env-output-setter/internal/config"
	"github.com/somaz94/env-output-setter/internal/printer"
	"github.com/somaz94/env-output-setter/internal/transformer"
)

const (
	errMismatchedPairs = "env_key and env_value must have the same number of entries"
	errWriteFile       = "failed to write to %s file"
	errEmptyValue      = "empty value found for key: %s"
	errDuplicateKey    = "duplicate key found: %s"
	errMaxRetries      = "failed to write after %d retries"
	localExecMsg       = "Local Execution - %s is not set, skipping writing to GitHub Actions %s"
)

// SetEnv sets environment variables in GitHub Actions environment file
func SetEnv(cfg *config.Config) (int, error) {
	return setVariables(cfg, "GITHUB_ENV", "env")
}

// SetOutput sets output variables in GitHub Actions output file
func SetOutput(cfg *config.Config) (int, error) {
	return setVariables(cfg, "GITHUB_OUTPUT", "output")
}

// setVariables handles setting variables for both env and output files
func setVariables(cfg *config.Config, envVar, varType string) (int, error) {
	// Get input value
	keys, values := getInputValues(cfg, envVar)

	// Debug logging - original input values
	logInputValues(cfg, varType, keys, values)

	// Process input values
	keyList, valueList, err := processInputValues(cfg, keys, values)
	if err != nil {
		return 0, err
	}

	// Debug logging - processed input values
	logProcessedValues(cfg, keyList, valueList)

	// Validate input values
	if err := validateInputs(cfg, keyList, valueList); err != nil {
		return 0, err
	}

	// Check file path and process
	filePath := os.Getenv(envVar)
	if filePath == "" {
		return handleLocalExecution(envVar, varType, keyList, valueList)
	}

	// Write to file
	return writeToFile(cfg, filePath, keyList, valueList, varType)
}

// getInputValues returns the appropriate keys and values based on the environment variable
func getInputValues(cfg *config.Config, envVar string) (string, string) {
	switch envVar {
	case "GITHUB_ENV":
		return cfg.EnvKeys, cfg.EnvValues
	case "GITHUB_OUTPUT":
		return cfg.OutputKeys, cfg.OutputValues
	default:
		return "", ""
	}
}

// logInputValues logs the original input values if debug mode is enabled
func logInputValues(cfg *config.Config, varType, keys, values string) {
	if cfg.DebugMode {
		printer.PrintDebugSection(strings.Title(varType))
		printer.PrintDebugInfo("📥 Input Values:\n")
		printer.PrintDebugInfo("  • Keys:      %q\n", keys)
		printer.PrintDebugInfo("  • Values:    %q\n", values)
		printer.PrintDebugInfo("  • Delimiter: %q\n\n", cfg.Delimiter)
	}
}

// processInputValues processes the input strings into lists with proper formatting
func processInputValues(cfg *config.Config, keys, values string) ([]string, []string, error) {
	// Split by delimiter first
	keyList := strings.Split(keys, cfg.Delimiter)
	valueList := strings.Split(values, cfg.Delimiter)

	// Normalize whitespace for each item
	for i := range keyList {
		keyList[i] = normalizeWhitespace(keyList[i])
	}
	for i := range valueList {
		valueList[i] = normalizeWhitespace(valueList[i])
	}

	// Remove whitespace from each item
	for i := range keyList {
		keyList[i] = strings.TrimSpace(keyList[i])
	}
	for i := range valueList {
		valueList[i] = strings.TrimSpace(valueList[i])
	}

	// Remove empty items (do not remove if allow_empty is true)
	keyList = removeEmptyEntries(keyList, cfg.AllowEmpty)
	valueList = removeEmptyEntries(valueList, cfg.AllowEmpty)

	// Check if the number of keys and values match
	if len(keyList) != len(valueList) {
		return nil, nil, fmt.Errorf("%s (keys: %d, values: %d)", errMismatchedPairs, len(keyList), len(valueList))
	}

	return keyList, valueList, nil
}

// logProcessedValues logs the processed values if debug mode is enabled
func logProcessedValues(cfg *config.Config, keyList, valueList []string) {
	if cfg.DebugMode {
		printer.PrintDebugInfo("📋 Processed Values:\n")
		printer.PrintDebugInfo("  • Keys:   %v\n", keyList)
		printer.PrintDebugInfo("  • Values: %v\n\n", valueList)
	}
}

// handleLocalExecution handles the case when running outside of GitHub Actions
func handleLocalExecution(envVar, varType string, keyList, valueList []string) (int, error) {
	fmt.Printf(localExecMsg, envVar, varType)
	for i, key := range keyList {
		printer.PrintSuccess(varType, key, valueList[i])
	}
	return len(keyList), nil
}

// normalizeWhitespace normalizes all whitespace including newlines
func normalizeWhitespace(s string) string {
	// Convert actual line breaks to spaces
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")

	// Convert consecutive spaces to a single space
	for strings.Contains(s, "  ") {
		s = strings.ReplaceAll(s, "  ", " ")
	}

	return s
}

// removeEmptyEntries removes empty entries from slice if allow_empty is false
func removeEmptyEntries(entries []string, allowEmpty bool) []string {
	if allowEmpty {
		return entries
	}

	result := make([]string, 0, len(entries))
	for _, entry := range entries {
		if trimmed := strings.TrimSpace(entry); trimmed != "" {
			result = append(result, entry)
		}
	}
	return result
}

// validateInputs validates the input keys and values based on configuration
func validateInputs(cfg *config.Config, keys, values []string) error {
	seenKeys := make(map[string]bool)

	for i, key := range keys {
		if cfg.TrimWhitespace {
			key = strings.TrimSpace(key)
			keys[i] = key
		}

		// Case-insensitive comparison if case_sensitive is false
		lookupKey := key
		if !cfg.CaseSensitive {
			lookupKey = strings.ToLower(key)
		}

		// Check for empty values
		if cfg.FailOnEmpty && !cfg.AllowEmpty && (key == "" || values[i] == "") {
			return fmt.Errorf(errEmptyValue, key)
		}

		// Check for duplicate keys
		if cfg.ErrorOnDuplicate {
			if seenKeys[lookupKey] {
				return fmt.Errorf(errDuplicateKey, key)
			}
			seenKeys[lookupKey] = true
		}
	}

	return nil
}

// writeToFile writes the key-value pairs to the specified file with retry logic
func writeToFile(cfg *config.Config, filePath string, keys, values []string, varType string) (int, error) {
	maxRetries := 3
	retryDelay := time.Second

	for retry := 0; retry < maxRetries; retry++ {
		count, err := performWrite(cfg, filePath, keys, values, varType)
		if err != nil {
			if retry < maxRetries-1 {
				printer.PrintError(fmt.Sprintf("Retry %d/%d: Failed to write to file: %v", retry+1, maxRetries, err))
				time.Sleep(retryDelay)
				continue
			}
			return count, fmt.Errorf("%s: %w", fmt.Sprintf(errWriteFile, filePath), err)
		}
		return count, nil
	}
	return 0, fmt.Errorf(errMaxRetries, maxRetries)
}

// performWrite performs the actual file writing operation
func performWrite(cfg *config.Config, filePath string, keys, values []string, varType string) (int, error) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return 0, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Create value transformer
	valueTransformer := transformer.New(
		cfg.MaskSecrets,
		cfg.MaskPattern,
		cfg.ToUpper,
		cfg.ToLower,
		cfg.EncodeURL,
		cfg.EscapeNewlines,
		cfg.MaxLength,
	)

	if cfg.DebugMode {
		fmt.Printf("✍️  Writing Values:\n")
	}

	count := 0
	for i, key := range keys {
		if key == "" && !cfg.AllowEmpty {
			continue
		}

		if cfg.TrimWhitespace {
			key = strings.TrimSpace(key)
			values[i] = strings.TrimSpace(values[i])
		}

		// Transform value
		transformedValue := valueTransformer.TransformValue(values[i])

		// Write in GitHub Actions format
		if err := writeGitHubActionsFormat(file, key, transformedValue); err != nil {
			return count, err
		}

		// Print success message (with masking applied)
		maskedValue := valueTransformer.MaskValue(transformedValue)
		printer.PrintSuccess(varType, key, maskedValue)
		count++
	}

	if cfg.DebugMode {
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")
	}

	return count, nil
}

// writeGitHubActionsFormat writes a key-value pair in GitHub Actions format
func writeGitHubActionsFormat(file *os.File, key, value string) error {
	// GitHub Actions New Line Format
	line := fmt.Sprintf("%s<<%s\n%s\n%s\n", key, "EOF", value, "EOF")
	if _, err := file.WriteString(line); err != nil {
		return fmt.Errorf("failed to write line: %w", err)
	}
	return nil
}
