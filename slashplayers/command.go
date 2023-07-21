package slashplayers

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	dContext "go.formulabun.club/discord/context"
	"go.formulabun.club/discord/env"
	"go.formulabun.club/functional/array"
	translator "go.formulabun.club/translator/client"

	"github.com/bwmarrin/discordgo"
)

var trClient *translator.APIClient
var logger = log.New(os.Stdout, "/players | ", log.LstdFlags)
var request translator.ApiPlayerinfoGetRequest

var command = &discordgo.ApplicationCommand{
	Name:        "players",
	Description: "Get player info",
}

func Start(c *dContext.DiscordContext) {
	config := translator.NewConfiguration()
	trClient = translator.NewAPIClient(config)

	logger.Println("Trying to get a connection to the translator service.")
	request = trClient.DefaultApi.PlayerinfoGet(context.Background())
	_, _, err := request.Execute()
	for err != nil {
		<-time.After(time.Second * 5)
		_, _, err = request.Execute()
		logger.Println(err)
	}
	logger.Println("Connection gotten, registering command.")

	command, err = c.S.ApplicationCommandCreate(env.APPLICATIONID, env.TESTGUILD, command)
	if err != nil {
		logger.Fatal(err)
	}

	destroy := c.S.AddHandler(reply)

	for _ = range c.Cancel {
	}

	logger.Println("Destroying command")
	destroy()
}
func reply(s *discordgo.Session, interaction *discordgo.InteractionCreate) {
	if interaction.Interaction.Type != discordgo.InteractionApplicationCommand {
		return
	}
	if interaction.ApplicationCommandData().ID != command.ID {
		return
	}

	logger.Println("Interaction created")

	data, _, err := request.Execute()

	if err != nil {
		logger.Printf("request failed: %s\n", err)
		respondError(s, interaction.Interaction)
		return
	}

	players := array.Filter(data, func(p translator.PlayerInfoEntry) bool {
		return *p.Team != 255
	})
	playerNames := array.Map(players, func(p translator.PlayerInfoEntry) string {
		return p.GetName()
	})
	spectators := array.Filter(data, func(p translator.PlayerInfoEntry) bool {
		return *p.Team == 255
	})
	spectatorNames := array.Map(spectators, func(p translator.PlayerInfoEntry) string {
		return p.GetName()
	})
	response := formatResponse(playerNames, spectatorNames)

	logger.Printf("Sending %s\n", response)
	respondPlayers(s, interaction.Interaction, response)
	logger.Println("Interaction Ending")
}

func respondError(s *discordgo.Session, interaction *discordgo.Interaction) {
	s.InteractionRespond(interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Could not contact the srb2kart server. Seek help!",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

func respondPlayers(s *discordgo.Session, interaction *discordgo.Interaction, response string) {
	s.InteractionRespond(interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("%s", response),
		},
	})
}
