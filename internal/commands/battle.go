package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/varubogu/gbf_discord_bot_go/internal/gbf"
	"github.com/varubogu/gbf_discord_bot_go/internal/log"
)

// BattleCommand handles battle-related command functionality
type BattleCommand struct {
	logger        *log.Logger
	battleManager *gbf.BattleManager
}

// NewBattleCommand creates a new battle command handler
func NewBattleCommand(logger *log.Logger) *BattleCommand {
	return &BattleCommand{
		logger:        logger,
		battleManager: gbf.NewBattleManager(),
	}
}

// HandleBattlesListPrefixCommand handles the prefix version of battles list command (!battles)
func (b *BattleCommand) HandleBattlesListPrefixCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	logger := b.logger.WithDiscordContext(m.GuildID, m.ChannelID, m.Author.ID).WithCommand("battles")

	// Parse arguments - check if specific type requested
	args := strings.Fields(m.Content)
	var battleType gbf.BattleType
	if len(args) > 1 {
		battleType = gbf.BattleType(strings.ToLower(args[1]))
	}

	embed := b.buildBattlesListEmbed(battleType)

	_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		logger.WithError(err).Error("Failed to send battles list response")
		return
	}

	logger.Info("Battles list prefix command executed successfully", "battle_type", string(battleType))
}

// HandleBattlesListSlashCommand handles the slash version of battles list command (/battles)
func (b *BattleCommand) HandleBattlesListSlashCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	logger := b.logger.WithDiscordContext(i.GuildID, i.ChannelID, i.Member.User.ID).WithCommand("battles")

	// Parse options
	var battleType gbf.BattleType
	if len(i.ApplicationCommandData().Options) > 0 {
		battleType = gbf.BattleType(i.ApplicationCommandData().Options[0].StringValue())
	}

	embed := b.buildBattlesListEmbed(battleType)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
	if err != nil {
		logger.WithError(err).Error("Failed to respond to battles list slash command")
		return
	}

	logger.Info("Battles list slash command executed successfully", "battle_type", string(battleType))
}

// HandleBattleInfoPrefixCommand handles the prefix version of battle info command (!battle <id>)
func (b *BattleCommand) HandleBattleInfoPrefixCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	logger := b.logger.WithDiscordContext(m.GuildID, m.ChannelID, m.Author.ID).WithCommand("battle")

	// Parse arguments
	args := strings.Fields(m.Content)
	if len(args) < 2 {
		_, err := s.ChannelMessageSend(m.ChannelID, "❌ Please specify a battle ID. Usage: `!battle <id>`")
		if err != nil {
			logger.WithError(err).Error("Failed to send usage message")
		}
		return
	}

	battleID := args[1]
	embed := b.buildBattleInfoEmbed(battleID)

	_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		logger.WithError(err).Error("Failed to send battle info response")
		return
	}

	logger.Info("Battle info prefix command executed successfully", "battle_id", battleID)
}

// HandleBattleInfoSlashCommand handles the slash version of battle info command (/battle)
func (b *BattleCommand) HandleBattleInfoSlashCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	logger := b.logger.WithDiscordContext(i.GuildID, i.ChannelID, i.Member.User.ID).WithCommand("battle")

	// Get battle ID from options
	var battleID string
	if len(i.ApplicationCommandData().Options) > 0 {
		battleID = i.ApplicationCommandData().Options[0].StringValue()
	}

	embed := b.buildBattleInfoEmbed(battleID)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
	if err != nil {
		logger.WithError(err).Error("Failed to respond to battle info slash command")
		return
	}

	logger.Info("Battle info slash command executed successfully", "battle_id", battleID)
}

// GetBattlesListSlashCommandDefinition returns the slash command definition for battles list
func (b *BattleCommand) GetBattlesListSlashCommandDefinition() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "battles",
		Description: "Shows list of available battles",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "type",
				Description: "Filter by battle type",
				Required:    false,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{Name: "High Level", Value: "hl"},
					{Name: "Faa HL", Value: "faa_hl"},
					{Name: "Baha HL", Value: "baha_hl"},
					{Name: "UBaha HL", Value: "ubaha_hl"},
					{Name: "Akasha HL", Value: "akasha_hl"},
					{Name: "Luci HL", Value: "luci_hl"},
					{Name: "Guild War", Value: "gw"},
					{Name: "Guild War NM", Value: "gw_nm"},
					{Name: "Event", Value: "event"},
					{Name: "Event HL", Value: "event_hl"},
					{Name: "Training", Value: "train"},
					{Name: "Custom", Value: "custom"},
				},
			},
		},
	}
}

