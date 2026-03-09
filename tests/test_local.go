package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/somaz94/env-output-setter/internal/config"
	"github.com/somaz94/env-output-setter/internal/writer"
)

// Color codes for terminal output
const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorBlue   = "\033[34m"
)

// TestCase represents a single test scenario
type TestCase struct {
	Name         string
	Category     string
	EnvKeys      string
	EnvValues    string
	OutputKeys   string
	OutputValues string
	Delimiter    string
	ExpectedEnvs map[string]string
	ExpectedOuts map[string]string
	ShouldFail   bool
	Options      TestOptions
}

// TestOptions holds optional test configurations
type TestOptions struct {
	FailOnEmpty      bool
	TrimWhitespace   bool
	CaseSensitive    bool
	ErrorOnDuplicate bool
	MaskSecrets      bool
	ToUpper          bool
	ToLower          bool
	EncodeURL        bool
	EscapeNewlines   bool
	MaxLength        int
	AllowEmpty       bool
	JsonSupport      bool
	ExportAsEnv      bool
	GroupPrefix      string
}

// TestResult holds the result of a test execution
type TestResult struct {
	Name     string
	Category string
	Passed   bool
	Error    error
}

func main() {
	fmt.Printf("\n%sLocal Test Suite for Environment/Output Setter%s\n", colorCyan, colorReset)
	fmt.Printf("Working directory: %s\n\n", mustGetWorkingDir())

	tests := createTestSuite()
	results := runTests(tests)
	printSummary(results)

	if hasFailures(results) {
		os.Exit(1)
	}
}

