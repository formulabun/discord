package common

import (
	"github.com/bwmarrin/discordgo"
  "go.formulabun.club/discord/context"
)

type CommandOption struct {
	Option *discordgo.ApplicationCommandOption
	Call   func(*context.DiscordContext,
		[]*discordgo.ApplicationCommandInteractionDataOption) discordgo.InteractionResponse
}
