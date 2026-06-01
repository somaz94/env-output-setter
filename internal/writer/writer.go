package writer

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
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

// GitHub environment variables (aliased to the config package to avoid literal divergence)
const (
	githubEnvVar    = config.GithubEnvVar
	githubOutputVar = config.GithubOutputVar
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

	// Validate output values against rules if configured
	if err := w.validator.ValidateOutputs(keyList, valueList); err != nil {
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

// Status key names written to $GITHUB_ENV/$GITHUB_OUTPUT.
const (
	statusKey  = "action_status"
	errMsgKey  = "error_message"
	statusOK   = "success"
	statusFail = "failure"
)

// writeToFile writes key-value pairs to a file with retry logic.
// It builds the full payload in a buffer first and appends atomically per attempt,
// so a failed attempt never leaves partial lines behind for the next retry to duplicate.
func (w *Writer) writeToFile(filePath string, keys, values []string, varType string) (int, error) {
	maxRetries := 3
	retryDelay := time.Second
	var lastError error

	for retry := 0; retry < maxRetries; retry++ {
		count, err := w.performWrite(filePath, keys, values, varType)
		if err == nil {
			// Success - write action status (best-effort)
			if _, statusErr := w.performWrite(filePath, []string{statusKey}, []string{statusOK}, varType); statusErr != nil {
				printer.PrintWarning(fmt.Sprintf("Warning: failed to write success status: %v", statusErr))
			}
			return count, nil
		}

		lastError = err
		if retry < maxRetries-1 {
			printer.PrintError(fmt.Sprintf("Retry %d/%d: Failed to write to file: %v",
				retry+1, maxRetries, err))
			time.Sleep(retryDelay)
		}
	}

	// Write failure status after exhausting retries (best-effort)
	failMsg := ""
	if lastError != nil {
		failMsg = lastError.Error()
	}
	if _, err := w.performWrite(filePath,
		[]string{statusKey, errMsgKey},
		[]string{statusFail, failMsg},
		varType); err != nil {
		printer.PrintWarning(fmt.Sprintf("Warning: failed to write failure status: %v", err))
	}

	return 0, fmt.Errorf(errMaxRetries, maxRetries)
}

// performWrite writes key-value pairs to a file in GitHub Actions format.
// All lines are serialized into an in-memory buffer first and flushed in a single
// WriteString call so a partial failure leaves the file untouched (atomicity per call).
// The file's Close error is propagated via the named return so disk-full / NFS errors
// surface to the caller instead of being silently discarded.
func (w *Writer) performWrite(filePath string, keys, values []string, varType string) (count int, err error) {
	// Build the full payload in memory before opening the file.
	valueTransformer := transformer.New(
		w.cfg.MaskSecrets,
		w.cfg.MaskPattern,
		w.cfg.ToUpper,
		w.cfg.ToLower,
		w.cfg.EncodeURL,
		w.cfg.EscapeNewlines,
		w.cfg.MaxLength,
	)

	if w.cfg.DebugMode {
		fmt.Printf("Writing Values:\n")
	}

	var buf bytes.Buffer
	type successMsg struct{ key, masked string }
	var successes []successMsg

	for i, key := range keys {
		// Skip empty keys unless allowed
		if key == "" && !w.cfg.AllowEmpty {
			continue
		}

		// Apply whitespace trimming if configured (non-mutating: trim local copies only)
		k := key
		v := values[i]
		if w.cfg.TrimWhitespace {
			k = strings.TrimSpace(k)
			v = strings.TrimSpace(v)
		}

		transformedValue := valueTransformer.TransformValue(v, w.cfg.JsonSupport)
		if werr := appendGitHubActionsFormat(&buf, k, transformedValue); werr != nil {
			return 0, werr
		}

		if k != statusKey && k != errMsgKey {
			successes = append(successes, successMsg{k, valueTransformer.MaskValue(transformedValue)})
		}
		count++
	}

	// Now open the file and flush the buffer in one write.
	file, openErr := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if openErr != nil {
		return 0, fmt.Errorf("failed to open file: %w", openErr)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("failed to close file: %w", cerr)
			count = 0
		}
	}()

	if _, werr := file.Write(buf.Bytes()); werr != nil {
		return 0, fmt.Errorf("failed to write payload: %w", werr)
	}
	if serr := file.Sync(); serr != nil {
		return 0, fmt.Errorf("failed to sync file: %w", serr)
	}

	for _, s := range successes {
		printer.PrintSuccess(varType, s.key, s.masked)
	}

	if w.cfg.DebugMode {
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")
	}

	return count, nil
}

// appendGitHubActionsFormat appends a key-value pair to buf in GitHub Actions multiline format.
// Uses a random delimiter to avoid collisions with value content.
func appendGitHubActionsFormat(buf *bytes.Buffer, key, value string) error {
	delimiter, err := randomDelimiter()
	if err != nil {
		return fmt.Errorf("failed to generate delimiter: %w", err)
	}
	fmt.Fprintf(buf, "%s<<%s\n%s\n%s\n", key, delimiter, value, delimiter)
	return nil
}

// randomDelimiter generates a unique delimiter for GitHub Actions multiline format.
func randomDelimiter() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return "EOF_" + hex.EncodeToString(b), nil
}
