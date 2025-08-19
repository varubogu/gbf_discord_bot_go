package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/varubogu/gbf_discord_bot_go/internal/log"
)

// HelpCommand handles help command functionality
type HelpCommand struct {
	logger   *log.Logger
	commands []CommandInfo
}

// CommandInfo represents information about a command
type CommandInfo struct {
	Name        string
	Description string
	Usage       string
	Category    string
	IsSlash     bool
	IsPrefix    bool
}

// NewHelpCommand creates a new help command handler
func NewHelpCommand(logger *log.Logger) *HelpCommand {
	return &HelpCommand{
		logger: logger,
		commands: []CommandInfo{
			{
				Name:        "ping",
				Description: "Pings the bot and returns response time",
				Usage:       "!ping or /ping",
				Category:    "General",
				IsSlash:     true,
				IsPrefix:    true,
			},
			{
				Name:        "help",
				Description: "Shows this help message",
				Usage:       "!help [command] or /help [command]",
				Category:    "General",
				IsSlash:     true,
				IsPrefix:    true,
			},
			{
				Name:        "status",
				Description: "Shows bot status information",
				Usage:       "!status or /status",
				Category:    "General",
				IsSlash:     true,
				IsPrefix:    true,
			},
			{
				Name:        "reload",
				Description: "Reloads bot configuration and components",
				Usage:       "!reload or /reload",
				Category:    "Admin",
				IsSlash:     true,
				IsPrefix:    true,
			},
			{
				Name:        "battles",
				Description: "Shows list of available battles",
				Usage:       "!battles [type] or /battles [type]",
				Category:    "GBF",
				IsSlash:     true,
				IsPrefix:    true,
			},
			{
				Name:        "battle",
				Description: "Shows detailed information about a specific battle",
				Usage:       "!battle <id> or /battle <id>",
				Category:    "GBF",
				IsSlash:     true,
				IsPrefix:    true,
			},
		},
	}
}

// HandlePrefixCommand handles the prefix version of help command (!help)
func (h *HelpCommand) HandlePrefixCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	logger := h.logger.WithDiscordContext(m.GuildID, m.ChannelID, m.Author.ID).WithCommand("help")

	// Parse arguments
	args := strings.Fields(m.Content)
	var commandName string
	if len(args) > 1 {
		commandName = args[1]
	}

	embed := h.buildHelpEmbed(commandName, false)

	_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		logger.WithError(err).Error("Failed to send help response")
		return
	}

	logger.Info("Help prefix command executed successfully", "requested_command", commandName)
}

// HandleSlashCommand handles the slash version of help command (/help)
func (h *HelpCommand) HandleSlashCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	logger := h.logger.WithDiscordContext(i.GuildID, i.ChannelID, i.Member.User.ID).WithCommand("help")

	// Parse options
	var commandName string
	if len(i.ApplicationCommandData().Options) > 0 {
		commandName = i.ApplicationCommandData().Options[0].StringValue()
	}

	embed := h.buildHelpEmbed(commandName, true)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
	if err != nil {
		logger.WithError(err).Error("Failed to respond to help slash command")
		return
	}

	logger.Info("Help slash command executed successfully", "requested_command", commandName)
}

// GetSlashCommandDefinition returns the slash command definition for registration
func (h *HelpCommand) GetSlashCommandDefinition() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "help",
		Description: "Shows help information about commands",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "command",
				Description: "Specific command to get help for",
				Required:    false,
			},
		},
	}
}

// buildHelpEmbed builds the help embed message
func (h *HelpCommand) buildHelpEmbed(commandName string, isSlash bool) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Title: "GBF Discord Bot Help",
		Color: 0x3498db,
	}

	if commandName != "" {
		// Show help for specific command
		for _, cmd := range h.commands {
			if strings.EqualFold(cmd.Name, commandName) {
				embed.Title = fmt.Sprintf("Help: %s", cmd.Name)
				embed.Description = cmd.Description
				embed.Fields = []*discordgo.MessageEmbedField{
					{
						Name:   "Usage",
						Value:  cmd.Usage,
						Inline: false,
					},
					{
						Name:   "Category",
						Value:  cmd.Category,
						Inline: true,
					},
					{
						Name:   "Available as",
						Value:  h.getAvailabilityString(cmd),
						Inline: true,
					},
				}
				return embed
			}
		}

		// Command not found
		embed.Title = "Command Not Found"
		embed.Description = fmt.Sprintf("No command found with name: `%s`", commandName)
		embed.Color = 0xe74c3c
		return embed
	}

	// Show all commands grouped by category
	categories := make(map[string][]CommandInfo)
	for _, cmd := range h.commands {
		categories[cmd.Category] = append(categories[cmd.Category], cmd)
	}

	embed.Description = "Available commands grouped by category. Use `!help <command>` or `/help <command>` for detailed information."

	for category, commands := range categories {
		var commandList []string
		for _, cmd := range commands {
			availability := ""
			if isSlash && cmd.IsSlash {
				availability = "/"
			} else if !isSlash && cmd.IsPrefix {
				availability = "!"
			} else if cmd.IsSlash && cmd.IsPrefix {
				availability = "!//"
			}

			if availability != "" {
				commandList = append(commandList, fmt.Sprintf("`%s%s` - %s", availability, cmd.Name, cmd.Description))
			}
		}

		if len(commandList) > 0 {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   category,
				Value:  strings.Join(commandList, "\n"),
				Inline: false,
			})
		}
	}

	embed.Footer = &discordgo.MessageEmbedFooter{
		Text: "GBF Discord Bot - Go Edition",
	}

	return embed
}

// getAvailabilityString returns a string describing command availability
func (h *HelpCommand) getAvailabilityString(cmd CommandInfo) string {
	var parts []string
	if cmd.IsPrefix {
		parts = append(parts, "Prefix (!)")
	}
	if cmd.IsSlash {
		parts = append(parts, "Slash (/)")
	}
	return strings.Join(parts, ", ")
}

// AddCommand adds a new command to the help system
func (h *HelpCommand) AddCommand(info CommandInfo) {
	h.commands = append(h.commands, info)
}

// SetCommands sets the complete list of commands
func (h *HelpCommand) SetCommands(commands []CommandInfo) {
	h.commands = commands
}
