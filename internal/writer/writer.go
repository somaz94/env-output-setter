package writer

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/somaz94/env-output-setter/internal/config"
	"github.com/somaz94/env-output-setter/internal/printer"
)

const (
	errMismatchedPairs = "env_key and env_value must have the same number of entries"
	errWriteFile       = "failed to write to %s file: %v"
	errEmptyValue      = "empty value found for key: %s"
	errDuplicateKey    = "duplicate key found: %s"
	localExecMsg       = "Local Execution - %s is not set, skipping writing to GitHub Actions %s."
)

func SetEnv(cfg *config.Config) (int, error) {
	return setVariables(cfg, "GITHUB_ENV", "environment variable")
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

	keyList := strings.Split(strings.TrimSpace(keys), cfg.Delimiter)
	valueList := strings.Split(strings.TrimSpace(values), cfg.Delimiter)

	if len(keyList) != len(valueList) {
		return 0, fmt.Errorf(errMismatchedPairs)
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

		if cfg.FailOnEmpty && (key == "" || values[i] == "") {
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
			return count, fmt.Errorf(errWriteFile, filePath, err)
		}
		return count, nil
	}
	return 0, fmt.Errorf("failed to write after %d retries", maxRetries)
}

func doWrite(cfg *config.Config, filePath string, keys, values []string, varType string) (int, error) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	printer.PrintSection(fmt.Sprintf("Setting %s Variables", strings.Title(varType)))

	count := 0
	for i, key := range keys {
		if cfg.TrimWhitespace {
			key = strings.TrimSpace(key)
			values[i] = strings.TrimSpace(values[i])
		}

		line := fmt.Sprintf("%s=%s\n", key, values[i])
		if _, err := file.WriteString(line); err != nil {
			return count, err
		}
		printer.PrintSuccess(varType, key, values[i])
		count++
	}
	fmt.Println()
	return count, nil
}
