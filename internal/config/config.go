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
	MaskSecrets      bool
	MaskPattern      string
	ToUpper          bool
	ToLower          bool
	EncodeURL        bool
	EscapeNewlines   bool
	MaxLength        int
	AllowEmpty       bool
	DebugMode        bool
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
		MaskSecrets:      getBoolEnv("INPUT_MASK_SECRETS", false),
		MaskPattern:      getEnvWithDefault("INPUT_MASK_PATTERN", ""),
		ToUpper:          getBoolEnv("INPUT_TO_UPPER", false),
		ToLower:          getBoolEnv("INPUT_TO_LOWER", false),
		EncodeURL:        getBoolEnv("INPUT_ENCODE_URL", false),
		EscapeNewlines:   getBoolEnv("INPUT_ESCAPE_NEWLINES", true),
		MaxLength:        getIntEnv("INPUT_MAX_LENGTH", 0),
		AllowEmpty:       getBoolEnv("INPUT_ALLOW_EMPTY", false),
		DebugMode:        getBoolEnv("DEBUG_MODE", false),
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

func getIntEnv(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return i
}
