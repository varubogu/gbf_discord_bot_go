package config

import (
	"fmt"
	"os"
	"strings"
)

// Config holds all configuration for the bot
type Config struct {
	// Discord settings (required)
	DiscordToken string

	// Logging settings (optional)
	LogLevel string
}

// Load reads configuration from environment variables and validates required fields
func Load() (*Config, error) {
	config := &Config{
		DiscordToken: os.Getenv("DISCORD_TOKEN"),
		LogLevel:     getEnvWithDefault("LOG_LEVEL", "info"),
	}

	// Validate required fields
	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return config, nil
}

// validate checks that all required configuration is present
func (c *Config) validate() error {
	var missingFields []string

	if c.DiscordToken == "" {
		missingFields = append(missingFields, "DISCORD_TOKEN")
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("missing required environment variables: %s", strings.Join(missingFields, ", "))
	}

	// Validate log level
	validLogLevels := []string{"debug", "info", "warn", "error"}
	logLevel := strings.ToLower(c.LogLevel)
	isValidLogLevel := false
	for _, level := range validLogLevels {
		if logLevel == level {
			isValidLogLevel = true
			break
		}
	}
	if !isValidLogLevel {
		return fmt.Errorf("invalid LOG_LEVEL: %s, must be one of: %s", c.LogLevel, strings.Join(validLogLevels, ", "))
	}

	return nil
}

// getEnvWithDefault returns the environment variable value or a default if not set
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
