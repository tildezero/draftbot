package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgo/handler/middleware"

	// "github.com/disgoorg/paginator"
	"github.com/disgoorg/snowflake/v2"
	_ "github.com/disgoorg/snowflake/v2"
	"github.com/tildezero/draftbot/bot"
	"github.com/tildezero/draftbot/commands"
)

func main() {

	sync := flag.String("sync", "test", "whether to sync commands or not. debug for the guild in DEBUG_GUILD_ID, global for all guilds, anything else for no sync")

	// set up handler for commands
	h := handler.New()
	h.Use(middleware.Go)
	h.Use(middleware.Logger)

	bot, err := bot.New()
	if err != nil {
		panic("couldn't initialize bot!" + err.Error())
	}

	commands.RegisterCommands(h, bot)

	bot.Client.AddEventListeners(h)

	var guildIDs []snowflake.ID
	var shouldSync bool

	switch *sync {
	case "debug":
		guildID := os.Getenv("TEST_GUILD_ID")
		if guildID == "" {
			panic("TEST_GUILD_ID needs to be set")
		}
		guildIDs = []snowflake.ID{snowflake.MustParse(guildID)}
		shouldSync = true

	case "global":
		guildIDs = []snowflake.ID{}
		shouldSync = true
	}

	if shouldSync {
		if err := handler.SyncCommands(bot.Client, commands.Commands, guildIDs); err != nil {
			slog.Error("couldn't sync commands",
				slog.String("mode", *sync),
				slog.Any("err", err),
			)
		}
	}

	// connect to the gateway
	if err = bot.Client.OpenGateway(context.Background()); err != nil {
		panic(err)
	}

	me, _ := bot.Client.Rest.GetCurrentUser("")
	slog.Info(fmt.Sprintf("Bot is now running, press C-c to exit. Username: %s#%s\n", me.Username, me.Discriminator))
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	<-s
}
