package filereader

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
)

// Supported encoding types
const (
	EncodingRaw    = "raw"
	EncodingBase64 = "base64"
)

// FilePrefix is the prefix that identifies a value as a file reference.
const FilePrefix = "file://"

// Reader handles reading values from files with optional encoding.
type Reader struct {
	encoding string
}

// New creates a new Reader with the specified encoding.
func New(encoding string) *Reader {
	if encoding == "" {
		encoding = EncodingRaw
	}
	return &Reader{encoding: strings.ToLower(encoding)}
}

// IsFileReference checks if a value is a file reference (starts with file://).
func IsFileReference(value string) bool {
	return strings.HasPrefix(value, FilePrefix)
}

// GetFilePath extracts the file path from a file reference.
func GetFilePath(value string) string {
	return strings.TrimPrefix(value, FilePrefix)
}

// ReadValue reads a value from a file if it's a file reference.
// If the value is not a file reference, it's returned as-is.
func (r *Reader) ReadValue(value string) (string, error) {
	if !IsFileReference(value) {
		return value, nil
	}

	filePath := GetFilePath(value)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	return r.decode(content)
}

// ReadValues processes a list of values, reading from files where applicable.
func (r *Reader) ReadValues(values []string) ([]string, error) {
	result := make([]string, len(values))
	for i, v := range values {
		read, err := r.ReadValue(v)
		if err != nil {
			return nil, fmt.Errorf("failed to read value at index %d: %w", i, err)
		}
		result[i] = read
	}
	return result, nil
}

// decode processes file content based on the configured encoding.
func (r *Reader) decode(content []byte) (string, error) {
	switch r.encoding {
	case EncodingBase64:
		decoded, err := base64.StdEncoding.DecodeString(strings.TrimSpace(string(content)))
		if err != nil {
			return "", fmt.Errorf("failed to decode base64 content: %w", err)
		}
		return string(decoded), nil
	case EncodingRaw:
		return strings.TrimSpace(string(content)), nil
	default:
		return "", fmt.Errorf("unsupported encoding: %s", r.encoding)
	}
}
