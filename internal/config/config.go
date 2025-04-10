package config

import (
	"os"
	"strconv"
	"strings"
)

// Environment Variable Names Constants
const (
	// Input Environment Variables
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

	// GitHub Environment Variables
	GithubEnvVar    = "GITHUB_ENV"
	GithubOutputVar = "GITHUB_OUTPUT"
)

// Default Values Constants
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

// Config holds the application configuration
type Config struct {
	// Input/Output Keys and Values
	EnvKeys      string
	EnvValues    string
	OutputKeys   string
	OutputValues string

	// GitHub File Paths
	GithubEnv    string
	GithubOutput string

	// Input Processing Options
	Delimiter        string
	FailOnEmpty      bool
	TrimWhitespace   bool
	CaseSensitive    bool
	ErrorOnDuplicate bool
	AllowEmpty       bool

	// Value Transformation Options
	ToUpper        bool
	ToLower        bool
	EncodeURL      bool
	EscapeNewlines bool
	MaxLength      int

	// Security Options
	MaskSecrets bool
	MaskPattern string

	// Debug Options
	DebugMode bool

	// New configuration fields
	GroupPrefix string
	JsonSupport bool
	ExportAsEnv bool
}

// Load loads configuration from environment variables
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

		// New configuration fields
		GroupPrefix: getEnvWithDefault(GroupPrefixInput, DefaultGroupPrefix),
		JsonSupport: getBoolEnv(JsonSupportInput, DefaultJsonSupport),
		ExportAsEnv: getBoolEnv(ExportAsEnvInput, DefaultExportAsEnv),
	}
}

// getEnvWithDefault returns the value of the environment variable or the default value if not set
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getBoolEnv parses a boolean environment variable
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

// getIntEnv parses an integer environment variable
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
