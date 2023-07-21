package upload

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/bwmarrin/discordgo"
	client "go.formulabun.club/replays/ingest/client"
)

var logger = log.New(os.Stdout, "/records upload ", log.LstdFlags)

var fileOption = &discordgo.ApplicationCommandOption{
	Type:        discordgo.ApplicationCommandOptionAttachment,
	Name:        "replay",
	Description: "srb2kart replay file",
	Required:    true,
}

var CommandOption = &discordgo.ApplicationCommandOption{
	Type:        discordgo.ApplicationCommandOptionSubCommand,
	Name:        "upload",
	Description: "add your record",
	Options:     []*discordgo.ApplicationCommandOption{fileOption},
}

func errorResponse(err error) discordgo.InteractionResponse {
	logger.Print(err)
	return discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Something went wrong, ||%s||", err),
		},
	}
}

func Reply(options []*discordgo.ApplicationCommandInteractionDataOption, interactionData *discordgo.ApplicationCommandInteractionDataResolved, request client.ApiRootPostRequest) discordgo.InteractionResponse {
	attachmentId := options[0].Value.(string)

	file := interactionData.Attachments[attachmentId]

	tempFile, err := os.CreateTemp("", file.Filename+"*")
	if err != nil {
		return errorResponse(fmt.Errorf("Couldn't create a temporary file for the record: %s", err))
	}

	resp, err := http.Get(file.URL)
	if err != nil {
		return errorResponse(fmt.Errorf("Couldn't download the record: %s", err))
	}

	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		return errorResponse(fmt.Errorf("Couldn't save the record: %s", err))
	}
	resp.Body.Close()

	tempFile.Seek(0, 0)

	request = request.Body(tempFile)
	resp, err = request.Execute()
	if err != nil {
		if resp.StatusCode == http.StatusConflict {
			return discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "It looks like this replay was already added.",
				},
			}
		}
		return errorResponse(fmt.Errorf("could not process replay: %s", resp.Body))
	}

	logger.Print("upload file response: ", resp.Status)

	return discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Thank you for your replay",
		},
	}
}
