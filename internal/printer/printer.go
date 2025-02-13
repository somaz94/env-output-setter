package printer

import (
	"fmt"
	"strings"
)

func PrintSection(title string) {
	fmt.Printf("\n%s\n%s\n", strings.Repeat("=", 50), title)
}

func PrintSuccess(varType, key, value string) {
	fmt.Printf("  • %s: %s = %s\n", varType, key, value)
}

func PrintLine() {
	fmt.Println(strings.Repeat("=", 50))
}
