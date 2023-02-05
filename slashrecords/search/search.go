package search

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.formulabun.club/discord/context"
	"go.formulabun.club/discord/slashrecords/common"
	"go.formulabun.club/replays/store"
	"go.formulabun.club/srb2kart/conversion"
)

var logger = log.New(os.Stdout, "/records search ", log.LstdFlags)

var oneFloat = float64(1)

var mapOption = &discordgo.ApplicationCommandOption{
	Type:        discordgo.ApplicationCommandOptionString,
	Name:        "map",
	Description: "the name of the map",
	MinValue:    &oneFloat,
}

var commandOption = &discordgo.ApplicationCommandOption{
	Type:        discordgo.ApplicationCommandOptionSubCommand,
	Name:        "search",
	Description: "Search records",
	Options:     []*discordgo.ApplicationCommandOption{mapOption},
}

var Option = common.CommandOption{commandOption, reply}

func reply(c *context.DiscordContext, options []*discordgo.ApplicationCommandInteractionDataOption) discordgo.InteractionResponse {
	var searchData store.Replay

	for _, option := range options {
		logger.Println(option.Name, option.Value)
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

	embed := format(replays)

	return discordgo.InteractionResponse{
		discordgo.InteractionResponseChannelMessageWithSource,
		&discordgo.InteractionResponseData{Embeds: embed},
	}
}

func format(replays []store.Replay) []*discordgo.MessageEmbed {
	result := make([]*discordgo.MessageEmbed, len(replays))

	for i, r := range replays {
		recordTime := conversion.FramesToTime(uint(r.Time))
		result[i] = &discordgo.MessageEmbed{
			Type:        discordgo.EmbedTypeRich,
			Title:       fmt.Sprintf("%s (%d, %d)", r.PlayerName[:], r.Speed, r.Weight),
			Description: fmt.Sprintf("%s", recordTime.Round(time.Millisecond)),
		}
	}

	return result
}
