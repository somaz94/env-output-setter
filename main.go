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
		fmt.Println("GITHUB_ENV is not set, skipping writing to GitHub Actions environment.")
		// Local execution - print values
		for i, key := range keyList {
			fmt.Printf("Setting environment variable locally: %s=%s\n", key, valueList[i])
		}
	} else {
		// GitHub Actions - write to GITHUB_ENV
		file, err := os.OpenFile(envPath, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open GITHUB_ENV file: %v", err)
		}
		defer file.Close()

		for i, key := range keyList {
			line := fmt.Sprintf("%s=%s\n", key, valueList[i])
			if _, err := file.WriteString(line); err != nil {
				return fmt.Errorf("failed to write to GITHUB_ENV file: %v", err)
			}
		}
	}
	return nil
}

func setOutput(keys, values string) (string, error) {
	keyList := strings.Split(keys, ",")
	valueList := strings.Split(values, ",")
	if len(keyList) != len(valueList) {
		return "", fmt.Errorf("output_key and output_value must have the same number of entries")
	}

	outputPath := os.Getenv("GITHUB_OUTPUT")
	var outputSummary strings.Builder

	if outputPath == "" {
		fmt.Println("GITHUB_OUTPUT is not set, skipping writing to GitHub Actions output.")
		// Local execution - print values
		for i, key := range keyList {
			entry := fmt.Sprintf("%s=%s", key, valueList[i])
			fmt.Printf("Setting output variable locally: %s\n", entry)
			outputSummary.WriteString(entry + " ")
		}
	} else {
		// GitHub Actions - write to GITHUB_OUTPUT
		file, err := os.OpenFile(outputPath, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return "", fmt.Errorf("failed to open GITHUB_OUTPUT file: %v", err)
		}
		defer file.Close()

		for i, key := range keyList {
			line := fmt.Sprintf("%s=%s\n", key, valueList[i])
			if _, err := file.WriteString(line); err != nil {
				return "", fmt.Errorf("failed to write to GITHUB_OUTPUT file: %v", err)
			}
			outputSummary.WriteString(line)
		}
	}

	return outputSummary.String(), nil
}

func main() {
	envKeys := os.Getenv("INPUT_ENV_KEY")
	envValues := os.Getenv("INPUT_ENV_VALUE")
	outputKeys := os.Getenv("INPUT_OUTPUT_KEY")
	outputValues := os.Getenv("INPUT_OUTPUT_VALUE")

	if err := setEnv(envKeys, envValues); err != nil {
		fmt.Printf("Error setting environment variables: %v\n", err)
		os.Exit(1)
	}

	outputSummary, err := setOutput(outputKeys, outputValues)
	if err != nil {
		fmt.Printf("Error setting output variables: %v\n", err)
		os.Exit(1)
	}

	if os.Getenv("GITHUB_ENV") == "" && os.Getenv("GITHUB_OUTPUT") == "" {
		// Local execution - print the success message
		fmt.Printf("Local Execution - Environment and outputs set successfully: %s\n", outputSummary)
	} else {
		// GitHub Actions - write success message
		fmt.Printf("success_message=Environment and outputs set successfully: %s\n", outputSummary)
	}
}
