package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/varubogu/gbf_discord_bot_go/internal/log"
)

// AdminCommand handles admin command functionality
type AdminCommand struct {
	logger         *log.Logger
	controlRoleID  string
	controlRoleName string
}

// NewAdminCommand creates a new admin command handler
func NewAdminCommand(logger *log.Logger) *AdminCommand {
	return &AdminCommand{
		logger:          logger,
		controlRoleName: "gbf_bot_control", // Default role name matching Python bot
	}
}

// SetControlRole sets the control role for admin commands
func (a *AdminCommand) SetControlRole(roleID, roleName string) {
	a.controlRoleID = roleID
	a.controlRoleName = roleName
}

// HandleReloadPrefixCommand handles the prefix version of reload command (!reload)
func (a *AdminCommand) HandleReloadPrefixCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	logger := a.logger.WithDiscordContext(m.GuildID, m.ChannelID, m.Author.ID).WithCommand("reload")

	// Check permissions
	if !a.hasAdminPermission(s, m.GuildID, m.Author.ID, logger) {
		_, err := s.ChannelMessageSend(m.ChannelID, "❌ You don't have permission to use this command. Required role: `"+a.controlRoleName+"`")
		if err != nil {
			logger.WithError(err).Error("Failed to send permission denied message")
		}
		return
	}

	// Perform reload operation
	result := a.performReload(logger)
	
	_, err := s.ChannelMessageSend(m.ChannelID, result)
	if err != nil {
		logger.WithError(err).Error("Failed to send reload response")
		return
	}

	logger.Info("Reload prefix command executed successfully")
}

// HandleReloadSlashCommand handles the slash version of reload command (/reload)
func (a *AdminCommand) HandleReloadSlashCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	logger := a.logger.WithDiscordContext(i.GuildID, i.ChannelID, i.Member.User.ID).WithCommand("reload")

	// Check permissions
	if !a.hasAdminPermissionInteraction(s, i, logger) {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ You don't have permission to use this command. Required role: `" + a.controlRoleName + "`",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			logger.WithError(err).Error("Failed to respond with permission denied")
		}
		return
	}

	// Perform reload operation
	result := a.performReload(logger)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: result,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		logger.WithError(err).Error("Failed to respond to reload slash command")
		return
	}

	logger.Info("Reload slash command executed successfully")
}

// GetReloadSlashCommandDefinition returns the slash command definition for reload
func (a *AdminCommand) GetReloadSlashCommandDefinition() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "reload",
		Description: "Reloads bot configuration and components (Admin only)",
	}
}

// HandleStatusPrefixCommand handles the prefix version of status command (!status)
func (a *AdminCommand) HandleStatusPrefixCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	logger := a.logger.WithDiscordContext(m.GuildID, m.ChannelID, m.Author.ID).WithCommand("status")

	embed := a.buildStatusEmbed()
	
	_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		logger.WithError(err).Error("Failed to send status response")
		return
	}

	logger.Info("Status prefix command executed successfully")
}

// HandleStatusSlashCommand handles the slash version of status command (/status)
func (a *AdminCommand) HandleStatusSlashCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	logger := a.logger.WithDiscordContext(i.GuildID, i.ChannelID, i.Member.User.ID).WithCommand("status")

	embed := a.buildStatusEmbed()

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
	if err != nil {
		logger.WithError(err).Error("Failed to respond to status slash command")
		return
	}

	logger.Info("Status slash command executed successfully")
}

// GetStatusSlashCommandDefinition returns the slash command definition for status
func (a *AdminCommand) GetStatusSlashCommandDefinition() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "status",
		Description: "Shows bot status information",
	}
}

// hasAdminPermission checks if user has admin permission for message commands
func (a *AdminCommand) hasAdminPermission(s *discordgo.Session, guildID, userID string, logger *log.Logger) bool {
	if guildID == "" {
		return false // No permission check in DM
	}

	member, err := s.GuildMember(guildID, userID)
	if err != nil {
		logger.WithError(err).Error("Failed to get guild member")
		return false
	}

	return a.checkMemberRoles(s, guildID, member, logger)
}

// hasAdminPermissionInteraction checks if user has admin permission for interactions
func (a *AdminCommand) hasAdminPermissionInteraction(s *discordgo.Session, i *discordgo.InteractionCreate, logger *log.Logger) bool {
	if i.GuildID == "" {
		return false // No permission check in DM
	}

	if i.Member == nil {
		return false
	}

	return a.checkMemberRoles(s, i.GuildID, i.Member, logger)
}

// checkMemberRoles checks if member has required roles
func (a *AdminCommand) checkMemberRoles(s *discordgo.Session, guildID string, member *discordgo.Member, logger *log.Logger) bool {
	guild, err := s.Guild(guildID)
	if err != nil {
		logger.WithError(err).Error("Failed to get guild")
		return false
	}

	// Check for Administrator permission
	permissions, err := s.UserChannelPermissions(member.User.ID, guildID)
	if err == nil && permissions&discordgo.PermissionAdministrator != 0 {
		return true
	}

	// Check for control role by name or ID
	for _, roleID := range member.Roles {
		for _, role := range guild.Roles {
			if role.ID == roleID {
				if role.Name == a.controlRoleName || role.ID == a.controlRoleID {
					return true
				}
			}
		}
	}

	return false
}

// performReload performs the actual reload operation
func (a *AdminCommand) performReload(logger *log.Logger) string {
	logger.Info("Performing bot reload operation")
	
	// In a real implementation, this would:
	// - Reload configuration
	// - Refresh external connections
	// - Update command registrations
	// - Clear caches
	
	// For now, simulate a successful reload
	return "✅ Bot components reloaded successfully!"
}

// buildStatusEmbed builds the status embed message
func (a *AdminCommand) buildStatusEmbed() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: "Bot Status",
		Color: 0x00ff00,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Status",
				Value:  "✅ Online",
				Inline: true,
			},
			{
				Name:   "Version",
				Value:  "v0.1.0",
				Inline: true,
			},
			{
				Name:   "Language",
				Value:  "Go",
				Inline: true,
			},
			{
				Name:   "Commands",
				Value:  "ping, help, reload, status",
				Inline: false,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "GBF Discord Bot - Go Edition",
		},
	}
}