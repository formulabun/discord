package search

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.formulabun.club/replays/store"
	"go.formulabun.club/srb2kart/conversion"
)

func format(searchData store.Replay, replays []store.Replay) []*discordgo.MessageEmbed {
	result := make([]*discordgo.MessageEmbed, 1+len(replays))
	if len(replays) == 0 {
		result[0] = &discordgo.MessageEmbed{
			Type:        discordgo.EmbedTypeRich,
			Title:       "there are no records yet.",
			Description: "Go take this easy first place!",
		}
		return result
	}
	searchGenerator := formatSearchData(searchData)
	recordGenerator := formatRecordData(replays)

	i := 1
	header := <-searchGenerator
	result[0] = &header
	for embed := range recordGenerator {
		em := embed
		result[i] = &em
		i++
	}

	return result
}

func formatSearchData(searchData store.Replay) chan discordgo.MessageEmbed {
	generate := make(chan discordgo.MessageEmbed)
	mapID, _ := conversion.NumberToMapId(uint(searchData.GameMap))
	go func() {
		generate <- discordgo.MessageEmbed{
			Type:        discordgo.EmbedTypeRich,
			Title:       "Records for you",
			Description: fmt.Sprintf("On map %s", mapID),
		}
		close(generate)
	}()

	return generate
}

func formatRecordData(replays []store.Replay) chan discordgo.MessageEmbed {
	generate := make(chan discordgo.MessageEmbed)
	go func() {
		fmt.Println(replays)
		for _, r := range replays {
			recordTime := conversion.FramesToTime(uint(r.Time))
			embed := discordgo.MessageEmbed{
				Type:        discordgo.EmbedTypeRich,
				Title:       fmt.Sprintf("%s - %s (%d, %d)", r.PlayerName[:], r.PlayerSkin, r.Speed, r.Weight),
				Description: fmt.Sprintf("%s", recordTime.Round(time.Millisecond)),
			}
			generate <- embed
		}
		close(generate)
	}()
	return generate
}
