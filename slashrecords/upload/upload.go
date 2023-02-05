package upload

import (
	"github.com/bwmarrin/discordgo"
	"go.formulabun.club/discord/slashrecords/common"
	"go.formulabun.club/discord/context"
)

var commandOption = &discordgo.ApplicationCommandOption{
	Type:        discordgo.ApplicationCommandOptionSubCommand,
	Name:        "upload",
	Description: "add your record",
}

var Option = common.CommandOption{commandOption, reply}

func reply(c *context.DiscordContext, options []*discordgo.ApplicationCommandInteractionDataOption) discordgo.InteractionResponse{
  return discordgo.InteractionResponse{}
}
