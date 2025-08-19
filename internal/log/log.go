package log

import (
	"log/slog"
	"os"
	"strings"
)

// Logger wraps slog.Logger with additional context methods
type Logger struct {
	*slog.Logger
}

// InitLogger initializes the global logger with the specified level
func InitLogger(level string) *Logger {
	var logLevel slog.Level

	switch strings.ToLower(level) {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	// Create structured logger with JSON output
	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)

	// Set as default logger
	slog.SetDefault(logger)

	return &Logger{Logger: logger}
}

// WithGuild returns a logger with guild_id context
func (l *Logger) WithGuild(guildID string) *Logger {
	return &Logger{Logger: l.Logger.With("guild_id", guildID)}
}

// WithUser returns a logger with user_id context
func (l *Logger) WithUser(userID string) *Logger {
	return &Logger{Logger: l.Logger.With("user_id", userID)}
}

// WithCommand returns a logger with command context
func (l *Logger) WithCommand(command string) *Logger {
	return &Logger{Logger: l.Logger.With("command", command)}
}

// WithChannel returns a logger with channel_id context
func (l *Logger) WithChannel(channelID string) *Logger {
	return &Logger{Logger: l.Logger.With("channel_id", channelID)}
}

// WithDiscordContext returns a logger with common Discord context fields
func (l *Logger) WithDiscordContext(guildID, channelID, userID string) *Logger {
	return &Logger{
		Logger: l.Logger.With(
			"guild_id", guildID,
			"channel_id", channelID,
			"user_id", userID,
		),
	}
}

// WithError returns a logger with error context
func (l *Logger) WithError(err error) *Logger {
	return &Logger{Logger: l.Logger.With("error", err.Error())}
}

// Global logger instance
var globalLogger *Logger

// SetGlobalLogger sets the global logger instance
func SetGlobalLogger(logger *Logger) {
	globalLogger = logger
}

// Global returns the global logger instance
func Global() *Logger {
	if globalLogger == nil {
		globalLogger = InitLogger("info")
	}
	return globalLogger
}