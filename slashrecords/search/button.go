package search

import (
	"bytes"
	"encoding/binary"

	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
	"go.formulabun.club/discord/slashrecords/download"
	"go.formulabun.club/replays/store"
)

var node = initNode()

func initNode() *snowflake.Node {
	var name = "Formula bun"
	buff := bytes.NewBufferString(name)
	byteSlice := make([]byte, 64)
	_, err := buff.Read(byteSlice)
	if err != nil {
		logger.Fatal(err)
	}
	value, n := binary.Varint(byteSlice)
	if n == 0 {
		logger.Fatal("Could not create snowflake node")
	}
	node, err := snowflake.NewNode(value)
	if err != nil {
		logger.Fatal(err)
	}
	return node
}

func makeButton(replay store.Replay) discordgo.Button {
	id := node.Generate()
	download.AddPendingDownload(id, replay)
	return discordgo.Button{
		Label:    "Download top replay",
		Style:    discordgo.PrimaryButton,
		Disabled: false,
		CustomID: id.String(),
	}
}
