package search

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"go.formulabun.club/discord/context"
	"go.formulabun.club/replays/store"
	"go.formulabun.club/srb2kart/conversion"
)

var logger = log.New(os.Stdout, "/records search ", log.LstdFlags)

var oneFloat = float64(1)

var mapOption = &discordgo.ApplicationCommandOption{
	Type:        discordgo.ApplicationCommandOptionString,
	Name:        "map",
	Description: "the name of the map",
	Required:    true,
	MinValue:    &oneFloat,
}

var CommandOption = &discordgo.ApplicationCommandOption{
	Type:        discordgo.ApplicationCommandOptionSubCommand,
	Name:        "search",
	Description: "Search records",
	Options:     []*discordgo.ApplicationCommandOption{mapOption},
}

func Reply(c *context.DiscordContext, options []*discordgo.ApplicationCommandInteractionDataOption) discordgo.InteractionResponse {
	var searchData store.Replay

	for _, option := range options {
		switch option.Name {
		case mapOption.Name:
			// TODO input error handling
			mapID, _ := conversion.MapIdToNumber(option.Value.(string))
			searchData.GameMap = uint16(mapID)
		}
	}

	replays, err := c.ReplayDB.FindReplay(searchData)
	if err != nil {
		return discordgo.InteractionResponse{
			discordgo.InteractionResponseChannelMessageWithSource,
			&discordgo.InteractionResponseData{Content: "Something terrible happened D:", Flags: discordgo.MessageFlagsEphemeral},
		}
	}

	var components []discordgo.MessageComponent
	embeds := format(searchData, replays)

	if len(replays) > 0 {
		button := makeButton(searchData)

		components = []discordgo.MessageComponent{discordgo.ActionsRow{
			[]discordgo.MessageComponent{
				button,
			},
		}}
	}

	return discordgo.InteractionResponse{
		discordgo.InteractionResponseChannelMessageWithSource,
		&discordgo.InteractionResponseData{
			Embeds:     embeds,
			Components: components,
		},
	}
}
