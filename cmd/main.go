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

	// Initialize counters
	var envCount, outputCount int
	var status = "success"
	var errorMsg string

	// Set environment variables
	envCount, err := writer.SetEnv(cfg)
	if err != nil {
		errorMsg = fmt.Sprintf("❌ Error setting environment variables: %v", err)
		printer.PrintError(errorMsg)
		setOutputAndExit(0, 0, "failure", errorMsg)
	}

	// Set output variables
	outputCount, err = writer.SetOutput(cfg)
	if err != nil {
		errorMsg = fmt.Sprintf("❌ Error setting output variables: %v", err)
		printer.PrintError(errorMsg)
		setOutputAndExit(envCount, outputCount, "failure", errorMsg)
	}

	// Print final status
	printer.PrintSection("✅ Execution Complete")
	if cfg.GithubEnv == "" && cfg.GithubOutput == "" {
		printer.PrintInfo("Mode: Local Execution (Simulation)")
	} else {
		printer.PrintInfo("Mode: GitHub Actions")
	}

	setOutputAndExit(envCount, outputCount, status, errorMsg)
}

func setOutputAndExit(envCount, outputCount int, status, errorMsg string) {
	if os.Getenv("GITHUB_OUTPUT") != "" {
		fmt.Printf("::set-output name=set_env_count::%d\n", envCount)
		fmt.Printf("::set-output name=set_output_count::%d\n", outputCount)
		fmt.Printf("::set-output name=status::%s\n", status)
		fmt.Printf("::set-output name=error_message::%s\n", errorMsg)
	}

	if status == "failure" {
		os.Exit(1)
	}
	os.Exit(0)
}