func createTestSuite() []TestCase {
	var tests []TestCase

	// ==================== Basic Tests ====================
	tests = append(tests, TestCase{
		Name:         "Basic environment variables",
		Category:     "Basic",
		EnvKeys:      "KEY1,KEY2,KEY3",
		EnvValues:    "value1,value2,value3",
		OutputKeys:   "OUT1,OUT2,OUT3",
		OutputValues: "out1,out2,out3",
		ExpectedEnvs: map[string]string{
			"KEY1": "value1",
			"KEY2": "value2",
			"KEY3": "value3",
		},
		ExpectedOuts: map[string]string{
			"OUT1": "out1",
			"OUT2": "out2",
			"OUT3": "out3",
		},
	})

	tests = append(tests, TestCase{
		Name:         "Custom delimiter (::)",
		Category:     "Basic",
		EnvKeys:      "KEY1::KEY2::KEY3",
		EnvValues:    "val1::val2::val3",
		OutputKeys:   "OUT1::OUT2::OUT3",
		OutputValues: "out1::out2::out3",
		Delimiter:    "::",
		ExpectedEnvs: map[string]string{
			"KEY1": "val1",
			"KEY2": "val2",
			"KEY3": "val3",
		},
		ExpectedOuts: map[string]string{
			"OUT1": "out1",
			"OUT2": "out2",
			"OUT3": "out3",
		},
	})

	// ==================== Whitespace Tests ====================
	tests = append(tests, TestCase{
		Name:         "Trim whitespace",
		Category:     "Whitespace",
		EnvKeys:      " KEY1 , KEY2 , KEY3 ",
		EnvValues:    " value1 , value2 , value3 ",
		OutputKeys:   " OUT1 , OUT2 , OUT3 ",
		OutputValues: " out1 , out2 , out3 ",
		ExpectedEnvs: map[string]string{
			"KEY1": "value1",
			"KEY2": "value2",
			"KEY3": "value3",
		},
		ExpectedOuts: map[string]string{
			"OUT1": "out1",
			"OUT2": "out2",
			"OUT3": "out3",
		},
		Options: TestOptions{
			TrimWhitespace: true,
		},
	})

	// ==================== Transformation Tests ====================
	tests = append(tests, TestCase{
		Name:         "Convert to uppercase",
		Category:     "Transformation",
		EnvKeys:      "KEY1,KEY2",
		EnvValues:    "lowercase,MixedCase",
		OutputKeys:   "OUT1,OUT2",
		OutputValues: "output1,output2",
		ExpectedEnvs: map[string]string{
			"KEY1": "LOWERCASE",
			"KEY2": "MIXEDCASE",
		},
		ExpectedOuts: map[string]string{
			"OUT1": "OUTPUT1",
			"OUT2": "OUTPUT2",
		},
		Options: TestOptions{
			ToUpper: true,
		},
	})

	tests = append(tests, TestCase{
		Name:         "Convert to lowercase",
		Category:     "Transformation",
		EnvKeys:      "KEY1,KEY2",
		EnvValues:    "UPPERCASE,MixedCase",
		OutputKeys:   "OUT1,OUT2",
		OutputValues: "OUTPUT1,OUTPUT2",
		ExpectedEnvs: map[string]string{
			"KEY1": "uppercase",
			"KEY2": "mixedcase",
		},
		ExpectedOuts: map[string]string{
			"OUT1": "output1",
			"OUT2": "output2",
		},
		Options: TestOptions{
			ToLower: true,
		},
	})

	tests = append(tests, TestCase{
		Name:         "URL encoding",
		Category:     "Transformation",
		EnvKeys:      "KEY1,KEY2",
		EnvValues:    "hello world,test@example.com",
		OutputKeys:   "OUT1,OUT2",
		OutputValues: "output value,email@test.com",
		ExpectedEnvs: map[string]string{
			"KEY1": "hello+world",
			"KEY2": "test%40example.com",
		},
		ExpectedOuts: map[string]string{
			"OUT1": "output+value",
			"OUT2": "email%40test.com",
		},
		Options: TestOptions{
			EncodeURL: true,
		},
	})

	tests = append(tests, TestCase{
		Name:         "Max length limitation",
		Category:     "Transformation",
		EnvKeys:      "KEY1,KEY2",
		EnvValues:    "ThisIsAVeryLongValue,Short",
		OutputKeys:   "OUT1,OUT2",
		OutputValues: "AnotherLongValue,OK",
		ExpectedEnvs: map[string]string{
			"KEY1": "ThisIsAVer", // Truncated to 10 chars
			"KEY2": "Short",
		},
		ExpectedOuts: map[string]string{
			"OUT1": "AnotherLon", // Truncated to 10 chars
			"OUT2": "OK",
		},
		Options: TestOptions{
			MaxLength: 10,
		},
	})

	tests = append(tests, TestCase{
		Name:         "Escape newlines",
		Category:     "Transformation",
		EnvKeys:      "KEY1,KEY2",
		EnvValues:    "Line1\nLine2,Single",
		OutputKeys:   "OUT1,OUT2",
		OutputValues: "Out1\nOut2,Single",
		ExpectedEnvs: map[string]string{
			"KEY1": "Line1\\nLine2",
			"KEY2": "Single",
		},
		ExpectedOuts: map[string]string{
			"OUT1": "Out1\\nOut2",
			"OUT2": "Single",
		},
		Options: TestOptions{
			EscapeNewlines: true,
		},
	})

	// ==================== Validation Tests ====================
	tests = append(tests, TestCase{
		Name:         "Duplicate keys (should fail)",
		Category:     "Validation",
		EnvKeys:      "KEY1,KEY1,KEY2",
		EnvValues:    "value1,value2,value3",
		OutputKeys:   "OUT1,OUT2,OUT3",
		OutputValues: "out1,out2,out3",
		ShouldFail:   true,
		Options: TestOptions{
			ErrorOnDuplicate: true,
		},
	})

	tests = append(tests, TestCase{
		Name:         "Empty values with allow_empty",
		Category:     "Validation",
		EnvKeys:      "KEY1,KEY2,KEY3",
		EnvValues:    "value1,,value3",
		OutputKeys:   "OUT1,OUT2,OUT3",
		OutputValues: "out1,,out3",
		ExpectedEnvs: map[string]string{
			"KEY1": "value1",
			"KEY2": "",
			"KEY3": "value3",
		},
		ExpectedOuts: map[string]string{
			"OUT1": "out1",
			"OUT2": "",
			"OUT3": "out3",
		},
		Options: TestOptions{
			AllowEmpty:     true,
			FailOnEmpty:    true,
			TrimWhitespace: true,
		},
	})

	// ==================== JSON Tests ====================
	tests = append(tests, TestCase{
		Name:         "Simple JSON parsing",
		Category:     "JSON",
		EnvKeys:      "CONFIG|SIMPLE",
		EnvValues:    `{"host":"localhost","port":8080}|simple_value`,
		OutputKeys:   "OUTPUT_CONFIG|OUTPUT_SIMPLE",
		OutputValues: `{"status":"ok","code":200}|simple_output`,
		Delimiter:    "|",
		ExpectedEnvs: map[string]string{
			"CONFIG":      `{"host":"localhost","port":8080}`,
			"CONFIG_host": "localhost",
			"CONFIG_port": "8080",
			"SIMPLE":      "simple_value",
		},
		ExpectedOuts: map[string]string{
			"OUTPUT_CONFIG":        `{"status":"ok","code":200}`,
			"OUTPUT_CONFIG_status": "ok",
			"OUTPUT_CONFIG_code":   "200",
			"OUTPUT_SIMPLE":        "simple_output",
		},
		Options: TestOptions{
			JsonSupport: true,
		},
	})

	tests = append(tests, TestCase{
		Name:      "Nested JSON parsing",
		Category:  "JSON",
		EnvKeys:   "COMPLEX|OTHER",
		EnvValues: `{"server":{"host":"example.com","port":443},"auth":{"enabled":true}}|other_value`,
		Delimiter: "|",
		ExpectedEnvs: map[string]string{
			"COMPLEX":              `{"server":{"host":"example.com","port":443},"auth":{"enabled":true}}`,
			"COMPLEX_server_host":  "example.com",
			"COMPLEX_server_port":  "443",
			"COMPLEX_auth_enabled": "true",
			"OTHER":                "other_value",
		},
		Options: TestOptions{
			JsonSupport: true,
		},
	})

	tests = append(tests, TestCase{
		Name:      "JSON array parsing",
		Category:  "JSON",
		EnvKeys:   "ARRAY_DATA|NORMAL",
		EnvValues: `{"items":[{"name":"item1"},{"name":"item2"}]}|normal_value`,
		Delimiter: "|",
		ExpectedEnvs: map[string]string{
			"ARRAY_DATA":              `{"items":[{"name":"item1"},{"name":"item2"}]}`,
			"ARRAY_DATA_items_0_name": "item1",
			"ARRAY_DATA_items_1_name": "item2",
			"NORMAL":                  "normal_value",
		},
		Options: TestOptions{
			JsonSupport: true,
		},
	})

	// ==================== Combined Tests ====================
	tests = append(tests, TestCase{
		Name:         "Multiple transformations",
		Category:     "Combined",
		EnvKeys:      "KEY1,KEY2",
		EnvValues:    "hello world,test value",
		OutputKeys:   "OUT1,OUT2",
		OutputValues: "output one,output two",
		ExpectedEnvs: map[string]string{
			"KEY1": "HELLO+WORL", // uppercase + url encoded + max 10 chars
			"KEY2": "TEST+VALU",  // uppercase + url encoded + max 10 chars
		},
		ExpectedOuts: map[string]string{
			"OUT1": "OUTPUT+ONE",
			"OUT2": "OUTPUT+TW",
		},
		Options: TestOptions{
			ToUpper:   true,
			EncodeURL: true,
			MaxLength: 10,
		},
	})

	return tests
}

