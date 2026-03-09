package printer

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestPrintSection(t *testing.T) {
	tests := []struct {
		name     string
		title    string
		expected []string
	}{
		{
			name:  "Simple section title",
			title: "Test Section",
			expected: []string{
				DoubleLine,
				"Test Section",
			},
		},
		{
			name:  "Section with special characters",
			title: "Starting Process",
			expected: []string{
				DoubleLine,
				"Starting Process",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				PrintSection(tt.title)
			})

			for _, expected := range tt.expected {
				if !strings.Contains(output, expected) {
					t.Errorf("PrintSection() output missing expected string: %v", expected)
				}
			}
		})
	}
}

func TestPrintSuccess(t *testing.T) {
	tests := []struct {
		name     string
		varType  string
		key      string
		value    string
		expected []string
	}{
		{
			name:    "Environment variable success",
			varType: "env",
			key:     "API_KEY",
			value:   "secret123",
			expected: []string{
				BulletPoint,
				"env",
				"API_KEY",
				"secret123",
			},
		},
		{
			name:    "Output variable success",
			varType: "output",
			key:     "status",
			value:   "success",
			expected: []string{
				BulletPoint,
				"output",
				"status",
				"success",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				PrintSuccess(tt.varType, tt.key, tt.value)
			})

			for _, expected := range tt.expected {
				if !strings.Contains(output, expected) {
					t.Errorf("PrintSuccess() output missing expected string: %v", expected)
				}
			}
		})
	}
}

func TestPrintError(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected []string
	}{
		{
			name:    "Simple error message",
			message: "Failed to process",
			expected: []string{
				ErrorSymbol,
				"Failed to process",
			},
		},
		{
			name:    "Error with details",
			message: "Invalid configuration: missing required field",
			expected: []string{
				ErrorSymbol,
				"Invalid configuration",
				"missing required field",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				PrintError(tt.message)
			})

			for _, expected := range tt.expected {
				if !strings.Contains(output, expected) {
					t.Errorf("PrintError() output missing expected string: %v", expected)
				}
			}
		})
	}
}

func TestPrintWarning(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected []string
	}{
		{
			name:    "Simple warning message",
			message: "Deprecated feature used",
			expected: []string{
				WarningSymbol,
				"Deprecated feature used",
			},
		},
		{
			name:    "Warning with recommendation",
			message: "Value truncated: consider using max_length",
			expected: []string{
				WarningSymbol,
				"Value truncated",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				PrintWarning(tt.message)
			})

			for _, expected := range tt.expected {
				if !strings.Contains(output, expected) {
					t.Errorf("PrintWarning() output missing expected string: %v", expected)
				}
			}
		})
	}
}

func TestPrintInfo(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected string
	}{
		{
			name:     "Simple info message",
			message:  "Processing started",
			expected: "Processing started",
		},
		{
			name:     "Info with details",
			message:  "Mode: GitHub Actions",
			expected: "Mode: GitHub Actions",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				PrintInfo(tt.message)
			})

			if !strings.Contains(output, tt.expected) {
				t.Errorf("PrintInfo() output missing expected string: %v", tt.expected)
			}
		})
	}
}

func TestPrintDebugSection(t *testing.T) {
	tests := []struct {
		name     string
		title    string
		expected []string
	}{
		{
			name:  "Debug section with title",
			title: "Configuration",
			expected: []string{
				DebugSymbol,
				"Debug Information",
				"Configuration",
				SingleLine,
			},
		},
		{
			name:  "Debug section for values",
			title: "Processed Values",
			expected: []string{
				DebugSymbol,
				"Debug Information",
				"Processed Values",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				PrintDebugSection(tt.title)
			})

			for _, expected := range tt.expected {
				if !strings.Contains(output, expected) {
					t.Errorf("PrintDebugSection() output missing expected string: %v", expected)
				}
			}
		})
	}
}

func TestPrintDebugInfo(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		args     []interface{}
		expected string
	}{
		{
			name:     "Debug info with no arguments",
			format:   "Simple debug message\n",
			args:     []interface{}{},
			expected: "Simple debug message",
		},
		{
			name:     "Debug info with formatting",
			format:   "Key: %s, Value: %s\n",
			args:     []interface{}{"API_KEY", "***"},
			expected: "Key: API_KEY, Value: ***",
		},
		{
			name:     "Debug info with multiple types",
			format:   "Count: %d, Status: %s, Enabled: %v\n",
			args:     []interface{}{5, "active", true},
			expected: "Count: 5, Status: active, Enabled: true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				PrintDebugInfo(tt.format, tt.args...)
			})

			if !strings.Contains(output, tt.expected) {
				t.Errorf("PrintDebugInfo() output missing expected string: %v", tt.expected)
			}
		})
	}
}

