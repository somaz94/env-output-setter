package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/somaz94/env-output-setter/internal/config"
	"github.com/somaz94/env-output-setter/internal/printer"
	"github.com/somaz94/env-output-setter/internal/writer"
)

func main() {
	cfg := config.Load()
	printer.PrintSection("🚀 GitHub Environment and Output Setter")

	// Initialize counters
	var envCount, outputCount int
	var status = "success"
	var errorMsg string

	// Set environment variables
	envCount, err := writer.SetEnv(cfg)
	if err != nil {
		errorMsg = fmt.Sprintf("❌ Error setting environment variables: %v", err)
		printer.PrintError(errorMsg)
		setOutputAndExit(0, 0, "failure", errorMsg)
	}

	// Set output variables
	outputCount, err = writer.SetOutput(cfg)
	if err != nil {
		errorMsg = fmt.Sprintf("❌ Error setting output variables: %v", err)
		printer.PrintError(errorMsg)
		setOutputAndExit(envCount, outputCount, "failure", errorMsg)
	}

	// Print final status
	printer.PrintSection("✅ Execution Complete")
	if cfg.GithubEnv == "" && cfg.GithubOutput == "" {
		printer.PrintInfo("Mode: Local Execution (Simulation)")
	} else {
		printer.PrintInfo("Mode: GitHub Actions")
	}

	setOutputAndExit(envCount, outputCount, status, errorMsg)
}

func setOutputAndExit(envCount, outputCount int, status, errorMsg string) {
	if outputFile := os.Getenv("GITHUB_OUTPUT"); outputFile != "" {
		outputs := map[string]string{
			"set_env_count":    strconv.Itoa(envCount),
			"set_output_count": strconv.Itoa(outputCount),
			"status":           status,
			"error_message":    errorMsg,
		}

		for key, value := range outputs {
			if err := appendToFile(outputFile, fmt.Sprintf("%s=%s", key, value)); err != nil {
				fmt.Printf("Error writing to GITHUB_OUTPUT: %v\n", err)
				os.Exit(1)
			}
		}
	}

	if status == "failure" {
		os.Exit(1)
	}
	os.Exit(0)
}

func appendToFile(filename, content string) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	if _, err := fmt.Fprintln(f, content); err != nil {
		return fmt.Errorf("failed to write content: %w", err)
	}
	return nil
}
