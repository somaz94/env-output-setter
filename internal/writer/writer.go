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

// Error messages
const (
	errWriteFile  = "failed to write to %s file"
	errMaxRetries = "failed to write after %d retries"
	localExecMsg  = "Local Execution - %s is not set, skipping writing to GitHub Actions %s"
)

// File types
const (
	envFileType    = "env"
	outputFileType = "output"
)

// GitHub environment variables
const (
	githubEnvVar    = "GITHUB_ENV"
	githubOutputVar = "GITHUB_OUTPUT"
)

// Writer handles writing environment variables and outputs to GitHub Actions files.
type Writer struct {
	cfg       *config.Config
	processor *Processor
	validator *Validator
}

// NewWriter creates a new Writer instance.
func NewWriter(cfg *config.Config) *Writer {
	return &Writer{
		cfg:       cfg,
		processor: NewProcessor(cfg),
		validator: NewValidator(cfg),
	}
}

// SetEnv sets environment variables in GitHub Actions environment file.
// It processes the env_key and env_value inputs and writes them to the GITHUB_ENV file.
func SetEnv(cfg *config.Config) (int, error) {
	w := NewWriter(cfg)
	return w.setVariables(githubEnvVar, envFileType)
}

// SetOutput sets output variables in GitHub Actions output file.
// It processes the output_key and output_value inputs and writes them to the GITHUB_OUTPUT file.
// If export_as_env is enabled, it also exports the output variables as environment variables.
func SetOutput(cfg *config.Config) (int, error) {
	w := NewWriter(cfg)

	// Set output variables
	count, err := w.setVariables(githubOutputVar, outputFileType)
	if err != nil {
		return count, err
	}

	// Export output variables as environment variables if enabled
	if cfg.ExportAsEnv {
		return w.exportOutputAsEnv(count)
	}

	return count, nil
}

// exportOutputAsEnv exports output variables as environment variables.
// It reads the output variables and writes them to the environment file.
func (w *Writer) exportOutputAsEnv(outputCount int) (int, error) {
	keys, values := w.getInputValues(githubOutputVar)
	keyList, valueList, err := w.processor.ProcessInputValues(keys, values)
	if err != nil {
		return outputCount, err
	}

	envFilePath := os.Getenv(githubEnvVar)
	// If we're not in GitHub Actions, just log the values
	if envFilePath == "" {
		return outputCount, nil
	}

	envCount, err := w.writeToFile(envFilePath, keyList, valueList, "env (from output)")
	if err != nil {
		return outputCount, err
	}

	return outputCount + envCount, nil
}

// setVariables handles setting variables for both env and output files.
// It's the core function that processes inputs and writes them to the appropriate file.
func (w *Writer) setVariables(envVar, varType string) (int, error) {
	// Get input values based on the variable type
	keys, values := w.getInputValues(envVar)

	// Log input values if debug mode is enabled
	w.processor.LogInputValues(varType, keys, values)

	// Process and validate input values
	keyList, valueList, err := w.processor.ProcessInputValues(keys, values)
	if err != nil {
		return 0, err
	}

	// Log processed values if debug mode is enabled
	w.processor.LogProcessedValues(keyList, valueList)

	// Validate pairs match
	if err := w.validator.ValidatePairs(keyList, valueList); err != nil {
		return 0, err
	}

	// Validate input constraints (empty values, duplicates, etc.)
	if err := w.validator.ValidateInputs(keyList, valueList); err != nil {
		return 0, err
	}

	// Get file path from environment variable
	filePath := os.Getenv(envVar)

	// Handle local execution (not in GitHub Actions)
	if filePath == "" {
		return w.handleLocalExecution(envVar, varType, keyList, valueList)
	}

	// Write variables to the file
	return w.writeToFile(filePath, keyList, valueList, varType)
}

