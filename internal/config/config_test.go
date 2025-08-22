package config

import (
	"os"
	"testing"
)

func TestLoad_RequiredFields(t *testing.T) {
	// Save original env vars and restore after test
	originalToken := os.Getenv("DISCORD_TOKEN")
	defer func() {
		if originalToken != "" {
			_ = os.Setenv("DISCORD_TOKEN", originalToken)
		} else {
			_ = os.Unsetenv("DISCORD_TOKEN")
		}
	}()

	t.Run("missing_discord_token", func(t *testing.T) {
		_ = os.Unsetenv("DISCORD_TOKEN")

		_, err := Load()
		if err == nil {
			t.Error("Expected error when DISCORD_TOKEN is missing, got nil")
		}

		expectedError := "missing required environment variables: DISCORD_TOKEN"
		if !containsError(err.Error(), expectedError) {
			t.Errorf("Expected error to contain %q, got %q", expectedError, err.Error())
		}
	})

	t.Run("valid_config", func(t *testing.T) {
		_ = os.Setenv("DISCORD_TOKEN", "test_token_123")

		cfg, err := Load()
		if err != nil {
			t.Fatalf("Expected no error with valid config, got %v", err)
		}

		if cfg.DiscordToken != "test_token_123" {
			t.Errorf("Expected DiscordToken to be 'test_token_123', got %q", cfg.DiscordToken)
		}

		if cfg.LogLevel != "info" {
			t.Errorf("Expected default LogLevel to be 'info', got %q", cfg.LogLevel)
		}
	})
}

func TestLoad_LogLevelValidation(t *testing.T) {
	// Save and restore env vars
	originalToken := os.Getenv("DISCORD_TOKEN")
	originalLogLevel := os.Getenv("LOG_LEVEL")
	defer func() {
		if originalToken != "" {
			_ = os.Setenv("DISCORD_TOKEN", originalToken)
		} else {
			_ = os.Unsetenv("DISCORD_TOKEN")
		}
		if originalLogLevel != "" {
			_ = os.Setenv("LOG_LEVEL", originalLogLevel)
		} else {
			_ = os.Unsetenv("LOG_LEVEL")
		}
	}()

	os.Setenv("DISCORD_TOKEN", "test_token")

	testCases := []struct {
		name        string
		logLevel    string
		expectError bool
	}{
		{"valid_debug", "debug", false},
		{"valid_info", "info", false},
		{"valid_warn", "warn", false},
		{"valid_error", "error", false},
		{"valid_uppercase", "INFO", false},
		{"invalid_level", "invalid", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_ = os.Setenv("LOG_LEVEL", tc.logLevel)

			_, err := Load()
			if tc.expectError && err == nil {
				t.Errorf("Expected error for log level %q, got nil", tc.logLevel)
			}
			if !tc.expectError && err != nil {
				t.Errorf("Expected no error for log level %q, got %v", tc.logLevel, err)
			}
		})
	}
}

func TestLoad_OptionalFields(t *testing.T) {
	// Save and restore env vars
	originalToken := os.Getenv("DISCORD_TOKEN")
	defer func() {
		if originalToken != "" {
			_ = os.Setenv("DISCORD_TOKEN", originalToken)
		} else {
			_ = os.Unsetenv("DISCORD_TOKEN")
		}
	}()

	_ = os.Setenv("DISCORD_TOKEN", "test_token")

	_, err := Load()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

// Helper function to check if error message contains expected text
func containsError(errorMsg, expected string) bool {
	return len(errorMsg) >= len(expected) &&
		errorMsg[len(errorMsg)-len(expected):] == expected ||
		errorMsg[:len(expected)] == expected ||
		(len(errorMsg) > len(expected) &&
			errorMsg[len(errorMsg)-len(expected)-1:len(errorMsg)-1] == expected)
}
