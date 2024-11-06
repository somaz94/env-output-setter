package main

import (
	"fmt"
	"os"
	"strings"
)

func setEnv(keys, values string) error {
	keyList := strings.Split(keys, ",")
	valueList := strings.Split(values, ",")
	if len(keyList) != len(valueList) {
		return fmt.Errorf("env_key and env_value must have the same number of entries")
	}

	envPath := os.Getenv("GITHUB_ENV")
	if envPath == "" {
		// Local execution - print values
		fmt.Println("GITHUB_ENV is not set, skipping writing to GitHub Actions environment.")
		for i, key := range keyList {
			fmt.Printf("Setting environment variable locally: %s=%s\n", key, valueList[i])
		}
	} else {
		// GitHub Actions - write to GITHUB_ENV
		for i, key := range keyList {
			line := fmt.Sprintf("%s=%s\n", key, valueList[i])
			if err := appendToFile(envPath, line); err != nil {
				return fmt.Errorf("failed to write to GITHUB_ENV file: %v", err)
			}
		}
	}
	return nil
}

func setOutput(keys, values string) error {
	keyList := strings.Split(keys, ",")
	valueList := strings.Split(values, ",")
	if len(keyList) != len(valueList) {
		return fmt.Errorf("output_key and output_value must have the same number of entries")
	}

	outputPath := os.Getenv("GITHUB_OUTPUT")
	if outputPath == "" {
		// Local execution - print values
		fmt.Println("GITHUB_OUTPUT is not set, skipping writing to GitHub Actions output.")
		for i, key := range keyList {
			fmt.Printf("Setting output variable locally: %s=%s\n", key, valueList[i])
		}
	} else {
		// GitHub Actions - write to GITHUB_OUTPUT
		for i, key := range keyList {
			line := fmt.Sprintf("%s=%s\n", key, valueList[i])
			if err := appendToFile(outputPath, line); err != nil {
				return fmt.Errorf("failed to write to GITHUB_OUTPUT file: %v", err)
			}
		}
	}
	return nil
}

func appendToFile(filePath, text string) error {
	// Append the provided text to the file specified by filePath
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(text); err != nil {
		return err
	}
	return nil
}

func main() {
	// Read input keys and values for environment variables and output
	envKeys := os.Getenv("INPUT_ENV_KEY")
	envValues := os.Getenv("INPUT_ENV_VALUE")
	outputKeys := os.Getenv("INPUT_OUTPUT_KEY")
	outputValues := os.Getenv("INPUT_OUTPUT_VALUE")

	// Set environment variables
	if err := setEnv(envKeys, envValues); err != nil {
		fmt.Printf("Error setting environment variables: %v\n", err)
		os.Exit(1)
	}

	// Set output variables
	if err := setOutput(outputKeys, outputValues); err != nil {
		fmt.Printf("Error setting output variables: %v\n", err)
		os.Exit(1)
	}

	// Print success message for local execution
	if os.Getenv("GITHUB_ENV") == "" && os.Getenv("GITHUB_OUTPUT") == "" {
		fmt.Printf("Local Execution - Environment and outputs set successfully.\n")
	} else {
		// Success message if running within GitHub Actions
		fmt.Printf("Environment and outputs set successfully.\n")
	}
}
