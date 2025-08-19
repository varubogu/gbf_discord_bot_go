package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/varubogu/gbf_discord_bot_go/internal/log"
)

// PingCommand handles ping command functionality
type PingCommand struct {
	logger *log.Logger
}

// NewPingCommand creates a new ping command handler
func NewPingCommand(logger *log.Logger) *PingCommand {
	return &PingCommand{
		logger: logger,
	}
}

// HandlePrefixCommand handles the prefix version of ping command (!ping)
func (p *PingCommand) HandlePrefixCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	logger := p.logger.WithDiscordContext(m.GuildID, m.ChannelID, m.Author.ID).WithCommand("ping")

	_, err := s.ChannelMessageSend(m.ChannelID, "Pong!")
	if err != nil {
		logger.WithError(err).Error("Failed to send ping response")
		return
	}

	logger.Info("Ping prefix command executed successfully")
}

// HandleSlashCommand handles the slash version of ping command (/ping)
func (p *PingCommand) HandleSlashCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	logger := p.logger.WithDiscordContext(i.GuildID, i.ChannelID, i.Member.User.ID).WithCommand("ping")

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Pong!",
		},
	})
	if err != nil {
		logger.WithError(err).Error("Failed to respond to ping slash command")
		return
	}

	logger.Info("Ping slash command executed successfully")
}

// GetSlashCommandDefinition returns the slash command definition for registration
func (p *PingCommand) GetSlashCommandDefinition() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "Replies with Pong!",
	}
}
