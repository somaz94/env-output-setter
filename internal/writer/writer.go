package writer

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/somaz94/env-output-setter/internal/config"
	"github.com/somaz94/env-output-setter/internal/printer"
	"github.com/somaz94/env-output-setter/internal/transformer"
)

// Error messages
const (
	errMismatchedPairs = "env_key and env_value must have the same number of entries"
	errWriteFile       = "failed to write to %s file"
	errEmptyValue      = "empty value found for key: %s"
	errDuplicateKey    = "duplicate key found: %s"
	errMaxRetries      = "failed to write after %d retries"
	localExecMsg       = "Local Execution - %s is not set, skipping writing to GitHub Actions %s"
)

// File types
const (
	envFileType    = "env"
	outputFileType = "output"
)

// GitHub environment variables
const (
	githubEnvVar    = "GITHUB_ENV"
	githubOutputVar = "GITHUB_OUTPUT"
)

// SetEnv sets environment variables in GitHub Actions environment file.
// It processes the env_key and env_value inputs and writes them to the GITHUB_ENV file.
func SetEnv(cfg *config.Config) (int, error) {
	return setVariables(cfg, githubEnvVar, envFileType)
}

// SetOutput sets output variables in GitHub Actions output file.
// It processes the output_key and output_value inputs and writes them to the GITHUB_OUTPUT file.
// If export_as_env is enabled, it also exports the output variables as environment variables.
func SetOutput(cfg *config.Config) (int, error) {
	// Set output variables
	count, err := setVariables(cfg, githubOutputVar, outputFileType)
	if err != nil {
		return count, err
	}

	// Export output variables as environment variables if enabled
	if cfg.ExportAsEnv {
		return exportOutputAsEnv(cfg, count)
	}

	return count, nil
}

// exportOutputAsEnv exports output variables as environment variables.
// It reads the output variables and writes them to the environment file.
func exportOutputAsEnv(cfg *config.Config, outputCount int) (int, error) {
	keys, values := getInputValues(cfg, githubOutputVar)
	keyList, valueList, err := processInputValues(cfg, keys, values)
	if err != nil {
		return outputCount, err
	}

	envFilePath := os.Getenv(githubEnvVar)
	// If we're not in GitHub Actions, just log the values
	if envFilePath == "" {
		return outputCount, nil
	}

	envCount, err := writeToFile(cfg, envFilePath, keyList, valueList, "env (from output)")
	if err != nil {
		return outputCount, err
	}

	return outputCount + envCount, nil
}

// setVariables handles setting variables for both env and output files.
// It's the core function that processes inputs and writes them to the appropriate file.
func setVariables(cfg *config.Config, envVar, varType string) (int, error) {
	// Get input values based on the variable type
	keys, values := getInputValues(cfg, envVar)

	// Log input values if debug mode is enabled
	logInputValues(cfg, varType, keys, values)

	// Process and validate input values
	keyList, valueList, err := processInputValues(cfg, keys, values)
	if err != nil {
		return 0, err
	}

	// Log processed values if debug mode is enabled
	logProcessedValues(cfg, keyList, valueList)

	// Validate input constraints (empty values, duplicates, etc.)
	if err := validateInputs(cfg, keyList, valueList); err != nil {
		return 0, err
	}

	// Get file path from environment variable
	filePath := os.Getenv(envVar)

	// Handle local execution (not in GitHub Actions)
	if filePath == "" {
		return handleLocalExecution(envVar, varType, keyList, valueList)
	}

	// Write variables to the file
	return writeToFile(cfg, filePath, keyList, valueList, varType)
}

// getInputValues returns the appropriate keys and values based on the variable type.
// It selects between env_key/env_value and output_key/output_value based on the envVar parameter.
func getInputValues(cfg *config.Config, envVar string) (string, string) {
	switch envVar {
	case githubEnvVar:
		return cfg.EnvKeys, cfg.EnvValues
	case githubOutputVar:
		return cfg.OutputKeys, cfg.OutputValues
	default:
		return "", ""
	}
}

// logInputValues logs the original input values if debug mode is enabled.
// It displays the raw input keys, values, and delimiter.
func logInputValues(cfg *config.Config, varType, keys, values string) {
	if !cfg.DebugMode {
		return
	}

	printer.PrintDebugSection(strings.Title(varType))
	printer.PrintDebugInfo("📥 Input Values:\n")
	printer.PrintDebugInfo("  • Keys:      %q\n", keys)
	printer.PrintDebugInfo("  • Values:    %q\n", values)
	printer.PrintDebugInfo("  • Delimiter: %q\n\n", cfg.Delimiter)
}

