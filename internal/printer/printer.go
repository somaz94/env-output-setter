package printer

import (
	"fmt"
	"strings"
)

const (
	infoColor    = "\033[1;34m"
	successColor = "\033[1;32m"
	errorColor   = "\033[1;31m"
	resetColor   = "\033[0m"
)

func PrintSection(title string) {
	fmt.Printf("\n%s==================================================\n", infoColor)
	fmt.Printf("🚀 %s%s\n", title, resetColor)
}

func PrintSuccess(varType, key, value string) {
	fmt.Printf("%s  • %s: %s = %s%s\n", successColor, varType, key, value, resetColor)
}

func PrintError(message string) {
	fmt.Printf("%s❌ %s%s\n", errorColor, message, resetColor)
}

func PrintInfo(message string) {
	fmt.Printf("%s%s%s\n", infoColor, message, resetColor)
}

func PrintDebugSection(title string) {
	fmt.Printf("\n%s🔍 Debug Information (%s)\n", infoColor, title)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n%s", resetColor)
}

func PrintDebugInfo(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func PrintComplete() {
	fmt.Printf("\n%s==================================================\n", infoColor)
	fmt.Printf("✅ Execution Complete\n")
	fmt.Printf("Mode: GitHub Actions%s\n", resetColor)
}

func PrintLine() {
	fmt.Println(strings.Repeat("=", 50))
}
