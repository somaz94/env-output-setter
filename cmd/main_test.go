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
