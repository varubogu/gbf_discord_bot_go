package main

import (
	"context"
	"os"

	"github.com/varubogu/gbf_discord_bot_go/internal/config"
	"github.com/varubogu/gbf_discord_bot_go/internal/discord"
	"github.com/varubogu/gbf_discord_bot_go/internal/log"
)

func main() {
	// Load configuration from environment variables
	cfg, err := config.Load()
	if err != nil {
		// Use basic logging since structured logger isn't initialized yet
		log.Global().Error("Failed to load configuration", "error", err.Error())
		os.Exit(1)
	}

	// Initialize structured logger
	logger := log.InitLogger(cfg.LogLevel)
	log.SetGlobalLogger(logger)

	logger.Info("Starting GBF Discord Bot", "version", "v0.1.0")
	logger.Debug("Configuration loaded", "log_level", cfg.LogLevel)

	// Create Discord bot instance
	bot, err := discord.New(cfg, logger)
	if err != nil {
		logger.Error("Failed to create Discord bot", "error", err.Error())
		os.Exit(1)
	}

	// Create context for graceful shutdown
	ctx := context.Background()

	// Start the bot (this will block until shutdown)
	if err := bot.Start(ctx); err != nil {
		logger.Error("Bot encountered an error", "error", err.Error())
		os.Exit(1)
	}

	logger.Info("Bot shutdown completed")
}