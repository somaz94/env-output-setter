package writer

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/somaz94/env-output-setter/internal/config"
)

func TestNewWriter(t *testing.T) {
	cfg := &config.Config{
		EnvKeys:   "KEY1",
		EnvValues: "VALUE1",
		Delimiter: ",",
	}

	writer := NewWriter(cfg)

	if writer == nil {
		t.Fatal("NewWriter() returned nil")
	}

	if writer.cfg != cfg {
		t.Error("NewWriter() config not set correctly")
	}

	if writer.processor == nil {
		t.Error("NewWriter() processor not initialized")
	}

	if writer.validator == nil {
		t.Error("NewWriter() validator not initialized")
	}
}

func TestGetInputValues(t *testing.T) {
	tests := []struct {
		name         string
		envVar       string
		cfg          *config.Config
		expectedKeys string
		expectedVals string
	}{
		{
			name:   "Get env variables",
			envVar: githubEnvVar,
			cfg: &config.Config{
				EnvKeys:   "ENV_KEY1,ENV_KEY2",
				EnvValues: "ENV_VAL1,ENV_VAL2",
			},
			expectedKeys: "ENV_KEY1,ENV_KEY2",
			expectedVals: "ENV_VAL1,ENV_VAL2",
		},
		{
			name:   "Get output variables",
			envVar: githubOutputVar,
			cfg: &config.Config{
				OutputKeys:   "OUT_KEY1,OUT_KEY2",
				OutputValues: "OUT_VAL1,OUT_VAL2",
			},
			expectedKeys: "OUT_KEY1,OUT_KEY2",
			expectedVals: "OUT_VAL1,OUT_VAL2",
		},
		{
			name:         "Unknown variable type",
			envVar:       "UNKNOWN_VAR",
			cfg:          &config.Config{},
			expectedKeys: "",
			expectedVals: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := NewWriter(tt.cfg)
			keys, values := writer.getInputValues(tt.envVar)

			if keys != tt.expectedKeys {
				t.Errorf("getInputValues() keys = %v, want %v", keys, tt.expectedKeys)
			}

			if values != tt.expectedVals {
				t.Errorf("getInputValues() values = %v, want %v", values, tt.expectedVals)
			}
		})
	}
}

func TestSetEnvLocalExecution(t *testing.T) {
	// Ensure GITHUB_ENV is not set
	os.Unsetenv(githubEnvVar)

	cfg := &config.Config{
		EnvKeys:        "TEST_KEY",
		EnvValues:      "TEST_VALUE",
		Delimiter:      ",",
		FailOnEmpty:    true,
		TrimWhitespace: true,
	}

	count, err := SetEnv(cfg)

	if err != nil {
		t.Errorf("SetEnv() local execution should not error: %v", err)
	}

	if count != 1 {
		t.Errorf("SetEnv() count = %d, want 1", count)
	}
}

func TestSetOutputLocalExecution(t *testing.T) {
	// Ensure GITHUB_OUTPUT is not set
	os.Unsetenv(githubOutputVar)

	cfg := &config.Config{
		OutputKeys:     "STATUS",
		OutputValues:   "SUCCESS",
		Delimiter:      ",",
		FailOnEmpty:    true,
		TrimWhitespace: true,
	}

	count, err := SetOutput(cfg)

	if err != nil {
		t.Errorf("SetOutput() local execution should not error: %v", err)
	}

	if count != 1 {
		t.Errorf("SetOutput() count = %d, want 1", count)
	}
}

func TestSetEnvWithFile(t *testing.T) {
	// Create temporary file
	tmpDir := t.TempDir()
	envFile := filepath.Join(tmpDir, "github_env")

	// Set GITHUB_ENV
	os.Setenv(githubEnvVar, envFile)
	defer os.Unsetenv(githubEnvVar)

	cfg := &config.Config{
		EnvKeys:        "TEST_KEY1,TEST_KEY2",
		EnvValues:      "TEST_VALUE1,TEST_VALUE2",
		Delimiter:      ",",
		FailOnEmpty:    true,
		TrimWhitespace: true,
		GithubEnv:      envFile,
	}

	count, err := SetEnv(cfg)

	if err != nil {
		t.Errorf("SetEnv() unexpected error: %v", err)
	}

	if count != 2 {
		t.Errorf("SetEnv() count = %d, want 2", count)
	}

	// Verify file exists
	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		t.Error("SetEnv() did not create the env file")
	}

	// Read and verify file content
	content, err := os.ReadFile(envFile)
	if err != nil {
		t.Fatalf("Failed to read env file: %v", err)
	}

	contentStr := string(content)
	if !contains(contentStr, "TEST_KEY1") {
		t.Error("SetEnv() file missing TEST_KEY1")
	}
	if !contains(contentStr, "TEST_VALUE1") {
		t.Error("SetEnv() file missing TEST_VALUE1")
	}
}

