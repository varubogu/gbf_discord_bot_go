package discord

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/varubogu/gbf_discord_bot_go/internal/commands"
	"github.com/varubogu/gbf_discord_bot_go/internal/config"
	"github.com/varubogu/gbf_discord_bot_go/internal/log"
)

// Bot represents the Discord bot instance
type Bot struct {
	session       *discordgo.Session
	config        *config.Config
	logger        *log.Logger
	pingCommand   *commands.PingCommand
	helpCommand   *commands.HelpCommand
	adminCommand  *commands.AdminCommand
	battleCommand *commands.BattleCommand
}

// New creates a new Discord bot instance
func New(cfg *config.Config, logger *log.Logger) (*Bot, error) {
	// Create Discord session with bot token
	session, err := discordgo.New("Bot " + cfg.DiscordToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create Discord session: %w", err)
	}

	// Set required intents
	session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent

	bot := &Bot{
		session:       session,
		config:        cfg,
		logger:        logger,
		pingCommand:   commands.NewPingCommand(logger),
		helpCommand:   commands.NewHelpCommand(logger),
		adminCommand:  commands.NewAdminCommand(logger),
		battleCommand: commands.NewBattleCommand(logger),
	}

	// Register event handlers
	bot.setupHandlers()

	return bot, nil
}

// setupHandlers registers event handlers for the bot
func (b *Bot) setupHandlers() {
	b.session.AddHandler(b.onReady)
	b.session.AddHandler(b.onMessageCreate)
	b.session.AddHandler(b.onInteractionCreate)
}

// onReady handles the ready event
func (b *Bot) onReady(s *discordgo.Session, r *discordgo.Ready) {
	b.logger.Info("Bot is ready", "user", r.User.Username, "guilds", len(r.Guilds))
}

// onMessageCreate handles new messages (for prefix commands)
func (b *Bot) onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages from bots
	if m.Author.Bot {
		return
	}

	// Handle prefix commands
	if m.Content == "!ping" {
		b.pingCommand.HandlePrefixCommand(s, m)
	} else if strings.HasPrefix(m.Content, "!help") {
		b.helpCommand.HandlePrefixCommand(s, m)
	} else if m.Content == "!reload" {
		b.adminCommand.HandleReloadPrefixCommand(s, m)
	} else if m.Content == "!status" {
		b.adminCommand.HandleStatusPrefixCommand(s, m)
	} else if strings.HasPrefix(m.Content, "!battles") {
		b.battleCommand.HandleBattlesListPrefixCommand(s, m)
	} else if strings.HasPrefix(m.Content, "!battle ") {
		b.battleCommand.HandleBattleInfoPrefixCommand(s, m)
	}
}

// onInteractionCreate handles slash command interactions
func (b *Bot) onInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.ApplicationCommandData().Name {
	case "ping":
		b.pingCommand.HandleSlashCommand(s, i)
	case "help":
		b.helpCommand.HandleSlashCommand(s, i)
	case "reload":
		b.adminCommand.HandleReloadSlashCommand(s, i)
	case "status":
		b.adminCommand.HandleStatusSlashCommand(s, i)
	case "battles":
		b.battleCommand.HandleBattlesListSlashCommand(s, i)
	case "battle":
		b.battleCommand.HandleBattleInfoSlashCommand(s, i)
	}
}

// Start starts the Discord bot
func (b *Bot) Start(ctx context.Context) error {
	// Open connection to Discord
	err := b.session.Open()
	if err != nil {
		return fmt.Errorf("failed to open Discord connection: %w", err)
	}

	b.logger.Info("Bot started successfully")

	// Wait for context cancellation or interrupt signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		b.logger.Info("Bot stopping due to context cancellation")
	case <-stop:
		b.logger.Info("Bot stopping due to interrupt signal")
	}

	return b.Close()
}

// Close closes the Discord connection
func (b *Bot) Close() error {
	b.logger.Info("Closing Discord connection")
	return b.session.Close()
}

// registerSlashCommands registers slash commands for the bot
func (b *Bot) registerSlashCommands() error {
	commands := []*discordgo.ApplicationCommand{
		b.pingCommand.GetSlashCommandDefinition(),
		b.helpCommand.GetSlashCommandDefinition(),
		b.adminCommand.GetReloadSlashCommandDefinition(),
		b.adminCommand.GetStatusSlashCommandDefinition(),
		b.battleCommand.GetBattlesListSlashCommandDefinition(),
		b.battleCommand.GetBattleInfoSlashCommandDefinition(),
	}

	for _, command := range commands {
		_, err := b.session.ApplicationCommandCreate(b.session.State.User.ID, b.config.TestGuildID, command)
		if err != nil {
			return fmt.Errorf("failed to create slash command %s: %w", command.Name, err)
		}
		b.logger.Info("Registered slash command", "command", command.Name, "guild", b.config.TestGuildID)
	}

	return nil
}