// getInputValues returns the appropriate keys and values based on the variable type.
// It selects between env_key/env_value and output_key/output_value based on the envVar parameter.
func (w *Writer) getInputValues(envVar string) (string, string) {
	switch envVar {
	case githubEnvVar:
		return w.cfg.EnvKeys, w.cfg.EnvValues
	case githubOutputVar:
		return w.cfg.OutputKeys, w.cfg.OutputValues
	default:
		return "", ""
	}
}

// handleLocalExecution handles variable setting when not running in GitHub Actions.
// It prints values to the console instead of writing to a file.
func (w *Writer) handleLocalExecution(envVar, varType string, keyList, valueList []string) (int, error) {
	fmt.Printf(localExecMsg, envVar, varType)
	for i, key := range keyList {
		printer.PrintSuccess(varType, key, valueList[i])
	}
	return len(keyList), nil
}

// writeToFile writes key-value pairs to a file with retry logic.
// It attempts to write up to maxRetries times with delays between attempts.
func (w *Writer) writeToFile(filePath string, keys, values []string, varType string) (int, error) {
	maxRetries := 3
	retryDelay := time.Second
	var lastError error

	// Attempt writing with retries
	for retry := 0; retry < maxRetries; retry++ {
		count, err := w.performWrite(filePath, keys, values, varType)
		if err == nil {
			// Success - write action status
			_, _ = w.performWrite(filePath, []string{"action_status"}, []string{"success"}, varType)
			return count, nil
		}

		// Handle error
		lastError = err
		if retry < maxRetries-1 {
			printer.PrintError(fmt.Sprintf("Retry %d/%d: Failed to write to file: %v",
				retry+1, maxRetries, err))
			time.Sleep(retryDelay)
		}
	}

	// Write failure status after exhausting retries
	_, _ = w.performWrite(filePath,
		[]string{"action_status", "error_message"},
		[]string{"failure", lastError.Error()},
		varType)

	return 0, fmt.Errorf(errMaxRetries, maxRetries)
}

// performWrite writes key-value pairs to a file in GitHub Actions format.
// It handles file opening, value transformation, and formatting.
func (w *Writer) performWrite(filePath string, keys, values []string, varType string) (int, error) {
	// Open the file
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return 0, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Create transformer for values
	valueTransformer := transformer.New(
		w.cfg.MaskSecrets,
		w.cfg.MaskPattern,
		w.cfg.ToUpper,
		w.cfg.ToLower,
		w.cfg.EncodeURL,
		w.cfg.EscapeNewlines,
		w.cfg.MaxLength,
	)

	// Write header in debug mode
	if w.cfg.DebugMode {
		fmt.Printf("✍️  Writing Values:\n")
	}

	// Write each key-value pair
	count := 0
	for i, key := range keys {
		// Skip empty keys unless allowed
		if key == "" && !w.cfg.AllowEmpty {
			continue
		}

		// Apply whitespace trimming if configured
		if w.cfg.TrimWhitespace {
			key = strings.TrimSpace(key)
			values[i] = strings.TrimSpace(values[i])
		}

		// Transform and write the value
		transformedValue := valueTransformer.TransformValue(values[i], w.cfg.JsonSupport)
		if err := writeGitHubActionsFormat(file, key, transformedValue); err != nil {
			return count, err
		}

		// Print success message with masking
		maskedValue := valueTransformer.MaskValue(transformedValue)
		printer.PrintSuccess(varType, key, maskedValue)
		count++
	}

	// Write footer in debug mode
	if w.cfg.DebugMode {
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")
	}

	return count, nil
}

// writeGitHubActionsFormat writes a key-value pair in GitHub Actions format.
// Format: key<<EOF\nvalue\nEOF
func writeGitHubActionsFormat(file *os.File, key, value string) error {
	line := fmt.Sprintf("%s<<%s\n%s\n%s\n", key, "EOF", value, "EOF")
	if _, err := file.WriteString(line); err != nil {
		return fmt.Errorf("failed to write line: %w", err)
	}
	return nil
}