func TestSetOutputWithFile(t *testing.T) {
	// Create temporary file
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "github_output")

	// Set GITHUB_OUTPUT
	os.Setenv(githubOutputVar, outputFile)
	defer os.Unsetenv(githubOutputVar)

	cfg := &config.Config{
		OutputKeys:     "STATUS,COUNT",
		OutputValues:   "SUCCESS,5",
		Delimiter:      ",",
		FailOnEmpty:    true,
		TrimWhitespace: true,
		GithubOutput:   outputFile,
	}

	count, err := SetOutput(cfg)

	if err != nil {
		t.Errorf("SetOutput() unexpected error: %v", err)
	}

	if count != 2 {
		t.Errorf("SetOutput() count = %d, want 2", count)
	}

	// Verify file exists
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Error("SetOutput() did not create the output file")
	}

	// Read and verify file content
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	contentStr := string(content)
	if !contains(contentStr, "STATUS") {
		t.Error("SetOutput() file missing STATUS")
	}
	if !contains(contentStr, "SUCCESS") {
		t.Error("SetOutput() file missing SUCCESS")
	}
}

func TestSetEnvValidationError(t *testing.T) {
	tmpDir := t.TempDir()
	envFile := filepath.Join(tmpDir, "github_env")
	os.Setenv(githubEnvVar, envFile)
	defer os.Unsetenv(githubEnvVar)

	cfg := &config.Config{
		EnvKeys:        "KEY1,KEY2",
		EnvValues:      "VALUE1", // Mismatched count
		Delimiter:      ",",
		FailOnEmpty:    true,
		TrimWhitespace: true,
		GithubEnv:      envFile,
	}

	count, err := SetEnv(cfg)

	if err == nil {
		t.Error("SetEnv() expected error for mismatched pairs, got nil")
	}

	if count != 0 {
		t.Errorf("SetEnv() count = %d, want 0 on error", count)
	}
}

func TestSetOutputWithExportAsEnv(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "github_output")
	envFile := filepath.Join(tmpDir, "github_env")

	os.Setenv(githubOutputVar, outputFile)
	os.Setenv(githubEnvVar, envFile)
	defer os.Unsetenv(githubOutputVar)
	defer os.Unsetenv(githubEnvVar)

	cfg := &config.Config{
		OutputKeys:     "STATUS",
		OutputValues:   "SUCCESS",
		Delimiter:      ",",
		FailOnEmpty:    true,
		TrimWhitespace: true,
		ExportAsEnv:    true,
		GithubOutput:   outputFile,
		GithubEnv:      envFile,
	}

	count, err := SetOutput(cfg)

	if err != nil {
		t.Errorf("SetOutput() with export_as_env unexpected error: %v", err)
	}

	if count < 1 {
		t.Errorf("SetOutput() count = %d, want at least 1", count)
	}

	// Verify both files exist
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Error("SetOutput() did not create the output file")
	}

	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		t.Error("SetOutput() with export_as_env did not create the env file")
	}
}

func TestWriteGitHubActionsFormat(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test_output")

	file, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer file.Close()

	err = writeGitHubActionsFormat(file, "TEST_KEY", "TEST_VALUE")
	if err != nil {
		t.Errorf("writeGitHubActionsFormat() unexpected error: %v", err)
	}

	file.Close()

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	expected := "TEST_KEY<<EOF\nTEST_VALUE\nEOF\n"
	if string(content) != expected {
		t.Errorf("writeGitHubActionsFormat() content = %q, want %q", string(content), expected)
	}
}

func TestSetEnvWithTransformations(t *testing.T) {
	tmpDir := t.TempDir()
	envFile := filepath.Join(tmpDir, "github_env")
	os.Setenv(githubEnvVar, envFile)
	defer os.Unsetenv(githubEnvVar)

	cfg := &config.Config{
		EnvKeys:        "test_key",
		EnvValues:      "test value",
		Delimiter:      ",",
		FailOnEmpty:    true,
		TrimWhitespace: true,
		ToUpper:        true,
		GithubEnv:      envFile,
	}

	count, err := SetEnv(cfg)

	if err != nil {
		t.Errorf("SetEnv() with transformations unexpected error: %v", err)
	}

	if count != 1 {
		t.Errorf("SetEnv() count = %d, want 1", count)
	}

	content, err := os.ReadFile(envFile)
	if err != nil {
		t.Fatalf("Failed to read env file: %v", err)
	}

	if !contains(string(content), "TEST VALUE") {
		t.Error("SetEnv() value not transformed to uppercase")
	}
}

