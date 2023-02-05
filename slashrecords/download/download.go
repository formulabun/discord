package download

import (
	"github.com/bwmarrin/discordgo"
	"go.formulabun.club/discord/slashrecords/common"
	"go.formulabun.club/discord/context"
)

var commandOption = &discordgo.ApplicationCommandOption{
	Type:        discordgo.ApplicationCommandOptionSubCommand,
	Name:        "get",
	Description: "download a record as guest record",
}

var Option = common.CommandOption{commandOption, reply}

func reply(c *context.DiscordContext, options []*discordgo.ApplicationCommandInteractionDataOption) discordgo.InteractionResponse{
  return discordgo.InteractionResponse{}
}
