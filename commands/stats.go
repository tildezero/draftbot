package commands

import (
	"fmt"
	"log/slog"
	"math"

	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/tildezero/draftbot/bot"
	"github.com/tildezero/draftbot/pkg/draftout"
)

var commandStats = discord.SlashCommandCreate{
	Name: "stats", Description: "get the stats for a draftout player",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionString{Name: "player", Description: "Minecraft username of the player", Required: true},
		discord.ApplicationCommandOptionString{Name: "type", Description: "Match type to get stats for", Required: true, Choices: []discord.ApplicationCommandOptionChoiceString{
			{Name: "Competitive", Value: "competitive"},
			{Name: "Quick Play", Value: "quick-play"},
			{Name: "Lobby", Value: "lobby"},
		}},
	},
}

func handleStats(e *handler.CommandEvent, b *bot.Draftbot) error {
	slog.Info("/stats called")
	player := e.SlashCommandInteractionData().String("player")
	typ := e.SlashCommandInteractionData().String("type")
	ps, err := b.Draftout.GetPlayerStats(e.Ctx, player, 1, draftout.MatchFilter(typ))
	if err != nil {
		slog.Error("draftout api call failed for /stats", slog.Any("err", err))
		return e.CreateMessage(discord.NewMessageCreate().
			WithContent("Draftout API call failed, Draftout might be down :(").
			WithEphemeral(true))
	}
	if ps.Player == nil {
		return e.CreateMessage(
			discord.NewMessageCreate().
				WithContent("player not found :(").
				WithEphemeral(true),
		)
	}

	if ps.Record.Matches == 0 {
		return e.CreateMessage(
			discord.NewMessageCreate().
				WithContentf("%s has played no matches of type %s!", strings.ReplaceAll(ps.Player.Username, "_", "\\_"), ps.Filter),
		)
	}

	description := ""

	if ps.Filter == "competitive" {
		description += fmt.Sprintf("**Elo**: %d (Peak Elo: %d)\n", ps.Player.Elo, *ps.Aggregate.PeakElo)
		description += fmt.Sprintf("**Rank**: %d\n", *ps.Player.Rank)
	}

	description += fmt.Sprintf("**Matches Played**: %d (%d W - %d D - %d L) \t **Win Rate**: %.2f\n", ps.Record.Matches, ps.Record.Wins, ps.Record.Draws, ps.Record.Losses, ps.Record.WinRate*100)
	description += fmt.Sprintf("**Best Streak**: %d\n", ps.Aggregate.BestStreak)
	description += fmt.Sprintf("**Average Goal Diff**: %.2f\n", *ps.Record.AverageGoals)
	description += fmt.Sprintf("**Forfeit Rate**: %.2f\n", (float64(ps.Aggregate.ForfeitCount) / float64(ps.Record.Matches) * 100))
	description += fmt.Sprintf("**Average Finish**: %s \t **PB**: %s", draftout.FormatDuration(int(*ps.Record.AverageFinishTime)), draftout.FormatDuration(*ps.Aggregate.FastestWinMs))

	var matchHistoryText strings.Builder

	for _, match := range ps.Matches[0:7] {
		var self draftout.Participant
		var other draftout.Participant
		if match.Participants[0].Username == player {
			self = match.Participants[0]
			other = match.Participants[1]
		} else {
			self = match.Participants[1]
			other = match.Participants[0]
		}
		var result string
		if match.Outcome == "forfeited" {
			result += "Forfeit "
		}
		if self.Won == true {
			result += "Win!"
		} else if self.Won == false && other.Won == true {
			result += "Loss :("
		} else {
			result += "Draw!"
		}
		// 															  ID       TIME   DUR       OPP  RES  SCORE
		fmt.Fprintf(&matchHistoryText, "`%d` - <t:%d:R> (%s) - vs. %s - **%s (%d-%d)**", match.ID, match.CompletedAt/1000, draftout.FormatDuration(match.DurationMs), strings.ReplaceAll(other.Username, "_", "\\_"), result, self.Score, other.Score)
		if ps.Filter == draftout.FilterCompetitive && self.EloChange > 0 {
			sign := "+"
			if self.EloChange < 0 {
				sign = "-"
			}
			fmt.Fprintf(&matchHistoryText, "\nElo: %d %s %d = %d", self.EloBefore, sign, int(math.Abs(float64(self.EloChange))), self.EloAfter)
		}
		matchHistoryText.WriteString("\n\n")
	}

	color := 0x7289da
	if ps.Player.RankColor != nil {
		color, err = draftout.ParseColorString(*ps.Player.RankColor)
		if err != nil {
			color = 0x7289da
		}
	}

	em := discord.NewEmbed().
		WithTitlef("Draftout Stats for %s", ps.Player.Username).
		WithThumbnail("https://mc-heads.net/head/"+ps.Player.Username+"/left/80").
		WithDescription(description).
		WithColor(color).
		WithFooterText(fmt.Sprintf("Stats for mode: %s", ps.Filter)).
		AddField("Matches", matchHistoryText.String(), false)

	return e.CreateMessage(
		discord.NewMessageCreate().WithEmbeds(em),
	)
}