func TestSetEnvEmptyInput(t *testing.T) {
	tmpDir := t.TempDir()
	envFile := filepath.Join(tmpDir, "github_env")
	os.Setenv(githubEnvVar, envFile)
	defer os.Unsetenv(githubEnvVar)

	cfg := &config.Config{
		EnvKeys:        "",
		EnvValues:      "",
		Delimiter:      ",",
		FailOnEmpty:    false,
		TrimWhitespace: true,
		GithubEnv:      envFile,
	}

	count, err := SetEnv(cfg)

	if err != nil {
		t.Errorf("SetEnv() with empty input unexpected error: %v", err)
	}

	if count != 0 {
		t.Errorf("SetEnv() count = %d, want 0 for empty input", count)
	}
}

func TestSetOutputEmptyInput(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "github_output")
	os.Setenv(githubOutputVar, outputFile)
	defer os.Unsetenv(githubOutputVar)

	cfg := &config.Config{
		OutputKeys:     "",
		OutputValues:   "",
		Delimiter:      ",",
		FailOnEmpty:    false,
		TrimWhitespace: true,
		GithubOutput:   outputFile,
	}

	count, err := SetOutput(cfg)

	if err != nil {
		t.Errorf("SetOutput() with empty input unexpected error: %v", err)
	}

	if count != 0 {
		t.Errorf("SetOutput() count = %d, want 0 for empty input", count)
	}
}

func TestSetEnvWithDebugMode(t *testing.T) {
	tmpDir := t.TempDir()
	envFile := filepath.Join(tmpDir, "github_env")
	os.Setenv(githubEnvVar, envFile)
	defer os.Unsetenv(githubEnvVar)

	cfg := &config.Config{
		EnvKeys:        "KEY1,KEY2",
		EnvValues:      "VALUE1,VALUE2",
		Delimiter:      ",",
		FailOnEmpty:    true,
		TrimWhitespace: true,
		DebugMode:      true,
		GithubEnv:      envFile,
	}

	count, err := SetEnv(cfg)

	if err != nil {
		t.Errorf("SetEnv() with debug mode unexpected error: %v", err)
	}

	if count != 2 {
		t.Errorf("SetEnv() count = %d, want 2", count)
	}
}

func TestSetEnvWithJsonSupport(t *testing.T) {
	tmpDir := t.TempDir()
	envFile := filepath.Join(tmpDir, "github_env")
	os.Setenv(githubEnvVar, envFile)
	defer os.Unsetenv(githubEnvVar)

	cfg := &config.Config{
		EnvKeys:        "CONFIG",
		EnvValues:      `{"host":"localhost"}`,
		Delimiter:      "|",
		TrimWhitespace: true,
		JsonSupport:    true,
		GithubEnv:      envFile,
	}

	count, err := SetEnv(cfg)

	if err != nil {
		t.Errorf("SetEnv() with JSON support unexpected error: %v", err)
	}

	if count < 1 {
		t.Errorf("SetEnv() count = %d, want at least 1", count)
	}

	content, err := os.ReadFile(envFile)
	if err != nil {
		t.Fatalf("Failed to read env file: %v", err)
	}

	contentStr := string(content)
	if !contains(contentStr, "CONFIG") {
		t.Error("SetEnv() file missing CONFIG key")
	}
}

func TestSetEnvWithAllowEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	envFile := filepath.Join(tmpDir, "github_env")
	os.Setenv(githubEnvVar, envFile)
	defer os.Unsetenv(githubEnvVar)

	cfg := &config.Config{
		EnvKeys:        ",KEY2",
		EnvValues:      "VALUE1,VALUE2",
		Delimiter:      ",",
		FailOnEmpty:    false,
		AllowEmpty:     true,
		TrimWhitespace: true,
		GithubEnv:      envFile,
	}

	count, err := SetEnv(cfg)

	if err != nil {
		t.Errorf("SetEnv() with allow empty unexpected error: %v", err)
	}

	if count < 1 {
		t.Errorf("SetEnv() count = %d, want at least 1", count)
	}
}

func TestWriteToFileInvalidPath(t *testing.T) {
	cfg := &config.Config{
		EnvKeys:        "KEY1",
		EnvValues:      "VALUE1",
		Delimiter:      ",",
		TrimWhitespace: true,
		GithubEnv:      "/nonexistent/path/file",
	}

	os.Setenv(githubEnvVar, "/nonexistent/path/file")
	defer os.Unsetenv(githubEnvVar)

	_, err := SetEnv(cfg)

	if err == nil {
		t.Error("SetEnv() expected error for invalid file path, got nil")
	}
}

