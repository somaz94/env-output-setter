package config

import (
	"os"
	"strconv"
	"strings"
)

// Config holds the application configuration
type Config struct {
	EnvKeys          string
	EnvValues        string
	OutputKeys       string
	OutputValues     string
	GithubEnv        string
	GithubOutput     string
	Delimiter        string
	FailOnEmpty      bool
	TrimWhitespace   bool
	CaseSensitive    bool
	ErrorOnDuplicate bool
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		EnvKeys:          os.Getenv("INPUT_ENV_KEY"),
		EnvValues:        os.Getenv("INPUT_ENV_VALUE"),
		OutputKeys:       os.Getenv("INPUT_OUTPUT_KEY"),
		OutputValues:     os.Getenv("INPUT_OUTPUT_VALUE"),
		GithubEnv:        os.Getenv("GITHUB_ENV"),
		GithubOutput:     os.Getenv("GITHUB_OUTPUT"),
		Delimiter:        getEnvWithDefault("INPUT_DELIMITER", ","),
		FailOnEmpty:      getBoolEnv("INPUT_FAIL_ON_EMPTY", true),
		TrimWhitespace:   getBoolEnv("INPUT_TRIM_WHITESPACE", true),
		CaseSensitive:    getBoolEnv("INPUT_CASE_SENSITIVE", true),
		ErrorOnDuplicate: getBoolEnv("INPUT_ERROR_ON_DUPLICATE", true),
	}
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	b, err := strconv.ParseBool(strings.ToLower(value))
	if err != nil {
		return defaultValue
	}
	return b
}
