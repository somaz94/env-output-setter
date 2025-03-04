package printer

import (
	"fmt"
	"strings"
)

// Color Constant
const (
	// Basic Colors
	InfoColor    = "\033[1;34m" // Blue
	SuccessColor = "\033[1;32m" // Green
	ErrorColor   = "\033[1;31m" // Red
	WarningColor = "\033[1;33m" // Yellow
	DebugColor   = "\033[1;36m" // Cyan
	ResetColor   = "\033[0m"    // Reset Color

	// Additional Colors
	HeaderColor    = "\033[1;35m" // Purple
	HighlightColor = "\033[1;37m" // Light White
)

// Separator Constants
const (
	DoubleLine = "=================================================="
	SingleLine = "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
)

// PrintSection prints a section header
func PrintSection(title string) {
	fmt.Printf("\n%s%s\n", InfoColor, DoubleLine)
	fmt.Printf("%s%s%s\n", InfoColor, title, ResetColor)
}

// PrintSuccess prints a success message for a variable
func PrintSuccess(varType, key, value string) {
	fmt.Printf("%s  • %s: %s = %s%s\n", SuccessColor, varType, key, value, ResetColor)
}

// PrintError prints an error message
func PrintError(message string) {
	fmt.Printf("%s❌ %s%s\n", ErrorColor, message, ResetColor)
}

// PrintWarning prints a warning message
func PrintWarning(message string) {
	fmt.Printf("%s⚠️ %s%s\n", WarningColor, message, ResetColor)
}

// PrintInfo prints an informational message
func PrintInfo(message string) {
	fmt.Printf("%s%s%s\n", InfoColor, message, ResetColor)
}

// PrintDebugSection prints a debug section header
func PrintDebugSection(title string) {
	fmt.Printf("\n%s🔍 Debug Information (%s)\n", DebugColor, title)
	fmt.Printf("%s%s\n", SingleLine, ResetColor)
}

// PrintDebugInfo prints debug information
func PrintDebugInfo(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

// PrintDebugHighlight prints highlighted debug information
func PrintDebugHighlight(format string, args ...interface{}) {
	fmt.Printf("%s%s%s", HighlightColor, fmt.Sprintf(format, args...), ResetColor)
}

// PrintComplete prints a completion message
func PrintComplete() {
	fmt.Printf("\n%s%s\n", InfoColor, DoubleLine)
	fmt.Printf("%s✅ Execution Complete\n", InfoColor)
	fmt.Printf("Mode: GitHub Actions%s\n", ResetColor)
}

// PrintLine prints a horizontal line
func PrintLine() {
	fmt.Println(strings.Repeat("=", 50))
}

// PrintEmptyLine prints an empty line
func PrintEmptyLine() {
	fmt.Println()
}
