package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

var token = flag.String("token", "", "discord bot token (secret)")
var applicationId = flag.String("app", "", "discord application id")
var commandId = flag.String("command", "", "discord slash command id")
var guildId = flag.String("guild", "", "discord guild, for test guilds")

func main() {
	flag.Parse()
	session, err := discordgo.New(fmt.Sprintf("Bot %s", *token))
	if err != nil {
		log.Fatal("Could not start session: ", err)
	}

	log.Fatal(session.ApplicationCommandDelete(*applicationId, *commandId, *guildId))
}