// GetBattleInfoSlashCommandDefinition returns the slash command definition for battle info
func (b *BattleCommand) GetBattleInfoSlashCommandDefinition() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "battle",
		Description: "Shows detailed information about a specific battle",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "id",
				Description: "Battle ID to get information for",
				Required:    true,
			},
		},
	}
}

// buildBattlesListEmbed builds the battles list embed message
func (b *BattleCommand) buildBattlesListEmbed(battleType gbf.BattleType) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Title: "GBF Battle List",
		Color: 0x3498db,
	}

	var battles []*gbf.BattleInfo
	if battleType != "" {
		if !gbf.IsValidBattleType(battleType) {
			embed.Title = "Invalid Battle Type"
			embed.Description = fmt.Sprintf("The battle type `%s` is not valid.", battleType)
			embed.Color = 0xe74c3c
			return embed
		}
		battles = b.battleManager.GetBattlesByType(battleType)
		embed.Title = fmt.Sprintf("GBF Battle List - %s", strings.ToUpper(string(battleType)))
	} else {
		battles = b.battleManager.GetActiveBattles()
		embed.Title = "GBF Battle List - Active Battles"
	}

	if len(battles) == 0 {
		embed.Description = "No battles found for the specified criteria."
		embed.Color = 0xf39c12
		return embed
	}

	// Group battles by type for better organization
	battlesByType := make(map[gbf.BattleType][]*gbf.BattleInfo)
	for _, battle := range battles {
		battlesByType[battle.Type] = append(battlesByType[battle.Type], battle)
	}

	for bType, typeBattles := range battlesByType {
		var battleList []string
		for _, battle := range typeBattles {
			status := "🔴"
			if battle.IsActive {
				status = "🟢"
			}
			battleList = append(battleList, fmt.Sprintf("%s `%s` - %s (Lv.%d)",
				status, battle.ID, battle.Name, battle.Level))
		}

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   fmt.Sprintf("%s Battles", strings.ToUpper(string(bType))),
			Value:  strings.Join(battleList, "\n"),
			Inline: false,
		})
	}

	embed.Footer = &discordgo.MessageEmbedFooter{
		Text: "Use !battle <id> or /battle <id> for detailed information",
	}

	return embed
}

// buildBattleInfoEmbed builds the battle info embed message
func (b *BattleCommand) buildBattleInfoEmbed(battleID string) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Color: 0x3498db,
	}

	if battleID == "" {
		embed.Title = "Battle ID Required"
		embed.Description = "Please specify a battle ID to get information."
		embed.Color = 0xe74c3c
		return embed
	}

	battle, err := b.battleManager.GetBattle(battleID)
	if err != nil {
		embed.Title = "Battle Not Found"
		embed.Description = fmt.Sprintf("No battle found with ID: `%s`", battleID)
		embed.Color = 0xe74c3c
		return embed
	}

	embed.Title = fmt.Sprintf("Battle Info: %s", battle.Name)
	embed.Description = battle.Description

	status := "🔴 Inactive"
	statusColor := 0xe74c3c
	if battle.IsActive {
		status = "🟢 Active"
		statusColor = 0x27ae60
	}
	embed.Color = statusColor

	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "Battle ID",
			Value:  battle.ID,
			Inline: true,
		},
		{
			Name:   "Level",
			Value:  fmt.Sprintf("%d", battle.Level),
			Inline: true,
		},
		{
			Name:   "Type",
			Value:  string(battle.Type),
			Inline: true,
		},
		{
			Name:   "Min Rank",
			Value:  fmt.Sprintf("%d", battle.MinRank),
			Inline: true,
		},
		{
			Name:   "Max Players",
			Value:  fmt.Sprintf("%d", battle.MaxPlayers),
			Inline: true,
		},
		{
			Name:   "Status",
			Value:  status,
			Inline: true,
		},
	}

	embed.Footer = &discordgo.MessageEmbedFooter{
		Text: fmt.Sprintf("Created: %s", battle.CreatedAt.Format("2006-01-02 15:04:05")),
	}

	return embed
}
