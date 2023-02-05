package slashrecords

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"go.formulabun.club/discord/env"

	dContext "go.formulabun.club/discord/context"

	"go.formulabun.club/discord/slashrecords/download"
	"go.formulabun.club/discord/slashrecords/search"
	"go.formulabun.club/discord/slashrecords/upload"
)

var logger = log.New(os.Stdout, "/records ", log.LstdFlags)
var command = &discordgo.ApplicationCommand{
	Name:        "records",
	Description: "add, search and get time attack records",
	Options:     []*discordgo.ApplicationCommandOption{search.Option.Option, upload.Option.Option, download.Option.Option},
}

var context *dContext.DiscordContext

func Start(c *dContext.DiscordContext) {
  context = c
	_, err := c.S.ApplicationCommandCreate(env.APPICATIONID, env.TESTGUILD, command)
	if err != nil {
		logger.Fatal(err)
	}

	destroy := c.S.AddHandler(reply)

	logger.Println("I'm ready")
	for _ = range c.Cancel {
	}

	logger.Println("Destroying command")
	destroy()
}

func reply(s *discordgo.Session, interaction *discordgo.InteractionCreate) {
	logger.Println("interaction received")
	data := interaction.ApplicationCommandData()
  option := data.Options[0]
  subcommand := option.Name

  var response discordgo.InteractionResponse
  if subcommand == search.Option.Option.Name {
    response = search.Option.Call(context, option.Options)
  } else if subcommand == download.Option.Option.Name {
    response = download.Option.Call(context, option.Options)
  } else if subcommand == upload.Option.Option.Name {
    response = upload.Option.Call(context, option.Options)
  }

  logger.Println("Sending response")
  s.InteractionRespond(interaction.Interaction, &response)
  logger.Println("Interaction ended")
}
