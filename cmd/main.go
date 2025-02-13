package main

import (
	"fmt"
	"os"

	"github.com/somaz94/env-output-setter/internal/config"
	"github.com/somaz94/env-output-setter/internal/printer"
	"github.com/somaz94/env-output-setter/internal/writer"
)

func main() {
	cfg := config.Load()

	printer.PrintSection("🚀 GitHub Environment and Output Setter")

	// Set environment variables
	if err := writer.SetEnv(cfg.EnvKeys, cfg.EnvValues); err != nil {
		fmt.Printf("❌ Error setting environment variables: %v\n", err)
		os.Exit(1)
	}

	// Set output variables
	if err := writer.SetOutput(cfg.OutputKeys, cfg.OutputValues); err != nil {
		fmt.Printf("❌ Error setting output variables: %v\n", err)
		os.Exit(1)
	}

	// Print final status
	printer.PrintSection("✅ Execution Complete")
	if cfg.GithubEnv == "" && cfg.GithubOutput == "" {
		fmt.Println("Mode: Local Execution (Simulation)")
	} else {
		fmt.Println("Mode: GitHub Actions")
	}
	printer.PrintLine()
}
