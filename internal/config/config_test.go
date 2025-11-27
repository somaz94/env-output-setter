package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected *Config
	}{
		{
			name: "Default values when no environment variables set",
			envVars: map[string]string{
				EnvKeyInput:    "",
				EnvValueInput:  "",
				DelimiterInput: "",
			},
			expected: &Config{
				EnvKeys:          "",
				EnvValues:        "",
				OutputKeys:       "",
				OutputValues:     "",
				GithubEnv:        "",
				GithubOutput:     "",
				Delimiter:        DefaultDelimiter,
				FailOnEmpty:      DefaultFailOnEmpty,
				TrimWhitespace:   DefaultTrimWhitespace,
				CaseSensitive:    DefaultCaseSensitive,
				ErrorOnDuplicate: DefaultErrorOnDuplicate,
				AllowEmpty:       DefaultAllowEmpty,
				ToUpper:          DefaultToUpper,
				ToLower:          DefaultToLower,
				EncodeURL:        DefaultEncodeURL,
				EscapeNewlines:   DefaultEscapeNewlines,
				MaxLength:        DefaultMaxLength,
				MaskSecrets:      DefaultMaskSecrets,
				MaskPattern:      DefaultMaskPattern,
				DebugMode:        DefaultDebugMode,
				GroupPrefix:      DefaultGroupPrefix,
				JsonSupport:      DefaultJsonSupport,
				ExportAsEnv:      DefaultExportAsEnv,
			},
		},
		{
			name: "Custom delimiter and basic inputs",
			envVars: map[string]string{
				EnvKeyInput:    "KEY1,KEY2",
				EnvValueInput:  "VALUE1,VALUE2",
				DelimiterInput: ",",
			},
			expected: &Config{
				EnvKeys:          "KEY1,KEY2",
				EnvValues:        "VALUE1,VALUE2",
				OutputKeys:       "",
				OutputValues:     "",
				GithubEnv:        "",
				GithubOutput:     "",
				Delimiter:        ",",
				FailOnEmpty:      DefaultFailOnEmpty,
				TrimWhitespace:   DefaultTrimWhitespace,
				CaseSensitive:    DefaultCaseSensitive,
				ErrorOnDuplicate: DefaultErrorOnDuplicate,
				AllowEmpty:       DefaultAllowEmpty,
				ToUpper:          DefaultToUpper,
				ToLower:          DefaultToLower,
				EncodeURL:        DefaultEncodeURL,
				EscapeNewlines:   DefaultEscapeNewlines,
				MaxLength:        DefaultMaxLength,
				MaskSecrets:      DefaultMaskSecrets,
				MaskPattern:      DefaultMaskPattern,
				DebugMode:        DefaultDebugMode,
				GroupPrefix:      DefaultGroupPrefix,
				JsonSupport:      DefaultJsonSupport,
				ExportAsEnv:      DefaultExportAsEnv,
			},
		},
		{
			name: "Boolean flags enabled",
			envVars: map[string]string{
				FailOnEmptyInput:      "false",
				TrimWhitespaceInput:   "false",
				CaseSensitiveInput:    "false",
				ErrorOnDuplicateInput: "false",
				MaskSecretsInput:      "true",
				ToUpperInput:          "true",
				DebugModeInput:        "true",
				JsonSupportInput:      "true",
				ExportAsEnvInput:      "true",
			},
			expected: &Config{
				EnvKeys:          "",
				EnvValues:        "",
				OutputKeys:       "",
				OutputValues:     "",
				GithubEnv:        "",
				GithubOutput:     "",
				Delimiter:        DefaultDelimiter,
				FailOnEmpty:      false,
				TrimWhitespace:   false,
				CaseSensitive:    false,
				ErrorOnDuplicate: false,
				AllowEmpty:       DefaultAllowEmpty,
				ToUpper:          true,
				ToLower:          DefaultToLower,
				EncodeURL:        DefaultEncodeURL,
				EscapeNewlines:   DefaultEscapeNewlines,
				MaxLength:        DefaultMaxLength,
				MaskSecrets:      true,
				MaskPattern:      DefaultMaskPattern,
				DebugMode:        true,
				GroupPrefix:      DefaultGroupPrefix,
				JsonSupport:      true,
				ExportAsEnv:      true,
			},
		},
		{
			name: "Advanced features enabled",
			envVars: map[string]string{
				GroupPrefixInput: "app",
				MaskPatternInput: "^secret_",
				MaxLengthInput:   "100",
				ToLowerInput:     "true",
				EncodeURLInput:   "true",
			},
			expected: &Config{
				EnvKeys:          "",
				EnvValues:        "",
				OutputKeys:       "",
				OutputValues:     "",
				GithubEnv:        "",
				GithubOutput:     "",
				Delimiter:        DefaultDelimiter,
				FailOnEmpty:      DefaultFailOnEmpty,
				TrimWhitespace:   DefaultTrimWhitespace,
				CaseSensitive:    DefaultCaseSensitive,
				ErrorOnDuplicate: DefaultErrorOnDuplicate,
				AllowEmpty:       DefaultAllowEmpty,
				ToUpper:          DefaultToUpper,
				ToLower:          true,
				EncodeURL:        true,
				EscapeNewlines:   DefaultEscapeNewlines,
				MaxLength:        100,
				MaskSecrets:      DefaultMaskSecrets,
				MaskPattern:      "^secret_",
				DebugMode:        DefaultDebugMode,
				GroupPrefix:      "app",
				JsonSupport:      DefaultJsonSupport,
				ExportAsEnv:      DefaultExportAsEnv,
			},
		},
		{
			name: "GitHub environment variables set",
			envVars: map[string]string{
				GithubEnvVar:     "/tmp/github_env",
				GithubOutputVar:  "/tmp/github_output",
				OutputKeyInput:   "STATUS",
				OutputValueInput: "SUCCESS",
			},
			expected: &Config{
				EnvKeys:          "",
				EnvValues:        "",
				OutputKeys:       "STATUS",
				OutputValues:     "SUCCESS",
				GithubEnv:        "/tmp/github_env",
				GithubOutput:     "/tmp/github_output",
				Delimiter:        DefaultDelimiter,
				FailOnEmpty:      DefaultFailOnEmpty,
				TrimWhitespace:   DefaultTrimWhitespace,
				CaseSensitive:    DefaultCaseSensitive,
				ErrorOnDuplicate: DefaultErrorOnDuplicate,
				AllowEmpty:       DefaultAllowEmpty,
				ToUpper:          DefaultToUpper,
				ToLower:          DefaultToLower,
				EncodeURL:        DefaultEncodeURL,
				EscapeNewlines:   DefaultEscapeNewlines,
				MaxLength:        DefaultMaxLength,
				MaskSecrets:      DefaultMaskSecrets,
				MaskPattern:      DefaultMaskPattern,
				DebugMode:        DefaultDebugMode,
				GroupPrefix:      DefaultGroupPrefix,
				JsonSupport:      DefaultJsonSupport,
				ExportAsEnv:      DefaultExportAsEnv,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean environment
			clearEnv()

			// Set test environment variables
			for key, value := range tt.envVars {
				if value != "" {
					os.Setenv(key, value)
				}
			}

			// Load configuration
			cfg := Load()

			// Verify all fields
			if cfg.EnvKeys != tt.expected.EnvKeys {
				t.Errorf("EnvKeys = %v, want %v", cfg.EnvKeys, tt.expected.EnvKeys)
			}
			if cfg.EnvValues != tt.expected.EnvValues {
				t.Errorf("EnvValues = %v, want %v", cfg.EnvValues, tt.expected.EnvValues)
			}
			if cfg.OutputKeys != tt.expected.OutputKeys {
				t.Errorf("OutputKeys = %v, want %v", cfg.OutputKeys, tt.expected.OutputKeys)
			}
			if cfg.OutputValues != tt.expected.OutputValues {
				t.Errorf("OutputValues = %v, want %v", cfg.OutputValues, tt.expected.OutputValues)
			}
			if cfg.GithubEnv != tt.expected.GithubEnv {
				t.Errorf("GithubEnv = %v, want %v", cfg.GithubEnv, tt.expected.GithubEnv)
			}
			if cfg.GithubOutput != tt.expected.GithubOutput {
				t.Errorf("GithubOutput = %v, want %v", cfg.GithubOutput, tt.expected.GithubOutput)
			}
			if cfg.Delimiter != tt.expected.Delimiter {
				t.Errorf("Delimiter = %v, want %v", cfg.Delimiter, tt.expected.Delimiter)
			}
			if cfg.FailOnEmpty != tt.expected.FailOnEmpty {
				t.Errorf("FailOnEmpty = %v, want %v", cfg.FailOnEmpty, tt.expected.FailOnEmpty)
			}
			if cfg.TrimWhitespace != tt.expected.TrimWhitespace {
				t.Errorf("TrimWhitespace = %v, want %v", cfg.TrimWhitespace, tt.expected.TrimWhitespace)
			}
			if cfg.CaseSensitive != tt.expected.CaseSensitive {
				t.Errorf("CaseSensitive = %v, want %v", cfg.CaseSensitive, tt.expected.CaseSensitive)
			}
			if cfg.ErrorOnDuplicate != tt.expected.ErrorOnDuplicate {
				t.Errorf("ErrorOnDuplicate = %v, want %v", cfg.ErrorOnDuplicate, tt.expected.ErrorOnDuplicate)
			}
			if cfg.AllowEmpty != tt.expected.AllowEmpty {
				t.Errorf("AllowEmpty = %v, want %v", cfg.AllowEmpty, tt.expected.AllowEmpty)
			}
			if cfg.ToUpper != tt.expected.ToUpper {
				t.Errorf("ToUpper = %v, want %v", cfg.ToUpper, tt.expected.ToUpper)
			}
			if cfg.ToLower != tt.expected.ToLower {
				t.Errorf("ToLower = %v, want %v", cfg.ToLower, tt.expected.ToLower)
			}
			if cfg.EncodeURL != tt.expected.EncodeURL {
				t.Errorf("EncodeURL = %v, want %v", cfg.EncodeURL, tt.expected.EncodeURL)
			}
			if cfg.EscapeNewlines != tt.expected.EscapeNewlines {
				t.Errorf("EscapeNewlines = %v, want %v", cfg.EscapeNewlines, tt.expected.EscapeNewlines)
			}
			if cfg.MaxLength != tt.expected.MaxLength {
				t.Errorf("MaxLength = %v, want %v", cfg.MaxLength, tt.expected.MaxLength)
			}
			if cfg.MaskSecrets != tt.expected.MaskSecrets {
				t.Errorf("MaskSecrets = %v, want %v", cfg.MaskSecrets, tt.expected.MaskSecrets)
			}
			if cfg.MaskPattern != tt.expected.MaskPattern {
				t.Errorf("MaskPattern = %v, want %v", cfg.MaskPattern, tt.expected.MaskPattern)
			}
			if cfg.DebugMode != tt.expected.DebugMode {
				t.Errorf("DebugMode = %v, want %v", cfg.DebugMode, tt.expected.DebugMode)
			}
			if cfg.GroupPrefix != tt.expected.GroupPrefix {
				t.Errorf("GroupPrefix = %v, want %v", cfg.GroupPrefix, tt.expected.GroupPrefix)
			}
			if cfg.JsonSupport != tt.expected.JsonSupport {
				t.Errorf("JsonSupport = %v, want %v", cfg.JsonSupport, tt.expected.JsonSupport)
			}
			if cfg.ExportAsEnv != tt.expected.ExportAsEnv {
				t.Errorf("ExportAsEnv = %v, want %v", cfg.ExportAsEnv, tt.expected.ExportAsEnv)
			}

			// Clean up
			clearEnv()
		})
	}
}