// processInputValues processes the input strings into lists with proper formatting.
// It handles splitting, trimming, and processing JSON values if json_support is enabled.
func processInputValues(cfg *config.Config, keys, values string) ([]string, []string, error) {
	// Split input strings by delimiter
	keyList := strings.Split(keys, cfg.Delimiter)
	valueList := strings.Split(values, cfg.Delimiter)

	// Process keys and values for whitespace
	keyList = processWhitespace(keyList)
	valueList = processWhitespace(valueList)

	// Filter out empty entries if not allowed
	keyList = removeEmptyEntries(keyList, cfg.AllowEmpty)
	valueList = removeEmptyEntries(valueList, cfg.AllowEmpty)

	// Process JSON values if enabled
	if cfg.JsonSupport {
		keyList, valueList = processJsonValues(keyList, valueList)
	}

	// Ensure key-value pairs match
	if len(keyList) != len(valueList) {
		return nil, nil, fmt.Errorf("%s (keys: %d, values: %d)",
			errMismatchedPairs, len(keyList), len(valueList))
	}

	return keyList, valueList, nil
}

// processWhitespace normalizes and trims whitespace from all entries in a list.
func processWhitespace(entries []string) []string {
	result := make([]string, len(entries))
	for i, entry := range entries {
		// Normalize whitespace (convert newlines to spaces, reduce multiple spaces)
		normalized := normalizeWhitespace(entry)
		// Trim any remaining leading/trailing whitespace
		result[i] = strings.TrimSpace(normalized)
	}
	return result
}

