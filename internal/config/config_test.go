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

func TestLoad_DatabaseSettings(t *testing.T) {
	// Save and restore env vars
	originalToken := os.Getenv("DISCORD_TOKEN")
	originalDBHost := os.Getenv("DB_HOST")
	originalDBUser := os.Getenv("DB_USER")
	originalDBPassword := os.Getenv("DB_PASSWORD")
	originalDBName := os.Getenv("DB_NAME")
	originalDBPort := os.Getenv("DB_PORT")
	
	defer func() {
		if originalToken != "" {
			_ = os.Setenv("DISCORD_TOKEN", originalToken)
		} else {
			_ = os.Unsetenv("DISCORD_TOKEN")
		}
		restoreEnv("DB_HOST", originalDBHost)
		restoreEnv("DB_USER", originalDBUser)
		restoreEnv("DB_PASSWORD", originalDBPassword)
		restoreEnv("DB_NAME", originalDBName)
		restoreEnv("DB_PORT", originalDBPort)
	}()

	_ = os.Setenv("DISCORD_TOKEN", "test_token")

	t.Run("default_values", func(t *testing.T) {
		// Clear all DB env vars
		_ = os.Unsetenv("DB_HOST")
		_ = os.Unsetenv("DB_USER")
		_ = os.Unsetenv("DB_PASSWORD")
		_ = os.Unsetenv("DB_NAME")
		_ = os.Unsetenv("DB_PORT")

		cfg, err := Load()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if cfg.DBHost != "localhost" {
			t.Errorf("Expected default DBHost to be 'localhost', got %q", cfg.DBHost)
		}
		if cfg.DBPort != "5432" {
			t.Errorf("Expected default DBPort to be '5432', got %q", cfg.DBPort)
		}
		if cfg.DBUser != "" {
			t.Errorf("Expected default DBUser to be empty, got %q", cfg.DBUser)
		}
	})

	t.Run("custom_values", func(t *testing.T) {
		_ = os.Setenv("DB_HOST", "custom.host")
		_ = os.Setenv("DB_USER", "custom_user")
		_ = os.Setenv("DB_PASSWORD", "custom_pass")
		_ = os.Setenv("DB_NAME", "custom_db")
		_ = os.Setenv("DB_PORT", "3306")

		cfg, err := Load()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if cfg.DBHost != "custom.host" {
			t.Errorf("Expected DBHost to be 'custom.host', got %q", cfg.DBHost)
		}
		if cfg.DBUser != "custom_user" {
			t.Errorf("Expected DBUser to be 'custom_user', got %q", cfg.DBUser)
		}
		if cfg.DBPassword != "custom_pass" {
			t.Errorf("Expected DBPassword to be 'custom_pass', got %q", cfg.DBPassword)
		}
		if cfg.DBName != "custom_db" {
			t.Errorf("Expected DBName to be 'custom_db', got %q", cfg.DBName)
		}
		if cfg.DBPort != "3306" {
			t.Errorf("Expected DBPort to be '3306', got %q", cfg.DBPort)
		}
	})
}

func TestLoad_TestEnvironmentSettings(t *testing.T) {
	// Save and restore env vars
	originalToken := os.Getenv("DISCORD_TOKEN")
	originalTestToken := os.Getenv("TEST_DISCORD_TOKEN")
	originalTestDBHost := os.Getenv("TEST_DBHOST")
	originalTestDBUser := os.Getenv("TEST_DBUSER")
	originalTestDBPassword := os.Getenv("TEST_DBPASSWORD")
	originalTestDBDatabase := os.Getenv("TEST_DBDATABASE")
	
	defer func() {
		if originalToken != "" {
			_ = os.Setenv("DISCORD_TOKEN", originalToken)
		} else {
			_ = os.Unsetenv("DISCORD_TOKEN")
		}
		restoreEnv("TEST_DISCORD_TOKEN", originalTestToken)
		restoreEnv("TEST_DBHOST", originalTestDBHost)
		restoreEnv("TEST_DBUSER", originalTestDBUser)
		restoreEnv("TEST_DBPASSWORD", originalTestDBPassword)
		restoreEnv("TEST_DBDATABASE", originalTestDBDatabase)
	}()

	_ = os.Setenv("DISCORD_TOKEN", "test_token")

	t.Run("default_values", func(t *testing.T) {
		// Clear all test env vars
		_ = os.Unsetenv("TEST_DISCORD_TOKEN")
		_ = os.Unsetenv("TEST_DBHOST")
		_ = os.Unsetenv("TEST_DBUSER")
		_ = os.Unsetenv("TEST_DBPASSWORD")
		_ = os.Unsetenv("TEST_DBDATABASE")

		cfg, err := Load()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if cfg.TestDBHost != "localhost" {
			t.Errorf("Expected default TestDBHost to be 'localhost', got %q", cfg.TestDBHost)
		}
		if cfg.TestDiscordToken != "" {
			t.Errorf("Expected default TestDiscordToken to be empty, got %q", cfg.TestDiscordToken)
		}
	})

	t.Run("custom_values", func(t *testing.T) {
		_ = os.Setenv("TEST_DISCORD_TOKEN", "test_token_custom")
		_ = os.Setenv("TEST_DBHOST", "test.host")
		_ = os.Setenv("TEST_DBUSER", "test_user")
		_ = os.Setenv("TEST_DBPASSWORD", "test_pass")
		_ = os.Setenv("TEST_DBDATABASE", "test_db")

		cfg, err := Load()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if cfg.TestDiscordToken != "test_token_custom" {
			t.Errorf("Expected TestDiscordToken to be 'test_token_custom', got %q", cfg.TestDiscordToken)
		}
		if cfg.TestDBHost != "test.host" {
			t.Errorf("Expected TestDBHost to be 'test.host', got %q", cfg.TestDBHost)
		}
		if cfg.TestDBUser != "test_user" {
			t.Errorf("Expected TestDBUser to be 'test_user', got %q", cfg.TestDBUser)
		}
		if cfg.TestDBPassword != "test_pass" {
			t.Errorf("Expected TestDBPassword to be 'test_pass', got %q", cfg.TestDBPassword)
		}
		if cfg.TestDBDatabase != "test_db" {
			t.Errorf("Expected TestDBDatabase to be 'test_db', got %q", cfg.TestDBDatabase)
		}
	})
}

// Helper function to restore environment variable
func restoreEnv(key, originalValue string) {
	if originalValue != "" {
		_ = os.Setenv(key, originalValue)
	} else {
		_ = os.Unsetenv(key)
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
