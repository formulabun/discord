package status

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	dContext "go.formulabun.club/discord/context"
	translator "go.formulabun.club/translator/client"

	"github.com/bwmarrin/discordgo"
)

var trClient *translator.APIClient
var logger = log.New(os.Stdout, "/status | ", log.LstdFlags)

func Start(c dContext.DiscordContext) {
	ticker := time.NewTicker(5 * time.Second)

	config := translator.NewConfiguration()
	trClient = translator.NewAPIClient(config)

  setIdleStatus(c.S)
	go setupTimer(ticker, c.S)

	for _ = range c.Cancel {
	}
	ticker.Stop()
}

func setupTimer(tick *time.Ticker, s *discordgo.Session) {
	for _ = range tick.C {
		updateStatus(s)
	}
}

func makeNoStatusData() discordgo.UpdateStatusData {
	i := 0
  status := discordgo.Activity{
    Name: "until you'll help me",
    Type: discordgo.ActivityTypeWatching,
    CreatedAt: time.Now(),
  }

	return discordgo.UpdateStatusData{
		&i,
		[]*discordgo.Activity{&status},
		false,
		string(discordgo.StatusDoNotDisturb),
	}
}

func makeStatusData(info *translator.ServerInfo) discordgo.UpdateStatusData {
	i := 0
	var statusText string
	playerCount := int(info.GetNumberOfPlayer())
	switch playerCount {
	case 0:
		statusText = "an empty map"
	case 1:
		statusText = fmt.Sprintf("%d player race", playerCount)
	default:
		statusText = fmt.Sprintf("%d players race", playerCount)
	}

	status := discordgo.Activity{
		Name:      statusText,
		Type:      discordgo.ActivityTypeWatching,
		CreatedAt: time.Now(),
	}

	return discordgo.UpdateStatusData{
		&i,
		[]*discordgo.Activity{&status},
		false,
		string(discordgo.StatusOnline),
	}
}

func setIdleStatus(s *discordgo.Session) {
  i := 0
  updateStatus := discordgo.UpdateStatusData{
    &i,
    []*discordgo.Activity{},
    true,
    string(discordgo.StatusIdle),
  }
	s.UpdateStatusComplex(updateStatus)
}

func updateStatus(s *discordgo.Session) {
	resp, _, err := trClient.DefaultApi.ServerinfoGet(context.Background()).Execute()
	var updateStatus discordgo.UpdateStatusData
	if err != nil {
		updateStatus = makeNoStatusData()
	} else {
		updateStatus = makeStatusData(resp)
	}
	err = s.UpdateStatusComplex(updateStatus)
	if err != nil {
		logger.Print(fmt.Sprintf("Could not update status: %s", err))
		return
	}
}
