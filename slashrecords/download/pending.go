package download

import (
	"github.com/bwmarrin/snowflake"
	"go.formulabun.club/replays/store"
)

var pendingDownloads = make(map[snowflake.ID]store.Replay)

func AddPendingDownload(id snowflake.ID, replay store.Replay) {
	pendingDownloads[id] = replay
}