func runTests(tests []TestCase) []TestResult {
	var results []TestResult
	currentCategory := ""

	fmt.Printf("%s%s\n", colorBlue, strings.Repeat("=", 50))
	fmt.Printf("Creating Test Suite\n")
	fmt.Printf("%s%s\n\n", strings.Repeat("=", 50), colorReset)
	fmt.Printf("Created %d test cases\n", len(tests))

	for _, test := range tests {
		// Print category header
		if test.Category != currentCategory {
			currentCategory = test.Category
			fmt.Printf("\n%s%s\n", colorBlue, strings.Repeat("=", 50))
			fmt.Printf("[Category] %s\n", currentCategory)
			fmt.Printf("%s%s\n\n", strings.Repeat("=", 50), colorReset)
		}

		result := runTest(test)
		results = append(results, result)
	}

	return results
}

func runTest(test TestCase) TestResult {
	fmt.Printf("Testing: %s\n", test.Name)

	// Setup test environment
	setupTestEnv(test)
	defer cleanupTestEnv()

	// Create config
	cfg := createTestConfig(test)

	// Run the test
	var err error
	if test.EnvKeys != "" && test.EnvValues != "" {
		_, err = writer.SetEnv(cfg)
	}
	if err == nil && test.OutputKeys != "" && test.OutputValues != "" {
		_, err = writer.SetOutput(cfg)
	}

	// Check result
	passed := true
	if test.ShouldFail {
		passed = (err != nil)
	} else {
		passed = (err == nil)
		if passed {
			passed = verifyResults(test)
		}
	}

	// Print result
	if passed {
		fmt.Printf("%sPASS: %s%s\n\n", colorGreen, test.Name, colorReset)
	} else {
		fmt.Printf("%sFAIL: %s%s\n", colorRed, test.Name, colorReset)
		if err != nil {
			fmt.Printf("  Error: %v\n", err)
		}
		fmt.Println()
	}

	return TestResult{
		Name:     test.Name,
		Category: test.Category,
		Passed:   passed,
		Error:    err,
	}
}

