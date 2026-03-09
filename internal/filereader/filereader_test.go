package filereader

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"
)

func TestIsFileReference(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{"File reference", "file:///tmp/data.txt", true},
		{"File reference relative", "file://config.json", true},
		{"Not a file reference", "just a value", false},
		{"Empty string", "", false},
		{"Partial prefix", "file:/", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsFileReference(tt.value); got != tt.expected {
				t.Errorf("IsFileReference(%q) = %v, want %v", tt.value, got, tt.expected)
			}
		})
	}
}

func TestGetFilePath(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{"Absolute path", "file:///tmp/data.txt", "/tmp/data.txt"},
		{"Relative path", "file://config.json", "config.json"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFilePath(tt.value); got != tt.expected {
				t.Errorf("GetFilePath(%q) = %q, want %q", tt.value, got, tt.expected)
			}
		})
	}
}

func TestReadValueRaw(t *testing.T) {
	t.Run("Read from file with raw encoding", func(t *testing.T) {
		tmpFile := filepath.Join(t.TempDir(), "test.txt")
		if err := os.WriteFile(tmpFile, []byte("hello world\n"), 0644); err != nil {
			t.Fatalf("failed to write test file: %v", err)
		}

		r := New(EncodingRaw)
		result, err := r.ReadValue("file://" + tmpFile)
		if err != nil {
			t.Fatalf("ReadValue() error = %v", err)
		}
		if result != "hello world" {
			t.Errorf("ReadValue() = %q, want %q", result, "hello world")
		}
	})

	t.Run("Non-file reference returns as-is", func(t *testing.T) {
		r := New(EncodingRaw)
		result, err := r.ReadValue("plain value")
		if err != nil {
			t.Fatalf("ReadValue() error = %v", err)
		}
		if result != "plain value" {
			t.Errorf("ReadValue() = %q, want %q", result, "plain value")
		}
	})

	t.Run("File not found returns error", func(t *testing.T) {
		r := New(EncodingRaw)
		_, err := r.ReadValue("file:///nonexistent/file.txt")
		if err == nil {
			t.Error("ReadValue() expected error for missing file, got nil")
		}
	})
}

func TestReadValueBase64(t *testing.T) {
	t.Run("Read base64 encoded file", func(t *testing.T) {
		encoded := base64.StdEncoding.EncodeToString([]byte("secret_value"))
		tmpFile := filepath.Join(t.TempDir(), "encoded.txt")
		if err := os.WriteFile(tmpFile, []byte(encoded+"\n"), 0644); err != nil {
			t.Fatalf("failed to write test file: %v", err)
		}

		r := New(EncodingBase64)
		result, err := r.ReadValue("file://" + tmpFile)
		if err != nil {
			t.Fatalf("ReadValue() error = %v", err)
		}
		if result != "secret_value" {
			t.Errorf("ReadValue() = %q, want %q", result, "secret_value")
		}
	})

	t.Run("Invalid base64 returns error", func(t *testing.T) {
		tmpFile := filepath.Join(t.TempDir(), "bad.txt")
		if err := os.WriteFile(tmpFile, []byte("not-valid-base64!!!"), 0644); err != nil {
			t.Fatalf("failed to write test file: %v", err)
		}

		r := New(EncodingBase64)
		_, err := r.ReadValue("file://" + tmpFile)
		if err == nil {
			t.Error("ReadValue() expected error for invalid base64, got nil")
		}
	})
}

func TestReadValues(t *testing.T) {
	tmpDir := t.TempDir()
	file1 := filepath.Join(tmpDir, "val1.txt")
	file2 := filepath.Join(tmpDir, "val2.txt")
	os.WriteFile(file1, []byte("from_file_1"), 0644)
	os.WriteFile(file2, []byte("from_file_2"), 0644)

	t.Run("Mixed file and plain values", func(t *testing.T) {
		r := New(EncodingRaw)
		values := []string{"file://" + file1, "plain_value", "file://" + file2}
		result, err := r.ReadValues(values)
		if err != nil {
			t.Fatalf("ReadValues() error = %v", err)
		}

		expected := []string{"from_file_1", "plain_value", "from_file_2"}
		for i, v := range result {
			if v != expected[i] {
				t.Errorf("ReadValues()[%d] = %q, want %q", i, v, expected[i])
			}
		}
	})

	t.Run("Error on missing file", func(t *testing.T) {
		r := New(EncodingRaw)
		values := []string{"ok", "file:///missing/file.txt"}
		_, err := r.ReadValues(values)
		if err == nil {
			t.Error("ReadValues() expected error, got nil")
		}
	})

	t.Run("Empty list", func(t *testing.T) {
		r := New(EncodingRaw)
		result, err := r.ReadValues([]string{})
		if err != nil {
			t.Fatalf("ReadValues() error = %v", err)
		}
		if len(result) != 0 {
			t.Errorf("ReadValues() length = %d, want 0", len(result))
		}
	})
}

func TestNewDefaultEncoding(t *testing.T) {
	r := New("")
	if r.encoding != EncodingRaw {
		t.Errorf("New(\"\").encoding = %q, want %q", r.encoding, EncodingRaw)
	}
}

func TestUnsupportedEncoding(t *testing.T) {
	r := New("gzip")
	tmpFile := filepath.Join(t.TempDir(), "test.txt")
	os.WriteFile(tmpFile, []byte("data"), 0644)

	_, err := r.ReadValue("file://" + tmpFile)
	if err == nil {
		t.Error("ReadValue() expected error for unsupported encoding, got nil")
	}
}

func BenchmarkReadValue(b *testing.B) {
	tmpFile := filepath.Join(b.TempDir(), "bench.txt")
	os.WriteFile(tmpFile, []byte("benchmark_value"), 0644)

	r := New(EncodingRaw)
	ref := "file://" + tmpFile

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.ReadValue(ref)
	}
}
