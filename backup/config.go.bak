package config

import "os"

// Config holds the application configuration
type Config struct {
	EnvKeys      string
	EnvValues    string
	OutputKeys   string
	OutputValues string
	GithubEnv    string
	GithubOutput string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		EnvKeys:      os.Getenv("INPUT_ENV_KEY"),
		EnvValues:    os.Getenv("INPUT_ENV_VALUE"),
		OutputKeys:   os.Getenv("INPUT_OUTPUT_KEY"),
		OutputValues: os.Getenv("INPUT_OUTPUT_VALUE"),
		GithubEnv:    os.Getenv("GITHUB_ENV"),
		GithubOutput: os.Getenv("GITHUB_OUTPUT"),
	}
}