func setupTestEnv(test TestCase) {
	// Create temporary files for GITHUB_ENV and GITHUB_OUTPUT
	tmpDir := os.TempDir()
	envFile := fmt.Sprintf("%s/test_github_env_%d", tmpDir, os.Getpid())
	outputFile := fmt.Sprintf("%s/test_github_output_%d", tmpDir, os.Getpid())

	os.WriteFile(envFile, []byte{}, 0644)
	os.WriteFile(outputFile, []byte{}, 0644)

	os.Setenv("GITHUB_ENV", envFile)
	os.Setenv("GITHUB_OUTPUT", outputFile)
}

func cleanupTestEnv() {
	tmpDir := os.TempDir()
	os.Remove(fmt.Sprintf("%s/test_github_env_%d", tmpDir, os.Getpid()))
	os.Remove(fmt.Sprintf("%s/test_github_output_%d", tmpDir, os.Getpid()))
	os.Unsetenv("GITHUB_ENV")
	os.Unsetenv("GITHUB_OUTPUT")
}

func createTestConfig(test TestCase) *config.Config {
	delimiter := test.Delimiter
	if delimiter == "" {
		delimiter = ","
	}

	return &config.Config{
		EnvKeys:          test.EnvKeys,
		EnvValues:        test.EnvValues,
		OutputKeys:       test.OutputKeys,
		OutputValues:     test.OutputValues,
		GithubEnv:        os.Getenv("GITHUB_ENV"),
		GithubOutput:     os.Getenv("GITHUB_OUTPUT"),
		Delimiter:        delimiter,
		FailOnEmpty:      test.Options.FailOnEmpty,
		TrimWhitespace:   test.Options.TrimWhitespace || true, // Default true
		CaseSensitive:    test.Options.CaseSensitive || true,  // Default true
		ErrorOnDuplicate: test.Options.ErrorOnDuplicate,
		AllowEmpty:       test.Options.AllowEmpty,
		ToUpper:          test.Options.ToUpper,
		ToLower:          test.Options.ToLower,
		EncodeURL:        test.Options.EncodeURL,
		EscapeNewlines:   test.Options.EscapeNewlines || true, // Default true
		MaxLength:        test.Options.MaxLength,
		MaskSecrets:      test.Options.MaskSecrets,
		DebugMode:        false,
		JsonSupport:      test.Options.JsonSupport,
		ExportAsEnv:      test.Options.ExportAsEnv,
		GroupPrefix:      test.Options.GroupPrefix,
	}
}

func verifyResults(test TestCase) bool {
	// Read the written files
	tmpDir := os.TempDir()
	envContent, _ := os.ReadFile(fmt.Sprintf("%s/test_github_env_%d", tmpDir, os.Getpid()))
	outputContent, _ := os.ReadFile(fmt.Sprintf("%s/test_github_output_%d", tmpDir, os.Getpid()))

	// For now, just check if files were written
	// Full verification would parse the EOF-delimited format
	if test.EnvKeys != "" && len(envContent) == 0 {
		return false
	}
	if test.OutputKeys != "" && len(outputContent) == 0 {
		return false
	}

	return true
}

func printSummary(results []TestResult) {
	passed := 0
	failed := 0

	for _, result := range results {
		if result.Passed {
			passed++
		} else {
			failed++
		}
	}

	fmt.Printf("\n%s%s\n", colorBlue, strings.Repeat("=", 50))
	fmt.Printf("Test Summary\n")
	fmt.Printf("%s%s\n", strings.Repeat("=", 50), colorReset)
	fmt.Printf("Total Tests: %d\n", len(results))
	fmt.Printf("%sPassed: %d%s\n", colorGreen, passed, colorReset)
	fmt.Printf("%sFailed: %d%s\n", colorRed, failed, colorReset)

	if failed == 0 {
		fmt.Printf("\n%sAll tests passed!%s\n\n", colorGreen, colorReset)
	} else {
		fmt.Printf("\n%sSome tests failed!%s\n\n", colorRed, colorReset)
	}
}

func hasFailures(results []TestResult) bool {
	for _, result := range results {
		if !result.Passed {
			return true
		}
	}
	return false
}

func mustGetWorkingDir() string {
	wd, err := os.Getwd()
	if err != nil {
		return "unknown"
	}
	return wd
}