// normalizeWhitespace converts all whitespace sequences to a single space.
// It handles newlines, carriage returns, and multiple consecutive spaces.
func normalizeWhitespace(s string) string {
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
// If allowEmpty is true, all entries are preserved.
func removeEmptyEntries(entries []string, allowEmpty bool) []string {
	if allowEmpty {
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

// processJsonValues extracts nested properties from JSON values.
// It processes JSON objects and arrays, creating flattened key-value pairs.
func processJsonValues(keyList, valueList []string) ([]string, []string) {
	// Make a copy of the original lists
	originalKeyCount := len(keyList)
	resultKeys := make([]string, len(keyList))
	resultValues := make([]string, len(valueList))
	copy(resultKeys, keyList)
	copy(resultValues, valueList)

	// Process each JSON value in the original list
	for i := 0; i < originalKeyCount; i++ {
		value := valueList[i]
		key := keyList[i]

		// Check if value looks like JSON
		if !isJsonLike(value) {
			continue
		}

		// Try to parse the JSON value
		var jsonData interface{}
		if err := json.Unmarshal([]byte(value), &jsonData); err != nil {
			printer.PrintWarning(fmt.Sprintf("Warning: Invalid JSON for key '%s': %v", key, err))
			continue
		}

		// Extract nested values based on the JSON type
		switch typedData := jsonData.(type) {
		case map[string]interface{}:
			// Handle JSON object
			nestedKeys, nestedValues := extractNestedJSON(key, typedData, "")
			resultKeys = append(resultKeys, nestedKeys...)
			resultValues = append(resultValues, nestedValues...)
		case []interface{}:
			// Handle JSON array
			for idx, item := range typedData {
				arrayKey := fmt.Sprintf("%s_%d", key, idx)
				resultKeys = append(resultKeys, arrayKey)
				resultValues = append(resultValues, fmt.Sprintf("%v", item))

				// Process nested objects in arrays
				if mapItem, ok := item.(map[string]interface{}); ok {
					objKeys, objValues := extractNestedJSON(arrayKey, mapItem, "")
					resultKeys = append(resultKeys, objKeys...)
					resultValues = append(resultValues, objValues...)
				}
			}
		}
	}

	return resultKeys, resultValues
}

// isJsonLike checks if a string looks like JSON (object or array).
func isJsonLike(value string) bool {
	value = strings.TrimSpace(value)
	return (strings.HasPrefix(value, "{") && strings.HasSuffix(value, "}")) ||
		(strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]"))
}

// extractNestedJSON flattens a nested JSON object into key-value pairs.
// It recursively processes nested objects and arrays, creating concatenated keys.
func extractNestedJSON(prefix string, jsonObj map[string]interface{}, groupPrefix string) ([]string, []string) {
	var keys []string
	var values []string

	// Prepare the key prefix with group prefix if provided
	keyPrefix := prefix
	if groupPrefix != "" && !strings.HasPrefix(prefix, groupPrefix) {
		keyPrefix = fmt.Sprintf("%s_%s", groupPrefix, prefix)
	}

	// Process each property in the JSON object
	for propKey, propValue := range jsonObj {
		nestedKey := fmt.Sprintf("%s_%s", keyPrefix, propKey)

		// Handle different value types
		switch typedValue := propValue.(type) {
		case map[string]interface{}:
			// Recursively process nested objects
			nestedKeys, nestedValues := extractNestedJSON(nestedKey, typedValue, "")
			keys = append(keys, nestedKeys...)
			values = append(values, nestedValues...)
		case []interface{}:
			// Process arrays
			for i, item := range typedValue {
				arrayKey := fmt.Sprintf("%s_%d", nestedKey, i)
				keys = append(keys, arrayKey)
				values = append(values, fmt.Sprintf("%v", item))

				// Process objects within arrays
				if mapItem, ok := item.(map[string]interface{}); ok {
					subKeys, subValues := extractNestedJSON(arrayKey, mapItem, "")
					keys = append(keys, subKeys...)
					values = append(values, subValues...)
				}
			}
		default:
			// Handle primitive values
			keys = append(keys, nestedKey)
			values = append(values, fmt.Sprintf("%v", typedValue))
		}
	}

	return keys, values
}

// logProcessedValues logs the processed key-value pairs if debug mode is enabled.
func logProcessedValues(cfg *config.Config, keyList, valueList []string) {
	if !cfg.DebugMode {
		return
	}

	printer.PrintDebugInfo("📋 Processed Values:\n")
	printer.PrintDebugInfo("  • Keys:   %v\n", keyList)
	printer.PrintDebugInfo("  • Values: %v\n\n", valueList)
}

// validateInputs checks for empty values and duplicate keys based on configuration.
func validateInputs(cfg *config.Config, keys, values []string) error {
	seenKeys := make(map[string]bool)

	for i, key := range keys {
		// Apply trimming if configured
		if cfg.TrimWhitespace {
			key = strings.TrimSpace(key)
			keys[i] = key
		}

		// Prepare key for duplicate checking
		lookupKey := key
		if !cfg.CaseSensitive {
			lookupKey = strings.ToLower(key)
		}

		// Check for empty values if configured to fail
		if cfg.FailOnEmpty && !cfg.AllowEmpty && (key == "" || values[i] == "") {
			return fmt.Errorf(errEmptyValue, key)
		}

		// Check for duplicate keys if configured
		if cfg.ErrorOnDuplicate {
			if seenKeys[lookupKey] {
				return fmt.Errorf(errDuplicateKey, key)
			}
			seenKeys[lookupKey] = true
		}
	}

	return nil
}

// handleLocalExecution handles variable setting when not running in GitHub Actions.
// It prints values to the console instead of writing to a file.
func handleLocalExecution(envVar, varType string, keyList, valueList []string) (int, error) {
	fmt.Printf(localExecMsg, envVar, varType)
	for i, key := range keyList {
		printer.PrintSuccess(varType, key, valueList[i])
	}
	return len(keyList), nil
}

// writeToFile writes key-value pairs to a file with retry logic.
// It attempts to write up to maxRetries times with delays between attempts.
func writeToFile(cfg *config.Config, filePath string, keys, values []string, varType string) (int, error) {
	maxRetries := 3
	retryDelay := time.Second
	var lastError error

	// Attempt writing with retries
	for retry := 0; retry < maxRetries; retry++ {
		count, err := performWrite(cfg, filePath, keys, values, varType)
		if err == nil {
			// Success - write status
			_, _ = performWrite(cfg, filePath, []string{"status"}, []string{"success"}, varType)
			return count, nil
		}

		// Handle error
		lastError = err
		if retry < maxRetries-1 {
			printer.PrintError(fmt.Sprintf("Retry %d/%d: Failed to write to file: %v",
				retry+1, maxRetries, err))
			time.Sleep(retryDelay)
		}
	}

	// Write failure status after exhausting retries
	_, _ = performWrite(cfg, filePath,
		[]string{"status", "error_message"},
		[]string{"failure", lastError.Error()},
		varType)

	return 0, fmt.Errorf(errMaxRetries, maxRetries)
}

// performWrite writes key-value pairs to a file in GitHub Actions format.
// It handles file opening, value transformation, and formatting.
func performWrite(cfg *config.Config, filePath string, keys, values []string, varType string) (int, error) {
	// Open the file
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return 0, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Create transformer for values
	valueTransformer := transformer.New(
		cfg.MaskSecrets,
		cfg.MaskPattern,
		cfg.ToUpper,
		cfg.ToLower,
		cfg.EncodeURL,
		cfg.EscapeNewlines,
		cfg.MaxLength,
	)

	// Write header in debug mode
	if cfg.DebugMode {
		fmt.Printf("✍️  Writing Values:\n")
	}

	// Write each key-value pair
	count := 0
	for i, key := range keys {
		// Skip empty keys unless allowed
		if key == "" && !cfg.AllowEmpty {
			continue
		}

		// Apply whitespace trimming if configured
		if cfg.TrimWhitespace {
			key = strings.TrimSpace(key)
			values[i] = strings.TrimSpace(values[i])
		}

		// Transform and write the value
		transformedValue := valueTransformer.TransformValue(values[i], cfg.JsonSupport)
		if err := writeGitHubActionsFormat(file, key, transformedValue); err != nil {
			return count, err
		}

		// Print success message with masking
		maskedValue := valueTransformer.MaskValue(transformedValue)
		printer.PrintSuccess(varType, key, maskedValue)
		count++
	}

	// Write footer in debug mode
	if cfg.DebugMode {
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")
	}

	return count, nil
}

// writeGitHubActionsFormat writes a key-value pair in GitHub Actions format.
// Format: key<<EOF\nvalue\nEOF
func writeGitHubActionsFormat(file *os.File, key, value string) error {
	line := fmt.Sprintf("%s<<%s\n%s\n%s\n", key, "EOF", value, "EOF")
	if _, err := file.WriteString(line); err != nil {
		return fmt.Errorf("failed to write line: %w", err)
	}
	return nil
}
