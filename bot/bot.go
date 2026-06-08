package bot

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/handler"
	"github.com/tildezero/draftbot/pkg/draftout"
)

type Draftbot struct {
	Client   *bot.Client
	Draftout *draftout.Client
}

func New() (*Draftbot, error) {
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		slog.Error("no token set!")
		return nil, errors.New("no token set")
	}
	cord, err := disgo.New(os.Getenv("DISCORD_TOKEN"),
		// set gateway options
		bot.WithGatewayConfigOpts(
			// set enabled intents
			gateway.WithIntents(
				gateway.IntentGuilds,
				gateway.IntentGuildMessages,
				gateway.IntentDirectMessages,
			),
		),
		// bot.WithEventListeners(paginator),
		bot.WithEventListenerFunc(func(e *events.Ready) {
			slog.Info("the bot is ready!")
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := e.Client().SetPresence(ctx, gateway.WithPlayingActivity("🔗 draftbot.suhas.one"), gateway.WithOnlineStatus(discord.OnlineStatusOnline)); err != nil {
				slog.Error("couldn't set presence", slog.Any("err", err))
			}
		}),
	)

	if err != nil {
		return nil, err
	}

	draft := draftout.New()

	return &Draftbot{Client: cord, Draftout: draft}, nil
}

type CmdFunction func(e *handler.CommandEvent, b *Draftbot) error

func (b *Draftbot) Cmd(h CmdFunction) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		return h(e, b)
	}
}
