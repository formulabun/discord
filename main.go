package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.formulabun.club/discord/context"
	"go.formulabun.club/discord/env"
	"go.formulabun.club/discord/slashplayers"
	"go.formulabun.club/discord/slashrecords"
	"go.formulabun.club/discord/status"
	"go.formulabun.club/replays/store"

	"github.com/bwmarrin/discordgo"
)

func makeDiscordSession() (*discordgo.Session, error) {
	session, err := discordgo.New(fmt.Sprintf("Bot %s", env.TOKEN))
	if err != nil {
		log.Fatal(fmt.Sprintf("Could not create a client: %s", err))
	}

	session.AddHandler(func(s *discordgo.Session, e *discordgo.Ready) {
		log.Print("Ready")
	})

	err = session.Open()
	if err != nil {
		log.Fatal(fmt.Sprintf("Could not start client: %s", err))
	}
	return session, nil
}

func makeReplayClient() (store.Client, error) {
	return store.NewClient()
}

type Maybe[T any] struct {
	value T
	err   error
}

func expensiveCreate[T any](create func() (T, error)) chan Maybe[T] {
	result := make(chan Maybe[T])
	go func() {
		res, err := create()
		result <- Maybe[T]{res, err}
	}()
	return result
}

func main() {
	env.ValidateEnvironment()

	flag.Parse()

	log.Print("Discord service is starting.")

	laterSession := expensiveCreate(makeDiscordSession)
	laterReplayClient := expensiveCreate(makeReplayClient)

	session := <-laterSession
	replayClient := <-laterReplayClient

	ctx := &context.DiscordContext{session.value, &replayClient.value, make(chan struct{})}

  go status.Start(ctx)
  go slashplayers.Start(ctx)

	if replayClient.err == nil {
		go slashrecords.Start(ctx)
	} else {
		log.Print(replayClient.err)
	}

	waitForShutdown(ctx)
}

func waitForShutdown(c *context.DiscordContext) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)
	<-signals

	log.Print("Closing")
	close(c.Cancel)
	c.S.Close()
}
