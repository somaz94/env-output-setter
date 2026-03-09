package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAppendToFile(t *testing.T) {
	t.Run("appends content to file", func(t *testing.T) {
		tmpFile := filepath.Join(t.TempDir(), "test_output")
		if err := appendToFile(tmpFile, "key1=value1"); err != nil {
			t.Fatalf("appendToFile() error = %v", err)
		}
		if err := appendToFile(tmpFile, "key2=value2"); err != nil {
			t.Fatalf("appendToFile() error = %v", err)
		}

		data, err := os.ReadFile(tmpFile)
		if err != nil {
			t.Fatalf("failed to read file: %v", err)
		}

		content := string(data)
		if !strings.Contains(content, "key1=value1") {
			t.Errorf("expected content to contain 'key1=value1', got %q", content)
		}
		if !strings.Contains(content, "key2=value2") {
			t.Errorf("expected content to contain 'key2=value2', got %q", content)
		}
	})

	t.Run("creates file if not exists", func(t *testing.T) {
		tmpFile := filepath.Join(t.TempDir(), "new_file")
		if err := appendToFile(tmpFile, "hello=world"); err != nil {
			t.Fatalf("appendToFile() error = %v", err)
		}

		data, err := os.ReadFile(tmpFile)
		if err != nil {
			t.Fatalf("failed to read file: %v", err)
		}
		if !strings.Contains(string(data), "hello=world") {
			t.Errorf("expected 'hello=world' in file content")
		}
	})

	t.Run("returns error for invalid path", func(t *testing.T) {
		err := appendToFile("/nonexistent/dir/file", "content")
		if err == nil {
			t.Error("expected error for invalid path, got nil")
		}
		if !strings.Contains(err.Error(), "failed to open file") {
			t.Errorf("expected 'failed to open file' error, got %v", err)
		}
	})
}

func TestWriteOutputs(t *testing.T) {
	t.Run("writes outputs to file", func(t *testing.T) {
		tmpFile := filepath.Join(t.TempDir(), "github_output")
		t.Setenv("GITHUB_OUTPUT", tmpFile)

		writeOutputs(3, 2, "success", "")

		data, err := os.ReadFile(tmpFile)
		if err != nil {
			t.Fatalf("failed to read file: %v", err)
		}

		content := string(data)
		if !strings.Contains(content, "set_env_count=3") {
			t.Errorf("expected 'set_env_count=3' in output, got %q", content)
		}
		if !strings.Contains(content, "set_output_count=2") {
			t.Errorf("expected 'set_output_count=2' in output, got %q", content)
		}
		if !strings.Contains(content, "status=success") {
			t.Errorf("expected 'status=success' in output, got %q", content)
		}
	})

	t.Run("skips when GITHUB_OUTPUT is not set", func(t *testing.T) {
		t.Setenv("GITHUB_OUTPUT", "")

		// Should not panic or error
		writeOutputs(1, 1, "success", "")
	})

	t.Run("writes failure status with error message", func(t *testing.T) {
		tmpFile := filepath.Join(t.TempDir(), "github_output")
		t.Setenv("GITHUB_OUTPUT", tmpFile)

		writeOutputs(0, 0, "failure", "something went wrong")

		data, err := os.ReadFile(tmpFile)
		if err != nil {
			t.Fatalf("failed to read file: %v", err)
		}

		content := string(data)
		if !strings.Contains(content, "status=failure") {
			t.Errorf("expected 'status=failure' in output, got %q", content)
		}
		if !strings.Contains(content, "error_message=something went wrong") {
			t.Errorf("expected error message in output, got %q", content)
		}
	})

	t.Run("handles invalid output path gracefully", func(t *testing.T) {
		t.Setenv("GITHUB_OUTPUT", "/nonexistent/dir/output")

		// Should not panic - just prints error
		writeOutputs(1, 1, "success", "")
	})
}

