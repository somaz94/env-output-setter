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
	if outputFile := os.Getenv("GITHUB_OUTPUT"); outputFile != "" {
		// 새로운 GITHUB_OUTPUT 파일 문법 사용
		outputs := []string{
			fmt.Sprintf("set_env_count=%d", envCount),
			fmt.Sprintf("set_output_count=%d", outputCount),
			fmt.Sprintf("status=%s", status),
			fmt.Sprintf("error_message=%s", errorMsg),
		}

		// 각 출력을 파일에 쓰기
		for _, output := range outputs {
			if err := appendToFile(outputFile, output); err != nil {
				fmt.Printf("Error writing to GITHUB_OUTPUT: %v\n", err)
			}
		}
	}

	if status == "failure" {
		os.Exit(1)
	}
	os.Exit(0)
}

// 파일에 출력을 추가하는 헬퍼 함수
func appendToFile(filename, content string) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(content + "\n"); err != nil {
		return err
	}
	return nil
}
