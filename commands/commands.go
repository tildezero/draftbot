package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/tildezero/draftbot/bot"
)

var Commands = []discord.ApplicationCommandCreate{commandStats}

func RegisterCommands(h *handler.Mux, b *bot.Draftbot) {
	h.Command("/stats", b.Cmd(handleStats))
}
