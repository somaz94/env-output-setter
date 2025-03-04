package config

import (
	"os"
	"strconv"
	"strings"
)

// 환경 변수 이름 상수
const (
	// 입력 환경 변수
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

	// GitHub 환경 변수
	GithubEnvVar    = "GITHUB_ENV"
	GithubOutputVar = "GITHUB_OUTPUT"
)

// 기본값 상수
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
)

// Config holds the application configuration
type Config struct {
	// 입력/출력 키와 값
	EnvKeys      string
	EnvValues    string
	OutputKeys   string
	OutputValues string

	// GitHub 파일 경로
	GithubEnv    string
	GithubOutput string

	// 입력 처리 옵션
	Delimiter        string
	FailOnEmpty      bool
	TrimWhitespace   bool
	CaseSensitive    bool
	ErrorOnDuplicate bool
	AllowEmpty       bool

	// 값 변환 옵션
	ToUpper        bool
	ToLower        bool
	EncodeURL      bool
	EscapeNewlines bool
	MaxLength      int

	// 보안 옵션
	MaskSecrets bool
	MaskPattern string

	// 디버그 옵션
	DebugMode bool
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		// 입력/출력 키와 값
		EnvKeys:      os.Getenv(EnvKeyInput),
		EnvValues:    os.Getenv(EnvValueInput),
		OutputKeys:   os.Getenv(OutputKeyInput),
		OutputValues: os.Getenv(OutputValueInput),

		// GitHub 파일 경로
		GithubEnv:    os.Getenv(GithubEnvVar),
		GithubOutput: os.Getenv(GithubOutputVar),

		// 입력 처리 옵션
		Delimiter:        getEnvWithDefault(DelimiterInput, DefaultDelimiter),
		FailOnEmpty:      getBoolEnv(FailOnEmptyInput, DefaultFailOnEmpty),
		TrimWhitespace:   getBoolEnv(TrimWhitespaceInput, DefaultTrimWhitespace),
		CaseSensitive:    getBoolEnv(CaseSensitiveInput, DefaultCaseSensitive),
		ErrorOnDuplicate: getBoolEnv(ErrorOnDuplicateInput, DefaultErrorOnDuplicate),
		AllowEmpty:       getBoolEnv(AllowEmptyInput, DefaultAllowEmpty),

		// 값 변환 옵션
		ToUpper:        getBoolEnv(ToUpperInput, DefaultToUpper),
		ToLower:        getBoolEnv(ToLowerInput, DefaultToLower),
		EncodeURL:      getBoolEnv(EncodeURLInput, DefaultEncodeURL),
		EscapeNewlines: getBoolEnv(EscapeNewlinesInput, DefaultEscapeNewlines),
		MaxLength:      getIntEnv(MaxLengthInput, DefaultMaxLength),

		// 보안 옵션
		MaskSecrets: getBoolEnv(MaskSecretsInput, DefaultMaskSecrets),
		MaskPattern: getEnvWithDefault(MaskPatternInput, DefaultMaskPattern),

		// 디버그 옵션
		DebugMode: getBoolEnv(DebugModeInput, DefaultDebugMode),
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
