package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/somaz94/env-output-setter/internal/config"
	"github.com/somaz94/env-output-setter/internal/printer"
	"github.com/somaz94/env-output-setter/internal/writer"
)

// Output keys written to $GITHUB_OUTPUT. Must match action.yml outputs.
const (
	outputSetEnvCount    = "set_env_count"
	outputSetOutputCount = "set_output_count"
	outputActionStatus   = "action_status"
	outputErrorMessage   = "error_message"
	statusSuccess        = "success"
	statusFailure        = "failure"
)

func main() {
	os.Exit(run())
}

// run executes the main logic and returns the exit code.
func run() int {
	cfg := config.Load()
	printer.PrintSection("GitHub Environment and Output Setter")

	logAdvancedFeatures(cfg)

	// Set environment variables
	envCount, err := writer.SetEnv(cfg)
	if err != nil {
		errorMsg := fmt.Sprintf("Error setting environment variables: %v", err)
		printer.PrintError(errorMsg)
		writeOutputs(0, 0, statusFailure, errorMsg)
		return 1
	}

	// Set output variables
	outputCount, err := writer.SetOutput(cfg)
	if err != nil {
		errorMsg := fmt.Sprintf("Error setting output variables: %v", err)
		printer.PrintError(errorMsg)
		writeOutputs(envCount, outputCount, statusFailure, errorMsg)
		return 1
	}

	// Print final status
	printer.PrintSection("Execution Complete")
	if cfg.GithubEnv == "" && cfg.GithubOutput == "" {
		printer.PrintInfo("Mode: Local Execution (Simulation)")
	} else {
		printer.PrintInfo("Mode: GitHub Actions")
	}

	writeOutputs(envCount, outputCount, statusSuccess, "")
	return 0
}

// logAdvancedFeatures prints which optional features are enabled when debug mode is on.
func logAdvancedFeatures(cfg *config.Config) {
	if !cfg.DebugMode {
		return
	}
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

// writeOutputs writes action result outputs to the GITHUB_OUTPUT file.
// Failures are reported to stderr (not stdout) so they do not pollute the
// action's regular log stream.
func writeOutputs(envCount, outputCount int, status, errorMsg string) {
	outputFile := os.Getenv(config.GithubOutputVar)
	if outputFile == "" {
		return
	}

	outputs := map[string]string{
		outputSetEnvCount:    strconv.Itoa(envCount),
		outputSetOutputCount: strconv.Itoa(outputCount),
		outputActionStatus:   status,
		outputErrorMessage:   errorMsg,
	}

	// These status keys are written in the plain key=value form on purpose;
	// the writer package emits user-provided values via the multiline EOF form.
	for key, value := range outputs {
		if err := appendToFile(outputFile, fmt.Sprintf("%s=%s", key, value)); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to GITHUB_OUTPUT: %v\n", err)
		}
	}
}

func appendToFile(filename, content string) (err error) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("failed to close file: %w", cerr)
		}
	}()

	if _, werr := fmt.Fprintln(f, content); werr != nil {
		return fmt.Errorf("failed to write content: %w", werr)
	}
	return nil
}
