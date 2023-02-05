package context

import (
	"github.com/bwmarrin/discordgo"
	"go.formulabun.club/replays/store"
)

type DiscordContext struct {
	S        *discordgo.Session
	ReplayDB *store.Client
	Cancel   chan struct{}
}
