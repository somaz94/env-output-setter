package printer

import (
	"fmt"
	"strings"
)

// Color constants for terminal output formatting
const (
	// Basic colors
	InfoColor    = "\033[1;34m" // Blue
	SuccessColor = "\033[1;32m" // Green
	ErrorColor   = "\033[1;31m" // Red
	WarningColor = "\033[1;33m" // Yellow
	DebugColor   = "\033[1;36m" // Cyan
	ResetColor   = "\033[0m"    // Reset color

	// Additional colors
	HeaderColor    = "\033[1;35m" // Purple
	HighlightColor = "\033[1;37m" // Light White
)

// Separator constants for visual organization
const (
	DoubleLine = "=================================================="
	SingleLine = "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
)

// Common output symbols
const (
	SuccessSymbol = "✅"
	ErrorSymbol   = "❌"
	WarningSymbol = "⚠️"
	InfoSymbol    = "ℹ️"
	DebugSymbol   = "🔍"
	BulletPoint   = "•"
)

// PrintSection prints a section header with a double line border.
// It creates a visual separation for a new section of output.
func PrintSection(title string) {
	fmt.Printf("\n%s%s\n", InfoColor, DoubleLine)
	fmt.Printf("%s%s%s\n", InfoColor, title, ResetColor)
}

// PrintSuccess prints a success message for a variable assignment.
// It formats the output with the variable type, key, and value.
func PrintSuccess(varType, key, value string) {
	fmt.Printf("%s  %s %s: %s = %s%s\n",
		SuccessColor,
		BulletPoint,
		varType,
		key,
		value,
		ResetColor)
}

// PrintError prints an error message with a distinguishing symbol.
func PrintError(message string) {
	fmt.Printf("%s%s %s%s\n",
		ErrorColor,
		ErrorSymbol,
		message,
		ResetColor)
}

// PrintWarning prints a warning message with a warning symbol.
func PrintWarning(message string) {
	fmt.Printf("%s%s %s%s\n",
		WarningColor,
		WarningSymbol,
		message,
		ResetColor)
}

// PrintInfo prints an informational message in blue color.
func PrintInfo(message string) {
	fmt.Printf("%s%s%s\n",
		InfoColor,
		message,
		ResetColor)
}

// PrintDebugSection prints a debug section header with a title.
// It creates a visual separation for debug information.
func PrintDebugSection(title string) {
	fmt.Printf("\n%s%s Debug Information (%s)\n",
		DebugColor,
		DebugSymbol,
		title)
	fmt.Printf("%s%s\n", SingleLine, ResetColor)
}

// PrintDebugInfo prints debug information without special formatting.
// It accepts a format string and variadic arguments like fmt.Printf.
func PrintDebugInfo(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

// PrintDebugHighlight prints highlighted debug information in light white color.
// It accepts a format string and variadic arguments like fmt.Printf.
func PrintDebugHighlight(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	fmt.Printf("%s%s%s", HighlightColor, message, ResetColor)
}

// PrintComplete prints a completion message with execution mode information.
func PrintComplete() {
	fmt.Printf("\n%s%s\n", InfoColor, DoubleLine)
	fmt.Printf("%s%s Execution Complete\n", InfoColor, SuccessSymbol)
	fmt.Printf("Mode: GitHub Actions%s\n", ResetColor)
}

// PrintLine prints a horizontal line of 50 equal signs for visual separation.
func PrintLine() {
	fmt.Println(strings.Repeat("=", 50))
}

// PrintEmptyLine prints an empty line for better visual spacing in output.
func PrintEmptyLine() {
	fmt.Println()
}

// FormatColor wraps the given text with the specified color and adds reset.
// This allows custom colored text formatting without direct printing.
func FormatColor(text string, color string) string {
	return color + text + ResetColor
}

// FormatSuccess formats text as a success message without printing.
func FormatSuccess(text string) string {
	return FormatColor(text, SuccessColor)
}

// FormatError formats text as an error message without printing.
func FormatError(text string) string {
	return FormatColor(text, ErrorColor)
}

// FormatWarning formats text as a warning message without printing.
func FormatWarning(text string) string {
	return FormatColor(text, WarningColor)
}

// FormatInfo formats text as an info message without printing.
func FormatInfo(text string) string {
	return FormatColor(text, InfoColor)
}

// FormatDebug formats text as a debug message without printing.
func FormatDebug(text string) string {
	return FormatColor(text, DebugColor)
}

// FormatHighlight formats text as highlighted without printing.
func FormatHighlight(text string) string {
	return FormatColor(text, HighlightColor)
}
