package download

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
	"go.formulabun.club/discord/context"
	"go.formulabun.club/replays/store"
	"go.formulabun.club/srb2kart/conversion"
)

var logger = log.New(os.Stdout, "download button ", log.LstdFlags)

func ReplyFunction(c *context.DiscordContext) func(*discordgo.Session, *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		reply(c, s, i)
	}
}

func reply(c *context.DiscordContext, s *discordgo.Session, interaction *discordgo.InteractionCreate) {
	if interaction.Type != discordgo.InteractionMessageComponent {
		return
	}
	interactionData := interaction.MessageComponentData()
	id, err := snowflake.ParseString(interactionData.CustomID)
	if err != nil {
		logger.Print("bad id")
		return
	}
	data, ok := pendingDownloads[id]
	if !ok {
		response := badButton()
		s.InteractionRespond(interaction.Interaction, &response)
		return
	}

	replays, err := getReplays(c, data)
  if err != nil {
    response := errorResponse(err)
    s.InteractionRespond(interaction.Interaction, &response)
    return
  }
	response, err := makeResponse(replays)
  if err != nil {
    response := errorResponse(err)
    s.InteractionRespond(interaction.Interaction, &response)
    return
  }

	s.InteractionRespond(interaction.Interaction, &response)
}

func getReplays(context *context.DiscordContext, replay store.Replay) ([]store.Replay, error) {
	return context.ReplayDB.FindReplay(replay)
}

func makeResponse(replays []store.Replay) (response discordgo.InteractionResponse, err error) {
	if len(replays) == 0 {
		err = errors.New("No replays to create a response with")
		return
	}
	filename := path.Join("/data", fmt.Sprintf("%v", replays[0].ReplayID))
	file, err := os.Open(filename)
	if err != nil {
		return
	}
  mapID, _ := conversion.NumberToMapId(uint(replays[0].GameMap))
  dFile := discordgo.File{
		Name:  fmt.Sprintf("MAP%s-guest.lmp", mapID),
    ContentType: "application/*",
		Reader: file,
	}

	response = discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Files: []*discordgo.File{&dFile},
		},
	}
	return
}

func badButton() discordgo.InteractionResponse {
	return discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "It seems this button is no longer valid.",
		},
	}
}

func errorResponse(err error) discordgo.InteractionResponse {
  logger.Print(err)
	return discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Something went wrong. ||%s||", err),
		},
	}
}