func TestSetOutputWithExportAsEnvLocalExecution(t *testing.T) {
	// No GITHUB_OUTPUT or GITHUB_ENV set
	os.Unsetenv(githubOutputVar)
	os.Unsetenv(githubEnvVar)

	cfg := &config.Config{
		OutputKeys:     "STATUS",
		OutputValues:   "SUCCESS",
		Delimiter:      ",",
		TrimWhitespace: true,
		ExportAsEnv:    true,
	}

	count, err := SetOutput(cfg)

	if err != nil {
		t.Errorf("SetOutput() with export_as_env local execution unexpected error: %v", err)
	}

	if count < 1 {
		t.Errorf("SetOutput() count = %d, want at least 1", count)
	}
}

func TestSetEnvValidateInputsError(t *testing.T) {
	tmpDir := t.TempDir()
	envFile := filepath.Join(tmpDir, "github_env")
	os.Setenv(githubEnvVar, envFile)
	defer os.Unsetenv(githubEnvVar)

	cfg := &config.Config{
		EnvKeys:          "KEY1,KEY1",
		EnvValues:        "VALUE1,VALUE2",
		Delimiter:        ",",
		TrimWhitespace:   true,
		ErrorOnDuplicate: true,
		CaseSensitive:    true,
		GithubEnv:        envFile,
	}

	_, err := SetEnv(cfg)

	if err == nil {
		t.Error("SetEnv() expected error for duplicate keys, got nil")
	}
}

func TestSetOutputValidationError(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "github_output")
	os.Setenv(githubOutputVar, outputFile)
	defer os.Unsetenv(githubOutputVar)

	cfg := &config.Config{
		OutputKeys:     "KEY1,KEY2",
		OutputValues:   "VALUE1",
		Delimiter:      ",",
		TrimWhitespace: true,
		GithubOutput:   outputFile,
	}

	count, err := SetOutput(cfg)

	if err == nil {
		t.Error("SetOutput() expected error for mismatched pairs, got nil")
	}

	if count != 0 {
		t.Errorf("SetOutput() count = %d, want 0 on error", count)
	}
}

func TestSetOutputWithExportAsEnvWriteError(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "github_output")

	os.Setenv(githubOutputVar, outputFile)
	os.Setenv(githubEnvVar, "/nonexistent/path/env_file")
	defer os.Unsetenv(githubOutputVar)
	defer os.Unsetenv(githubEnvVar)

	cfg := &config.Config{
		OutputKeys:     "STATUS",
		OutputValues:   "SUCCESS",
		Delimiter:      ",",
		TrimWhitespace: true,
		ExportAsEnv:    true,
		GithubOutput:   outputFile,
		GithubEnv:      "/nonexistent/path/env_file",
	}

	_, err := SetOutput(cfg)

	if err == nil {
		t.Error("SetOutput() with export_as_env expected error for invalid env path, got nil")
	}
}

func TestSetEnvWithMasking(t *testing.T) {
	tmpDir := t.TempDir()
	envFile := filepath.Join(tmpDir, "github_env")
	os.Setenv(githubEnvVar, envFile)
	defer os.Unsetenv(githubEnvVar)

	cfg := &config.Config{
		EnvKeys:        "API_KEY,USERNAME",
		EnvValues:      "secret123,admin",
		Delimiter:      ",",
		TrimWhitespace: true,
		MaskSecrets:    true,
		MaskPattern:    "secret.*",
		GithubEnv:      envFile,
	}

	count, err := SetEnv(cfg)

	if err != nil {
		t.Errorf("SetEnv() with masking unexpected error: %v", err)
	}

	if count != 2 {
		t.Errorf("SetEnv() count = %d, want 2", count)
	}
}

func BenchmarkSetEnv(b *testing.B) {
	tmpDir := b.TempDir()
	envFile := filepath.Join(tmpDir, "github_env")
	os.Setenv(githubEnvVar, envFile)
	defer os.Unsetenv(githubEnvVar)

	cfg := &config.Config{
		EnvKeys:        "KEY1,KEY2,KEY3",
		EnvValues:      "VALUE1,VALUE2,VALUE3",
		Delimiter:      ",",
		FailOnEmpty:    true,
		TrimWhitespace: true,
		GithubEnv:      envFile,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Clean file for each iteration
		os.Remove(envFile)
		SetEnv(cfg)
	}
}

func BenchmarkSetOutput(b *testing.B) {
	tmpDir := b.TempDir()
	outputFile := filepath.Join(tmpDir, "github_output")
	os.Setenv(githubOutputVar, outputFile)
	defer os.Unsetenv(githubOutputVar)

	cfg := &config.Config{
		OutputKeys:     "STATUS,COUNT,MESSAGE",
		OutputValues:   "SUCCESS,5,Done",
		Delimiter:      ",",
		FailOnEmpty:    true,
		TrimWhitespace: true,
		GithubOutput:   outputFile,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Clean file for each iteration
		os.Remove(outputFile)
		SetOutput(cfg)
	}
}
