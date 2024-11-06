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
	file, err := os.OpenFile(envPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	for i, key := range keyList {
		_, err := file.WriteString(fmt.Sprintf("%s=%s\n", key, valueList[i]))
		if err != nil {
			return err
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
	file, err := os.OpenFile(outputPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var outputSummary strings.Builder
	for i, key := range keyList {
		entry := fmt.Sprintf("%s=%s", key, valueList[i])
		outputSummary.WriteString(entry + " ")
		// Write each output key-value pair individually
		_, err := file.WriteString(fmt.Sprintf("%s=%s\n", key, valueList[i]))
		if err != nil {
			return "", err
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

	// Print success message for GitHub Action output
	fmt.Printf("success_message=Environment and outputs set successfully: %s\n", outputSummary)
}