func TestRun(t *testing.T) {
	t.Run("runs successfully in local mode", func(t *testing.T) {
		// Ensure we're in local mode (no GITHUB_ENV/OUTPUT)
		t.Setenv("GITHUB_ENV", "")
		t.Setenv("GITHUB_OUTPUT", "")
		t.Setenv("INPUT_ENV_KEY", "TEST_KEY")
		t.Setenv("INPUT_ENV_VALUE", "test_value")
		t.Setenv("INPUT_OUTPUT_KEY", "OUT_KEY")
		t.Setenv("INPUT_OUTPUT_VALUE", "out_value")
		t.Setenv("INPUT_DELIMITER", ",")

		exitCode := run()
		if exitCode != 0 {
			t.Errorf("expected exit code 0, got %d", exitCode)
		}
	})

	t.Run("runs with debug mode", func(t *testing.T) {
		t.Setenv("GITHUB_ENV", "")
		t.Setenv("GITHUB_OUTPUT", "")
		t.Setenv("INPUT_ENV_KEY", "KEY1")
		t.Setenv("INPUT_ENV_VALUE", "val1")
		t.Setenv("INPUT_OUTPUT_KEY", "OUT1")
		t.Setenv("INPUT_OUTPUT_VALUE", "out1")
		t.Setenv("INPUT_DELIMITER", ",")
		t.Setenv("INPUT_DEBUG_MODE", "true")
		t.Setenv("INPUT_GROUP_PREFIX", "APP")
		t.Setenv("INPUT_JSON_SUPPORT", "true")
		t.Setenv("INPUT_EXPORT_AS_ENV", "true")

		exitCode := run()
		if exitCode != 0 {
			t.Errorf("expected exit code 0, got %d", exitCode)
		}
	})

	t.Run("runs in github actions mode", func(t *testing.T) {
		tmpEnv := filepath.Join(t.TempDir(), "github_env")
		tmpOutput := filepath.Join(t.TempDir(), "github_output")
		t.Setenv("GITHUB_ENV", tmpEnv)
		t.Setenv("GITHUB_OUTPUT", tmpOutput)
		t.Setenv("INPUT_ENV_KEY", "MY_KEY")
		t.Setenv("INPUT_ENV_VALUE", "my_value")
		t.Setenv("INPUT_OUTPUT_KEY", "MY_OUT")
		t.Setenv("INPUT_OUTPUT_VALUE", "my_out_val")
		t.Setenv("INPUT_DELIMITER", ",")

		exitCode := run()
		if exitCode != 0 {
			t.Errorf("expected exit code 0, got %d", exitCode)
		}

		// Verify env file was written
		envData, err := os.ReadFile(tmpEnv)
		if err != nil {
			t.Fatalf("failed to read env file: %v", err)
		}
		if !strings.Contains(string(envData), "MY_KEY") {
			t.Errorf("expected env file to contain 'MY_KEY'")
		}

		// Verify output file was written
		outData, err := os.ReadFile(tmpOutput)
		if err != nil {
			t.Fatalf("failed to read output file: %v", err)
		}
		if !strings.Contains(string(outData), "MY_OUT") {
			t.Errorf("expected output file to contain 'MY_OUT'")
		}
	})

	t.Run("returns 1 on env write error", func(t *testing.T) {
		t.Setenv("GITHUB_ENV", "/nonexistent/path/env")
		t.Setenv("GITHUB_OUTPUT", "")
		t.Setenv("INPUT_ENV_KEY", "KEY")
		t.Setenv("INPUT_ENV_VALUE", "val")
		t.Setenv("INPUT_OUTPUT_KEY", "")
		t.Setenv("INPUT_OUTPUT_VALUE", "")
		t.Setenv("INPUT_DELIMITER", ",")

		exitCode := run()
		if exitCode != 1 {
			t.Errorf("expected exit code 1, got %d", exitCode)
		}
	})

	t.Run("returns 1 on output write error", func(t *testing.T) {
		tmpEnv := filepath.Join(t.TempDir(), "github_env")
		t.Setenv("GITHUB_ENV", tmpEnv)
		t.Setenv("GITHUB_OUTPUT", "/nonexistent/path/output")
		t.Setenv("INPUT_ENV_KEY", "KEY")
		t.Setenv("INPUT_ENV_VALUE", "val")
		t.Setenv("INPUT_OUTPUT_KEY", "OUT")
		t.Setenv("INPUT_OUTPUT_VALUE", "out")
		t.Setenv("INPUT_DELIMITER", ",")

		exitCode := run()
		if exitCode != 1 {
			t.Errorf("expected exit code 1, got %d", exitCode)
		}
	})
}
