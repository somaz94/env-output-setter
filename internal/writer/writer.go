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
	count, err := setVariables(cfg, "GITHUB_OUTPUT", "output")
	if err != nil {
		return count, err
	}

	// output 변수도 환경 변수로 내보내기 옵션이 활성화된 경우
	if cfg.ExportAsEnv {
		keys, values := getInputValues(cfg, "GITHUB_OUTPUT")
		keyList, valueList, err := processInputValues(cfg, keys, values)
		if err != nil {
			return count, err
		}

		envCount, err := writeToFile(cfg, os.Getenv("GITHUB_ENV"), keyList, valueList, "env (from output)")
		if err != nil {
			return count, err
		}

		return count + envCount, nil
	}

	return count, nil
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

	// JSON 지원이 활성화된 경우 JSON 객체를 처리
	if cfg.JsonSupport {
		// 원래 키와 값의 복사본을 만들어 새 항목 추가시 반복문에 영향을 주지 않도록 함
		originalKeyCount := len(keyList)

		for i := 0; i < originalKeyCount; i++ {
			value := valueList[i]
			key := keyList[i]

			// 값이 JSON 객체인지 확인
			if strings.HasPrefix(value, "{") && strings.HasSuffix(value, "}") {
				var jsonObj map[string]interface{}
				if err := json.Unmarshal([]byte(value), &jsonObj); err == nil {
					// 중첩된 JSON 속성을 평면화하고 환경 변수로 추가
					extractedKeys, extractedValues := extractNestedJSON(key, jsonObj, cfg.GroupPrefix)
					keyList = append(keyList, extractedKeys...)
					valueList = append(valueList, extractedValues...)
				}
			} else if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
				// JSON 배열 처리
				var jsonArray []interface{}
				if err := json.Unmarshal([]byte(value), &jsonArray); err == nil {
					// 배열 항목을 환경 변수로 추가
					for idx, item := range jsonArray {
						arrayKey := fmt.Sprintf("%s_%d", key, idx)
						keyList = append(keyList, arrayKey)
						valueList = append(valueList, fmt.Sprintf("%v", item))

						// 배열 항목이 객체인 경우 재귀적으로 처리
						if mapItem, ok := item.(map[string]interface{}); ok {
							nestedKeys, nestedValues := extractNestedJSON(arrayKey, mapItem, cfg.GroupPrefix)
							keyList = append(keyList, nestedKeys...)
							valueList = append(valueList, nestedValues...)
						}
					}
				}
			}
		}
	}

	// 이 시점에서 keyList와 valueList는 추가된 중첩 JSON 속성을 포함
	// Check if the number of keys and values match
	if len(keyList) != len(valueList) {
		return nil, nil, fmt.Errorf("%s (keys: %d, values: %d)", errMismatchedPairs, len(keyList), len(valueList))
	}

	return keyList, valueList, nil
}

// extractNestedJSON은 중첩된 JSON 객체의 속성을 추출하여 평면화된 키-값 쌍을 반환합니다.
func extractNestedJSON(prefix string, jsonObj map[string]interface{}, groupPrefix string) ([]string, []string) {
	var keys []string
	var values []string

	// 그룹 접두사 적용
	keyPrefix := prefix
	if groupPrefix != "" {
		if strings.HasPrefix(prefix, groupPrefix) {
			// 이미 접두사가 있으면 그대로 사용
			keyPrefix = prefix
		} else {
			// 접두사 추가
			keyPrefix = fmt.Sprintf("%s_%s", groupPrefix, prefix)
		}
	}

	// 객체의 각 속성을 평면화된 키-값 쌍으로 변환
	for k, v := range jsonObj {
		nestedKey := fmt.Sprintf("%s_%s", keyPrefix, k)

		// 값 유형에 따라 처리
		switch val := v.(type) {
		case map[string]interface{}:
			// 중첩된 객체는 재귀적으로 처리
			nestedKeys, nestedValues := extractNestedJSON(nestedKey, val, "")
			keys = append(keys, nestedKeys...)
			values = append(values, nestedValues...)
		case []interface{}:
			// 배열은 인덱스를 키에 추가하여 처리
			for i, item := range val {
				arrayKey := fmt.Sprintf("%s_%d", nestedKey, i)
				keys = append(keys, arrayKey)
				values = append(values, fmt.Sprintf("%v", item))

				// 배열 항목이 객체인 경우 재귀적으로 처리
				if mapItem, ok := item.(map[string]interface{}); ok {
					subKeys, subValues := extractNestedJSON(arrayKey, mapItem, "")
					keys = append(keys, subKeys...)
					values = append(values, subValues...)
				}
			}
		default:
			// 기본 값 유형은 직접 추가
			keys = append(keys, nestedKey)
			values = append(values, fmt.Sprintf("%v", val))
		}
	}

	return keys, values
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
		transformedValue := valueTransformer.TransformValue(values[i], cfg.JsonSupport)

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
