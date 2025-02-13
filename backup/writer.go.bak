package writer

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/somaz94/env-output-setter/internal/printer"
)

const (
	errMismatchedPairs = "env_key and env_value must have the same number of entries"
	errWriteFile       = "failed to write to %s file: %v"
	localExecMsg       = "Local Execution - %s is not set, skipping writing to GitHub Actions %s."
)

func SetEnv(keys, values string) error {
	return setVariables(keys, values, "GITHUB_ENV", "environment variable")
}

func SetOutput(keys, values string) error {
	return setVariables(keys, values, "GITHUB_OUTPUT", "output")
}

func setVariables(keys, values, envVar, varType string) error {
	keyList := strings.Split(strings.TrimSpace(keys), ",")
	valueList := strings.Split(strings.TrimSpace(values), ",")

	if len(keyList) != len(valueList) {
		return fmt.Errorf(errMismatchedPairs)
	}

	filePath := os.Getenv(envVar)
	if filePath == "" {
		fmt.Printf(localExecMsg, envVar, varType)
		for i, key := range keyList {
			printer.PrintSuccess(varType, strings.TrimSpace(key), strings.TrimSpace(valueList[i]))
		}
		return nil
	}

	return writeToFile(filePath, keyList, valueList, varType)
}

func writeToFile(filePath string, keys, values []string, varType string) error {
	maxRetries := 3
	retryDelay := time.Second

	for retry := 0; retry < maxRetries; retry++ {
		if err := doWrite(filePath, keys, values, varType); err != nil {
			if retry < maxRetries-1 {
				fmt.Printf("Retry %d/%d: Failed to write to file: %v\n", retry+1, maxRetries, err)
				time.Sleep(retryDelay)
				continue
			}
			return fmt.Errorf(errWriteFile, filePath, err)
		}
		return nil
	}
	return fmt.Errorf("failed to write after %d retries", maxRetries)
}

func doWrite(filePath string, keys, values []string, varType string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	printer.PrintSection(fmt.Sprintf("Setting %s Variables", strings.Title(varType)))

	for i, key := range keys {
		line := fmt.Sprintf("%s=%s\n", strings.TrimSpace(key), strings.TrimSpace(values[i]))
		if _, err := file.WriteString(line); err != nil {
			return err
		}
		printer.PrintSuccess(varType, key, values[i])
	}
	fmt.Println()
	return nil
}
