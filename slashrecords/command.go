package slashrecords

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"

	dContext "go.formulabun.club/discord/context"
	ingest "go.formulabun.club/replays/ingest/client"

	"go.formulabun.club/discord/env"
	"go.formulabun.club/discord/slashrecords/download"
	"go.formulabun.club/discord/slashrecords/search"
	"go.formulabun.club/discord/slashrecords/upload"
)

var ingestClient *ingest.APIClient
var request ingest.ApiRootPostRequest
var logger = log.New(os.Stdout, "/records ", log.LstdFlags)

var command = &discordgo.ApplicationCommand{
	Name:        "records",
	Description: "add, search and get time attack records",
	Options:     []*discordgo.ApplicationCommandOption{search.CommandOption, upload.CommandOption},
}

var localContext *dContext.DiscordContext

func Start(c *dContext.DiscordContext) {
	config := ingest.NewConfiguration()
	ingestClient = ingest.NewAPIClient(config)
	localContext = c
	request = ingestClient.DefaultApi.RootPost(context.Background())

	logger.Println("Trying to get a connection to the replays ingest service.")
	get := ingestClient.DefaultApi.ListGet(context.Background())
	_, _, err := get.Execute()
	for err != nil {
		<-time.After(time.Second * 5)
		_, _, err = get.Execute()
		logger.Println(err)
	}
	logger.Println("Connection gotten, registering command.")

	_, err = c.S.ApplicationCommandCreate(env.APPLICATIONID, env.TESTGUILD, command)
	if err != nil {
		logger.Fatal("Could not create command: ", err)
	}

	destroySlash := c.S.AddHandler(reply)
	destroyButton := c.S.AddHandler(download.ReplyFunction(c))

	logger.Println("I'm ready")
	for _ = range c.Cancel {
	}

	logger.Println("Destroying handlers")
	destroySlash()
	destroyButton()
}

func reply(s *discordgo.Session, interaction *discordgo.InteractionCreate) {
	if interaction.Type != discordgo.InteractionApplicationCommand {
		return
	}

	logger.Println("interaction received")
	data := interaction.ApplicationCommandData()
	option := data.Options[0]
	subcommand := option.Name

	var response discordgo.InteractionResponse
	if subcommand == search.CommandOption.Name {
		response = search.Reply(localContext, option.Options)
	} else if subcommand == upload.CommandOption.Name {
		response = upload.Reply(option.Options, data.Resolved, request)
	}

	logger.Println("Sending response")
	err := s.InteractionRespond(interaction.Interaction, &response)
	if err != nil {
		logger.Print(err)
	}
	logger.Println("Interaction ended")
}
