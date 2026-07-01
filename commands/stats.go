package commands

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"sync"
	"time"

	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/paginator"
	"github.com/tildezero/draftbot/bot"
	"github.com/tildezero/draftbot/pkg/draftout"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
	stats, err := b.Draftout.GetPlayerStats(e.Ctx, player, 1, draftout.MatchFilter(typ))
	if err != nil {
		slog.Error("draftout api call failed for /stats", slog.Any("err", err))
		return e.CreateMessage(discord.NewMessageCreate().
			WithContent("Draftout API call failed, Draftout might be down :(").
			WithEphemeral(true))
	}

	if stats.Player == nil {
		return e.CreateMessage(
			discord.NewMessageCreate().
				WithContent("player not found :(").
				WithEphemeral(true),
		)
	}

	properMode := strings.ReplaceAll(cases.Title(language.English).String(string(stats.Filter)), "-", " ")
	if stats.Record.Matches == 0 {
		return e.CreateMessage(
			discord.NewMessageCreate().
				WithContentf("%s has played no matches of type %s!", strings.ReplaceAll(stats.Player.Username, "_", "\\_"), properMode),
		)
	}

	// EMBED STARTS HERE

	description := buildStaticInfo(stats)

	const matchesPerPage = 5
	const matchesPerApiPage = 20

	totalEmbedPages := int(math.Ceil(float64(stats.Record.Matches) / matchesPerPage))

	type pageCache struct {
		m     sync.Mutex
		cache map[int][]draftout.MatchSummary
	}

	cache := &pageCache{cache: map[int][]draftout.MatchSummary{1: stats.Matches}}

	fetchMatches := func(embedPage int) ([]draftout.MatchSummary, error) {
		matchIdx := embedPage * matchesPerPage
		apiPage := matchIdx/matchesPerApiPage + 1

		cache.m.Lock()
		defer cache.m.Unlock()

		if _, ok := cache.cache[apiPage]; !ok {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			fetched, err := b.Draftout.GetPlayerStats(ctx, player, apiPage, draftout.MatchFilter(typ))
			if err != nil {
				return nil, err
			}
			cache.cache[apiPage] = fetched.Matches
		}

		allMatches := cache.cache[apiPage]
		offsetInApiPage := matchIdx % matchesPerApiPage
		end := min(offsetInApiPage+matchesPerPage, len(allMatches))

		return allMatches[offsetInApiPage:end], nil
	}

	// matchHistoryText := matchHistoryText(stats, player, stats.Filter)

	color := 0x7289da
	if stats.Player.RankColor != nil {
		color, err = draftout.ParseColorString(*stats.Player.RankColor)
		if err != nil {
			color = 0x7289da
		}
	}

	return b.Paginator.Create(e.Respond, paginator.Pages{
		ID: e.ID().String(),
		PageFunc: func(page int, embed discord.Embed) discord.Embed {
			matches, err := fetchMatches(page)
			if err != nil {
				slog.Error("failed to fetch matches for paginator page", slog.Int("page", page), slog.Any("err", err))
				return embed.
					WithTitle("Error").
					WithDescription("Failed to fetch matches from the Draftout API.")
			}

			hist := matchHistoryText(matches, stats.Player, stats.Filter)

			return embed.
				WithTitlef("Draftout Stats for %s (mode: %s)", stats.Player.Username, properMode).
				WithThumbnail("https://mc-heads.net/head/"+stats.Player.Username+"/left/80").
				WithDescription(description).
				WithColor(color).
				WithFooterText(fmt.Sprintf("Stats for mode: %s | Page %d/%d", stats.Filter, page+1, totalEmbedPages)).
				AddField("Matches", hist.String(), false)
		},
		Pages:      totalEmbedPages,
		Creator:    e.User().ID,
		ExpireMode: paginator.ExpireModeAfterLastUsage,
	}, false)

}

func matchHistoryText(matches []draftout.MatchSummary, player *draftout.Player, filter draftout.MatchFilter) strings.Builder {
	var matchHistoryText strings.Builder

	for _, match := range matches {
		if len(match.Participants) < 2 {
			slog.Warn("match has fewer than two participants", slog.Int("matchID", match.ID))
			continue
		}

		var self draftout.Participant
		var other draftout.Participant

		switch {
		case strings.EqualFold(match.Participants[0].UUID, player.UUID):
			self = match.Participants[0]
			other = match.Participants[1]
		case strings.EqualFold(match.Participants[1].UUID, player.UUID):
			self = match.Participants[1]
			other = match.Participants[0]
		case strings.EqualFold(match.Participants[0].Username, player.Username):
			self = match.Participants[0]
			other = match.Participants[1]
		case strings.EqualFold(match.Participants[1].Username, player.Username):
			self = match.Participants[1]
			other = match.Participants[0]
		default:
			slog.Warn("couldn't find player in match participants", slog.Int("matchID", match.ID), slog.String("player", player.Username))
			continue
		}

		var result string
		if match.Outcome == "forfeited" {
			result += "Forfeit "
		}
		if self.Won {
			result += "Win!"
		} else if other.Won {
			result += "Loss :("
		} else {
			result += "Draw!"
		}
		// 															  ID       TIME   DUR       OPP  RES  SCORE
		fmt.Fprintf(&matchHistoryText, "`%d` - <t:%d:R> (%s) - vs. %s - **%s (%d-%d)**", match.ID, match.CompletedAt/1000, draftout.FormatDuration(match.DurationMs), strings.ReplaceAll(other.Username, "_", "\\_"), result, self.Score, other.Score)
		if filter == draftout.FilterCompetitive && self.EloChange != 0 {
			sign := "+"
			if self.EloChange < 0 {
				sign = "-"
			}
			fmt.Fprintf(&matchHistoryText, "\nElo: %d %s %d = %d", self.EloBefore, sign, int(math.Abs(float64(self.EloChange))), self.EloAfter)
		}
		matchHistoryText.WriteString("\n\n")
	}
	return matchHistoryText
}

func buildStaticInfo(stats *draftout.PlayerStats) string {
	description := ""

	if stats.Filter == "competitive" {
		description += fmt.Sprintf("**Elo**: %d (Peak Elo: %d)\n", stats.Player.Elo, *stats.Aggregate.PeakElo)
		description += fmt.Sprintf("**Rank**: %d\n", *stats.Player.Rank)
	}

	description += fmt.Sprintf("**Matches Played**: %d (%d W - %d D - %d L) \t **Win Rate**: %.2f\n", stats.Record.Matches, stats.Record.Wins, stats.Record.Draws, stats.Record.Losses, stats.Record.WinRate*100)
	description += fmt.Sprintf("**Best Streak**: %d\n", stats.Aggregate.BestStreak)
	description += fmt.Sprintf("**Average Goal Diff**: %.2f\n", *stats.Record.AverageGoals)
	description += fmt.Sprintf("**Forfeit Rate**: %.2f\n", (float64(stats.Aggregate.ForfeitCount) / float64(stats.Record.Matches) * 100))
	description += fmt.Sprintf("**Average Finish**: %s \t **PB**: %s", draftout.FormatDuration(int(*stats.Record.AverageFinishTime)), draftout.FormatDuration(*stats.Aggregate.FastestWinMs))

	return description
}
