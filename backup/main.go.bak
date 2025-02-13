package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// Constants for error messages and defaults
const (
	errMismatchedPairs = "env_key and env_value must have the same number of entries"
	errWriteFile       = "failed to write to %s file: %v"
	localExecMsg       = "Local Execution - %s is not set, skipping writing to GitHub Actions %s."
	successMsg         = "Successfully set %s: %s=%s"
)

// Config holds the application configuration
type Config struct {
	envKeys      string
	envValues    string
	outputKeys   string
	outputValues string
	githubEnv    string
	githubOutput string
}

// loadConfig loads configuration from environment variables
func loadConfig() *Config {
	return &Config{
		envKeys:      os.Getenv("INPUT_ENV_KEY"),
		envValues:    os.Getenv("INPUT_ENV_VALUE"),
		outputKeys:   os.Getenv("INPUT_OUTPUT_KEY"),
		outputValues: os.Getenv("INPUT_OUTPUT_VALUE"),
		githubEnv:    os.Getenv("GITHUB_ENV"),
		githubOutput: os.Getenv("GITHUB_OUTPUT"),
	}
}

func setEnv(keys, values string) error {
	keyList := strings.Split(strings.TrimSpace(keys), ",")
	valueList := strings.Split(strings.TrimSpace(values), ",")

	if len(keyList) != len(valueList) {
		return fmt.Errorf(errMismatchedPairs)
	}

	envPath := os.Getenv("GITHUB_ENV")
	if envPath == "" {
		fmt.Printf(localExecMsg, "GITHUB_ENV", "environment")
		for i, key := range keyList {
			fmt.Printf(successMsg, "environment variable", strings.TrimSpace(key), strings.TrimSpace(valueList[i]))
		}
		return nil
	}

	// GitHub Actions - write to GITHUB_ENV with retry mechanism
	return writeToFile(envPath, keyList, valueList, "environment variable")
}

func setOutput(keys, values string) error {
	keyList := strings.Split(strings.TrimSpace(keys), ",")
	valueList := strings.Split(strings.TrimSpace(values), ",")

	if len(keyList) != len(valueList) {
		return fmt.Errorf(errMismatchedPairs)
	}

	outputPath := os.Getenv("GITHUB_OUTPUT")
	if outputPath == "" {
		fmt.Printf(localExecMsg, "GITHUB_OUTPUT", "output")
		for i, key := range keyList {
			fmt.Printf(successMsg, "output", strings.TrimSpace(key), strings.TrimSpace(valueList[i]))
		}
		return nil
	}

	// GitHub Actions - write to GITHUB_OUTPUT with retry mechanism
	return writeToFile(outputPath, keyList, valueList, "output")
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

func printSection(title string) {
	fmt.Printf("\n%s\n%s\n", strings.Repeat("=", 50), title)
}

func printSuccess(varType, key, value string) {
	fmt.Printf("  • %s: %s = %s\n", varType, key, value)
}

func doWrite(filePath string, keys, values []string, varType string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	printSection(fmt.Sprintf("Setting %s Variables", strings.Title(varType)))

	for i, key := range keys {
		line := fmt.Sprintf("%s=%s\n", strings.TrimSpace(key), strings.TrimSpace(values[i]))
		if _, err := file.WriteString(line); err != nil {
			return err
		}
		printSuccess(varType, key, values[i])
	}
	fmt.Println()
	return nil
}

func main() {
	config := loadConfig()

	printSection("🚀 GitHub Environment and Output Setter")

	// Set environment variables
	if err := setEnv(config.envKeys, config.envValues); err != nil {
		fmt.Printf("❌ Error setting environment variables: %v\n", err)
		os.Exit(1)
	}

	// Set output variables
	if err := setOutput(config.outputKeys, config.outputValues); err != nil {
		fmt.Printf("❌ Error setting output variables: %v\n", err)
		os.Exit(1)
	}

	// Print final status
	printSection("✅ Execution Complete")
	if config.githubEnv == "" && config.githubOutput == "" {
		fmt.Println("Mode: Local Execution (Simulation)")
	} else {
		fmt.Println("Mode: GitHub Actions")
	}
	fmt.Println(strings.Repeat("=", 50))
}
