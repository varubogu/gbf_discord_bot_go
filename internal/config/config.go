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

	// Database settings (required)
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string

	// Test environment settings (optional)
	TestDiscordToken string
	TestDBHost       string
	TestDBUser       string
	TestDBPassword   string
	TestDBDatabase   string
	TestDBPort       string
}

// Load reads configuration from environment variables and validates required fields
func Load() (*Config, error) {
	config := &Config{
		DiscordToken: os.Getenv("DISCORD_TOKEN"),
		LogLevel:     getEnvWithDefault("LOG_LEVEL", "info"),

		// Database settings
		DBHost:     getEnvWithDefault("DB_HOST", "localhost"),
		DBUser:     getEnvWithDefault("DB_USER", ""),
		DBPassword: getEnvWithDefault("DB_PASSWORD", ""),
		DBName:     getEnvWithDefault("DB_NAME", ""),
		DBPort:     getEnvWithDefault("DB_PORT", "5432"),

		// Test environment settings
		TestDiscordToken: os.Getenv("TEST_DISCORD_TOKEN"),
		TestDBHost:       getEnvWithDefault("TEST_DBHOST", "localhost"),
		TestDBUser:       getEnvWithDefault("TEST_DBUSER", ""),
		TestDBPassword:   getEnvWithDefault("TEST_DBPASSWORD", ""),
		TestDBDatabase:   getEnvWithDefault("TEST_DBDATABASE", ""),
		TestDBPort:       getEnvWithDefault("TEST_DB_PORT", "5432"),
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