func TestPrintDebugHighlight(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		args     []interface{}
		expected string
	}{
		{
			name:     "Highlighted simple message",
			format:   "Important: %s",
			args:     []interface{}{"Value changed"},
			expected: "Important: Value changed",
		},
		{
			name:     "Highlighted with number",
			format:   "Count: %d",
			args:     []interface{}{42},
			expected: "Count: 42",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				PrintDebugHighlight(tt.format, tt.args...)
			})

			if !strings.Contains(output, tt.expected) {
				t.Errorf("PrintDebugHighlight() output missing expected string: %v", tt.expected)
			}
		})
	}
}

func TestPrintComplete(t *testing.T) {
	output := captureOutput(func() {
		PrintComplete()
	})

	expected := []string{
		DoubleLine,
		SuccessSymbol,
		"Execution Complete",
		"Mode: GitHub Actions",
	}

	for _, exp := range expected {
		if !strings.Contains(output, exp) {
			t.Errorf("PrintComplete() output missing expected string: %v", exp)
		}
	}
}

func TestPrintLine(t *testing.T) {
	output := captureOutput(func() {
		PrintLine()
	})

	expected := strings.Repeat("=", 50)
	if !strings.Contains(output, expected) {
		t.Errorf("PrintLine() output should contain %v", expected)
	}
}

func TestPrintEmptyLine(t *testing.T) {
	output := captureOutput(func() {
		PrintEmptyLine()
	})

	if output != "\n" {
		t.Errorf("PrintEmptyLine() = %q, want %q", output, "\n")
	}
}

func TestFormatColor(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		color    string
		expected string
	}{
		{
			name:     "Format with success color",
			text:     "Success",
			color:    SuccessColor,
			expected: SuccessColor + "Success" + ResetColor,
		},
		{
			name:     "Format with error color",
			text:     "Error",
			color:    ErrorColor,
			expected: ErrorColor + "Error" + ResetColor,
		},
		{
			name:     "Format with info color",
			text:     "Information",
			color:    InfoColor,
			expected: InfoColor + "Information" + ResetColor,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatColor(tt.text, tt.color)
			if result != tt.expected {
				t.Errorf("FormatColor() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFormatSuccess(t *testing.T) {
	text := "Operation successful"
	expected := SuccessColor + text + ResetColor
	result := FormatSuccess(text)

	if result != expected {
		t.Errorf("FormatSuccess() = %v, want %v", result, expected)
	}
}

func TestFormatError(t *testing.T) {
	text := "Operation failed"
	expected := ErrorColor + text + ResetColor
	result := FormatError(text)

	if result != expected {
		t.Errorf("FormatError() = %v, want %v", result, expected)
	}
}

func TestFormatWarning(t *testing.T) {
	text := "Warning message"
	expected := WarningColor + text + ResetColor
	result := FormatWarning(text)

	if result != expected {
		t.Errorf("FormatWarning() = %v, want %v", result, expected)
	}
}

func TestFormatInfo(t *testing.T) {
	text := "Info message"
	expected := InfoColor + text + ResetColor
	result := FormatInfo(text)

	if result != expected {
		t.Errorf("FormatInfo() = %v, want %v", result, expected)
	}
}

func TestFormatDebug(t *testing.T) {
	text := "Debug message"
	expected := DebugColor + text + ResetColor
	result := FormatDebug(text)

	if result != expected {
		t.Errorf("FormatDebug() = %v, want %v", result, expected)
	}
}

func TestFormatHighlight(t *testing.T) {
	text := "Highlighted text"
	expected := HighlightColor + text + ResetColor
	result := FormatHighlight(text)

	if result != expected {
		t.Errorf("FormatHighlight() = %v, want %v", result, expected)
	}
}

func TestColorConstants(t *testing.T) {
	tests := []struct {
		name     string
		color    string
		notEmpty bool
	}{
		{"InfoColor", InfoColor, true},
		{"SuccessColor", SuccessColor, true},
		{"ErrorColor", ErrorColor, true},
		{"WarningColor", WarningColor, true},
		{"DebugColor", DebugColor, true},
		{"HeaderColor", HeaderColor, true},
		{"HighlightColor", HighlightColor, true},
		{"ResetColor", ResetColor, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.notEmpty && tt.color == "" {
				t.Errorf("%s should not be empty", tt.name)
			}
		})
	}
}

func TestSymbolConstants(t *testing.T) {
	tests := []struct {
		name     string
		symbol   string
		notEmpty bool
	}{
		{"SuccessSymbol", SuccessSymbol, true},
		{"ErrorSymbol", ErrorSymbol, true},
		{"WarningSymbol", WarningSymbol, true},
		{"InfoSymbol", InfoSymbol, true},
		{"DebugSymbol", DebugSymbol, true},
		{"BulletPoint", BulletPoint, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.notEmpty && tt.symbol == "" {
				t.Errorf("%s should not be empty", tt.name)
			}
		})
	}
}

func ExamplePrintSuccess() {
	PrintSuccess("env", "API_KEY", "***")
}

func ExamplePrintError() {
	PrintError("Failed to process configuration")
}

func ExamplePrintWarning() {
	PrintWarning("Deprecated feature used")
}

func ExampleFormatColor() {
	formatted := FormatColor("Important", SuccessColor)
	fmt.Println(formatted)
}
