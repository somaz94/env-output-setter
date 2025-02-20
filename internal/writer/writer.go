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

func SetEnv(cfg *config.Config) (int, error) {
	return setVariables(cfg, "GITHUB_ENV", "env")
}

func SetOutput(cfg *config.Config) (int, error) {
	return setVariables(cfg, "GITHUB_OUTPUT", "output")
}

func setVariables(cfg *config.Config, envVar, varType string) (int, error) {
	var keys, values string
	switch envVar {
	case "GITHUB_ENV":
		keys, values = cfg.EnvKeys, cfg.EnvValues
	case "GITHUB_OUTPUT":
		keys, values = cfg.OutputKeys, cfg.OutputValues
	}

	// 디버그 로깅
	if cfg.DebugMode {
		printer.PrintDebugSection(strings.Title(varType))
		printer.PrintDebugInfo("📥 Input Values:\n")
		printer.PrintDebugInfo("  • Keys:      %q\n", keys)
		printer.PrintDebugInfo("  • Values:    %q\n", values)
		printer.PrintDebugInfo("  • Delimiter: %q\n\n", cfg.Delimiter)
	}

	// 구분자로 먼저 분리
	keyList := strings.Split(keys, cfg.Delimiter)
	valueList := strings.Split(values, cfg.Delimiter)

	// 각 항목별로 whitespace 정규화
	for i := range keyList {
		keyList[i] = normalizeWhitespace(keyList[i])
	}
	for i := range valueList {
		valueList[i] = normalizeWhitespace(valueList[i])
	}

	// 각 항목의 앞뒤 공백 제거
	for i := range keyList {
		keyList[i] = strings.TrimSpace(keyList[i])
	}
	for i := range valueList {
		valueList[i] = strings.TrimSpace(valueList[i])
	}

	// 빈 항목 제거 (allow_empty가 true일 때는 제거하지 않음)
	keyList = removeEmptyEntries(keyList, cfg.AllowEmpty)
	valueList = removeEmptyEntries(valueList, cfg.AllowEmpty)

	// 디버그 로깅
	if cfg.DebugMode {
		printer.PrintDebugInfo("📋 Processed Values:\n")
		printer.PrintDebugInfo("  • Keys:   %v\n", keyList)
		printer.PrintDebugInfo("  • Values: %v\n\n", valueList)
	}

	if len(keyList) != len(valueList) {
		return 0, fmt.Errorf("%s (keys: %d, values: %d)", errMismatchedPairs, len(keyList), len(valueList))
	}

	if err := validateInputs(cfg, keyList, valueList); err != nil {
		return 0, err
	}

	filePath := os.Getenv(envVar)
	if filePath == "" {
		fmt.Printf(localExecMsg, envVar, varType)
		for i, key := range keyList {
			printer.PrintSuccess(varType, key, valueList[i])
		}
		return len(keyList), nil
	}

	count, err := writeToFile(cfg, filePath, keyList, valueList, varType)
	if err != nil {
		return count, err
	}

	return count, nil
}

// normalizeWhitespace normalizes all whitespace including newlines
func normalizeWhitespace(s string) string {
	// 실제 줄바꿈을 공백으로 변환
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")

	// 연속된 공백을 하나의 공백으로 변환
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

func validateInputs(cfg *config.Config, keys, values []string) error {
	seenKeys := make(map[string]bool)

	for i, key := range keys {
		if cfg.TrimWhitespace {
			key = strings.TrimSpace(key)
			keys[i] = key
		}

		if !cfg.CaseSensitive {
			key = strings.ToLower(key)
		}

		if cfg.FailOnEmpty && !cfg.AllowEmpty && (key == "" || values[i] == "") {
			return fmt.Errorf(errEmptyValue, key)
		}

		if cfg.ErrorOnDuplicate {
			if seenKeys[key] {
				return fmt.Errorf(errDuplicateKey, key)
			}
			seenKeys[key] = true
		}
	}

	return nil
}

func writeToFile(cfg *config.Config, filePath string, keys, values []string, varType string) (int, error) {
	maxRetries := 3
	retryDelay := time.Second

	for retry := 0; retry < maxRetries; retry++ {
		count, err := doWrite(cfg, filePath, keys, values, varType)
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

func doWrite(cfg *config.Config, filePath string, keys, values []string, varType string) (int, error) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return 0, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	trans := transformer.New(
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

		transformedValue := trans.TransformValue(values[i])

		// GitHub Actions New Line Format
		line := fmt.Sprintf("%s<<%s\n%s\n%s\n", key, "EOF", transformedValue, "EOF")
		if _, err := file.WriteString(line); err != nil {
			return count, fmt.Errorf("failed to write line: %w", err)
		}

		maskedValue := trans.MaskValue(transformedValue)
		printer.PrintSuccess(varType, key, maskedValue)
		count++
	}

	if cfg.DebugMode {
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")
	}

	return count, nil
}
