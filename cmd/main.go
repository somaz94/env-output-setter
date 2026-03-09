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
	os.Exit(run())
}

// run executes the main logic and returns the exit code.
func run() int {
	cfg := config.Load()
	printer.PrintSection("GitHub Environment and Output Setter")

	// Log advanced feature status
	if cfg.DebugMode {
		printer.PrintInfo("Advanced Features Status:")
		if cfg.GroupPrefix != "" {
			printer.PrintInfo(fmt.Sprintf("  * Group Prefix: %s", cfg.GroupPrefix))
		}
		if cfg.JsonSupport {
			printer.PrintInfo("  * JSON Support: Enabled")
		}
		if cfg.ExportAsEnv {
			printer.PrintInfo("  * Export Output as Env: Enabled")
		}
	}

	// Initialize counters
	var envCount, outputCount int
	var status = "success"
	var errorMsg string

	// Set environment variables
	envCount, err := writer.SetEnv(cfg)
	if err != nil {
		errorMsg = fmt.Sprintf("Error setting environment variables: %v", err)
		printer.PrintError(errorMsg)
		writeOutputs(0, 0, "failure", errorMsg)
		return 1
	}

	// Set output variables
	outputCount, err = writer.SetOutput(cfg)
	if err != nil {
		errorMsg = fmt.Sprintf("Error setting output variables: %v", err)
		printer.PrintError(errorMsg)
		writeOutputs(envCount, outputCount, "failure", errorMsg)
		return 1
	}

	// Print final status
	printer.PrintSection("Execution Complete")
	if cfg.GithubEnv == "" && cfg.GithubOutput == "" {
		printer.PrintInfo("Mode: Local Execution (Simulation)")
	} else {
		printer.PrintInfo("Mode: GitHub Actions")
	}

	writeOutputs(envCount, outputCount, status, errorMsg)
	return 0
}

// writeOutputs writes action result outputs to the GITHUB_OUTPUT file.
func writeOutputs(envCount, outputCount int, status, errorMsg string) {
	outputFile := os.Getenv("GITHUB_OUTPUT")
	if outputFile == "" {
		return
	}

	outputs := map[string]string{
		"set_env_count":    strconv.Itoa(envCount),
		"set_output_count": strconv.Itoa(outputCount),
		"status":           status,
		"error_message":    errorMsg,
	}

	for key, value := range outputs {
		if err := appendToFile(outputFile, fmt.Sprintf("%s=%s", key, value)); err != nil {
			fmt.Printf("Error writing to GITHUB_OUTPUT: %v\n", err)
		}
	}
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
