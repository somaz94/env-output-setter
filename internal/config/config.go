package config

import (
	"os"
	"strconv"
	"strings"
)

// Input environment variable names
const (
	EnvKeyInput           = "INPUT_ENV_KEY"
	EnvValueInput         = "INPUT_ENV_VALUE"
	OutputKeyInput        = "INPUT_OUTPUT_KEY"
	OutputValueInput      = "INPUT_OUTPUT_VALUE"
	DelimiterInput        = "INPUT_DELIMITER"
	FailOnEmptyInput      = "INPUT_FAIL_ON_EMPTY"
	TrimWhitespaceInput   = "INPUT_TRIM_WHITESPACE"
	CaseSensitiveInput    = "INPUT_CASE_SENSITIVE"
	ErrorOnDuplicateInput = "INPUT_ERROR_ON_DUPLICATE"
	MaskSecretsInput      = "INPUT_MASK_SECRETS"
	MaskPatternInput      = "INPUT_MASK_PATTERN"
	ToUpperInput          = "INPUT_TO_UPPER"
	ToLowerInput          = "INPUT_TO_LOWER"
	EncodeURLInput        = "INPUT_ENCODE_URL"
	EscapeNewlinesInput   = "INPUT_ESCAPE_NEWLINES"
	MaxLengthInput        = "INPUT_MAX_LENGTH"
	AllowEmptyInput       = "INPUT_ALLOW_EMPTY"
	DebugModeInput        = "DEBUG_MODE"
	GroupPrefixInput      = "INPUT_GROUP_PREFIX"
	JsonSupportInput      = "INPUT_JSON_SUPPORT"
	ExportAsEnvInput      = "INPUT_EXPORT_AS_ENV"
)

// GitHub environment variables
const (
	GithubEnvVar    = "GITHUB_ENV"
	GithubOutputVar = "GITHUB_OUTPUT"
)

// Default values for configuration parameters
const (
	DefaultDelimiter        = ","
	DefaultFailOnEmpty      = true
	DefaultTrimWhitespace   = true
	DefaultCaseSensitive    = true
	DefaultErrorOnDuplicate = true
	DefaultMaskSecrets      = false
	DefaultMaskPattern      = ""
	DefaultToUpper          = false
	DefaultToLower          = false
	DefaultEncodeURL        = false
	DefaultEscapeNewlines   = true
	DefaultMaxLength        = 0
	DefaultAllowEmpty       = false
	DefaultDebugMode        = false
	DefaultGroupPrefix      = ""
	DefaultJsonSupport      = false
	DefaultExportAsEnv      = false
)

// Config holds the application configuration settings loaded from environment variables.
// It includes settings for input/output handling, value transformations, security options,
// and debugging.
type Config struct {
	// Input/Output Keys and Values
	EnvKeys      string // Environment variable keys to process
	EnvValues    string // Environment variable values to process
	OutputKeys   string // Output keys for GitHub Actions
	OutputValues string // Output values for GitHub Actions

	// GitHub File Paths
	GithubEnv    string // Path to GITHUB_ENV file
	GithubOutput string // Path to GITHUB_OUTPUT file

	// Input Processing Options
	Delimiter        string // Delimiter for splitting multiple keys/values
	FailOnEmpty      bool   // Whether to fail when encountering empty values
	TrimWhitespace   bool   // Whether to trim whitespace from values
	CaseSensitive    bool   // Whether key comparisons are case sensitive
	ErrorOnDuplicate bool   // Whether to error on duplicate keys
	AllowEmpty       bool   // Whether empty values are allowed in the output

	// Value Transformation Options
	ToUpper        bool // Convert values to uppercase
	ToLower        bool // Convert values to lowercase
	EncodeURL      bool // URL-encode values
	EscapeNewlines bool // Escape newlines in values
	MaxLength      int  // Maximum length for values (0 = no limit)

	// Security Options
	MaskSecrets bool   // Whether to mask secret values in logs
	MaskPattern string // Regex pattern for identifying values to mask

	// Debug Options
	DebugMode bool // Enable debug mode for verbose logging

	// Advanced Options
	GroupPrefix string // Prefix for grouping related outputs
	JsonSupport bool   // Support for JSON values
	ExportAsEnv bool   // Whether to export values as environment variables
}

// Load creates a new Config instance with values loaded from environment variables.
// Default values are used for any settings not specified in the environment.
func Load() *Config {
	return &Config{
		// Input/Output Keys and Values
		EnvKeys:      os.Getenv(EnvKeyInput),
		EnvValues:    os.Getenv(EnvValueInput),
		OutputKeys:   os.Getenv(OutputKeyInput),
		OutputValues: os.Getenv(OutputValueInput),

		// GitHub File Paths
		GithubEnv:    os.Getenv(GithubEnvVar),
		GithubOutput: os.Getenv(GithubOutputVar),

		// Input Processing Options
		Delimiter:        getEnvWithDefault(DelimiterInput, DefaultDelimiter),
		FailOnEmpty:      getBoolEnv(FailOnEmptyInput, DefaultFailOnEmpty),
		TrimWhitespace:   getBoolEnv(TrimWhitespaceInput, DefaultTrimWhitespace),
		CaseSensitive:    getBoolEnv(CaseSensitiveInput, DefaultCaseSensitive),
		ErrorOnDuplicate: getBoolEnv(ErrorOnDuplicateInput, DefaultErrorOnDuplicate),
		AllowEmpty:       getBoolEnv(AllowEmptyInput, DefaultAllowEmpty),

		// Value Transformation Options
		ToUpper:        getBoolEnv(ToUpperInput, DefaultToUpper),
		ToLower:        getBoolEnv(ToLowerInput, DefaultToLower),
		EncodeURL:      getBoolEnv(EncodeURLInput, DefaultEncodeURL),
		EscapeNewlines: getBoolEnv(EscapeNewlinesInput, DefaultEscapeNewlines),
		MaxLength:      getIntEnv(MaxLengthInput, DefaultMaxLength),

		// Security Options
		MaskSecrets: getBoolEnv(MaskSecretsInput, DefaultMaskSecrets),
		MaskPattern: getEnvWithDefault(MaskPatternInput, DefaultMaskPattern),

		// Debug Options
		DebugMode: getBoolEnv(DebugModeInput, DefaultDebugMode),

		// Advanced Options
		GroupPrefix: getEnvWithDefault(GroupPrefixInput, DefaultGroupPrefix),
		JsonSupport: getBoolEnv(JsonSupportInput, DefaultJsonSupport),
		ExportAsEnv: getBoolEnv(ExportAsEnvInput, DefaultExportAsEnv),
	}
}

// getEnvWithDefault retrieves an environment variable value or returns
// the specified default value if the variable is not set or empty.
func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getBoolEnv retrieves a boolean environment variable value.
// It parses the string value to a boolean, returning the default value
// if the variable is not set, empty, or cannot be parsed as a boolean.
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

// getIntEnv retrieves an integer environment variable value.
// It parses the string value to an integer, returning the default value
// if the variable is not set, empty, or cannot be parsed as an integer.
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