func TestGetEnvWithDefault(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		expected     string
	}{
		{
			name:         "Environment variable not set",
			key:          "TEST_VAR",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
		{
			name:         "Environment variable set",
			key:          "TEST_VAR",
			defaultValue: "default",
			envValue:     "custom",
			expected:     "custom",
		},
		{
			name:         "Empty environment variable",
			key:          "TEST_VAR",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Unsetenv(tt.key)
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
			}

			result := getEnvWithDefault(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("getEnvWithDefault() = %v, want %v", result, tt.expected)
			}

			os.Unsetenv(tt.key)
		})
	}
}

func TestGetBoolEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue bool
		envValue     string
		expected     bool
	}{
		{
			name:         "Environment variable not set",
			key:          "TEST_BOOL",
			defaultValue: true,
			envValue:     "",
			expected:     true,
		},
		{
			name:         "Environment variable set to true",
			key:          "TEST_BOOL",
			defaultValue: false,
			envValue:     "true",
			expected:     true,
		},
		{
			name:         "Environment variable set to false",
			key:          "TEST_BOOL",
			defaultValue: true,
			envValue:     "false",
			expected:     false,
		},
		{
			name:         "Environment variable set to 1",
			key:          "TEST_BOOL",
			defaultValue: false,
			envValue:     "1",
			expected:     true,
		},
		{
			name:         "Environment variable set to 0",
			key:          "TEST_BOOL",
			defaultValue: true,
			envValue:     "0",
			expected:     false,
		},
		{
			name:         "Invalid boolean value",
			key:          "TEST_BOOL",
			defaultValue: true,
			envValue:     "invalid",
			expected:     true,
		},
		{
			name:         "Case insensitive TRUE",
			key:          "TEST_BOOL",
			defaultValue: false,
			envValue:     "TRUE",
			expected:     true,
		},
		{
			name:         "Case insensitive False",
			key:          "TEST_BOOL",
			defaultValue: true,
			envValue:     "False",
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Unsetenv(tt.key)
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
			}

			result := getBoolEnv(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("getBoolEnv() = %v, want %v", result, tt.expected)
			}

			os.Unsetenv(tt.key)
		})
	}
}

func TestGetIntEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue int
		envValue     string
		expected     int
	}{
		{
			name:         "Environment variable not set",
			key:          "TEST_INT",
			defaultValue: 42,
			envValue:     "",
			expected:     42,
		},
		{
			name:         "Environment variable set to positive integer",
			key:          "TEST_INT",
			defaultValue: 0,
			envValue:     "100",
			expected:     100,
		},
		{
			name:         "Environment variable set to zero",
			key:          "TEST_INT",
			defaultValue: 42,
			envValue:     "0",
			expected:     0,
		},
		{
			name:         "Environment variable set to negative integer",
			key:          "TEST_INT",
			defaultValue: 0,
			envValue:     "-50",
			expected:     -50,
		},
		{
			name:         "Invalid integer value",
			key:          "TEST_INT",
			defaultValue: 42,
			envValue:     "not-a-number",
			expected:     42,
		},
		{
			name:         "Float value (invalid for int)",
			key:          "TEST_INT",
			defaultValue: 42,
			envValue:     "3.14",
			expected:     42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Unsetenv(tt.key)
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
			}

			result := getIntEnv(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("getIntEnv() = %v, want %v", result, tt.expected)
			}

			os.Unsetenv(tt.key)
		})
	}
}

// clearEnv removes all test environment variables
func clearEnv() {
	envVars := []string{
		EnvKeyInput,
		EnvValueInput,
		OutputKeyInput,
		OutputValueInput,
		DelimiterInput,
		FailOnEmptyInput,
		TrimWhitespaceInput,
		CaseSensitiveInput,
		ErrorOnDuplicateInput,
		MaskSecretsInput,
		MaskPatternInput,
		ToUpperInput,
		ToLowerInput,
		EncodeURLInput,
		EscapeNewlinesInput,
		MaxLengthInput,
		AllowEmptyInput,
		DebugModeInput,
		GroupPrefixInput,
		JsonSupportInput,
		ExportAsEnvInput,
		GithubEnvVar,
		GithubOutputVar,
	}

	for _, envVar := range envVars {
		os.Unsetenv(envVar)
	}
}
