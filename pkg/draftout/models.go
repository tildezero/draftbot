package draftout

type MatchFilter string
type MatchType string
type MatchOutcome string
type GameMode string

const (
	FilterCompetitive MatchFilter = "competitive"
	FilterQuickPlay   MatchFilter = "quick-play"
	FilterLobby       MatchFilter = "lobby"

	MatchTypeCompetitive MatchType = "competitive"
	MatchTypeQuickPlay   MatchType = "quick_play"
	MatchTypeLobby       MatchType = "lobby"

	MatchOutcomeFinished   MatchOutcome = "finished"
	MatchOutcomeForfeited  MatchOutcome = "forfeited"
	MatchOutcomeDrawByVote MatchOutcome = "draw_by_vote"

	GameModeDraftout GameMode = "draftout"
	GameModeLockout  GameMode = "lockout"
	GameModeBlackout GameMode = "blackout"
)

// PlayerStatsResponse (/api/stats/{username})

type PlayerStats struct {
	Player     *Player        `json:"player"`
	Record     Record         `json:"record"`
	Aggregate  Aggregate      `json:"aggregate"`
	Matches    []MatchSummary `json:"matches"`
	Page       int            `json:"page"`
	TotalPages int            `json:"totalPages"`
	Filter     MatchFilter    `json:"filter"`
}

type Player struct {
	UUID      string  `json:"uuid"`
	Username  string  `json:"username"`
	Elo       int     `json:"elo"`
	Ranked    bool    `json:"ranked"`
	Rank      *int    `json:"rank"`
	RankName  string  `json:"rankName"`
	RankColor *string `json:"rankColor"`
}

type Record struct {
	Matches           int      `json:"matches"`
	CompletedMatches  int      `json:"completedMatches"`
	Wins              int      `json:"wins"`
	Losses            int      `json:"losses"`
	Draws             int      `json:"draws"`
	WinRate           float64  `json:"winRate"`
	AverageFinishTime *float64 `json:"averageFinishTime"`
	AverageGoals      *float64 `json:"averageGoals"`
}

type Aggregate struct {
	PeakElo      *int `json:"peakElo"`
	BestStreak   int  `json:"bestStreak"`
	FastestWinMs *int `json:"fastestWinMs"`
	ForfeitCount int  `json:"forfeitCount"`
}

type MatchSummary struct {
	ID           int           `json:"id"`
	MatchType    MatchType     `json:"matchType"`
	GameMode     GameMode      `json:"gameMode"`
	Outcome      MatchOutcome  `json:"outcome"`
	CompletedAt  int64         `json:"completedAt"`
	DurationMs   int           `json:"durationMs"`
	Participants []Participant `json:"participants"`
}

type Participant struct {
	UUID      string `json:"uuid"`
	Username  string `json:"username"`
	Won       bool   `json:"won"`
	Score     int    `json:"score"`
	EloBefore int    `json:"eloBefore"`
	EloChange int    `json:"eloChange"`
	EloAfter  int    `json:"eloAfter"`
}

// MatchDetailResponse (/api/stats/{username}/{matchid})
type MatchDetail struct {
	Player *Player `json:"player"`
	Match  *Match  `json:"match"`
}

type Match struct {
	ID           int           `json:"id"`
	MatchType    MatchType     `json:"matchType"`
	GameMode     GameMode      `json:"gameMode"`
	Outcome      MatchOutcome  `json:"outcome"`
	CompletedAt  int64         `json:"completedAt"`
	DurationMs   int           `json:"durationMs"`
	Participants []Participant `json:"participants"`
	Seed         string        `json:"seed"`
	Goals        []Goal        `json:"goals"`
	Draft        Draft         `json:"draft"`
}

type Goal struct {
	Index           int     `json:"index"`
	ID              string  `json:"id"`
	Data            *string `json:"data"`
	Completed       bool    `json:"completed"`
	CompletedByUUID *string `json:"completedByUuid"`
	CompletedAtMs   *int    `json:"completedAtMs"`
}

type Draft struct {
	PickedFirstUUID *string         `json:"pickedFirstUuid"`
	Pool            []DraftPoolItem `json:"pool"`
}

type DraftPoolItem struct {
	ID       string  `json:"id"`
	Data     *string `json:"data"`
	Picked   bool    `json:"picked"`
	TimedOut bool    `json:"timedOut"`
}

// Rank (/api/ranks)
type Rank struct {
	Name  string `json:"name"`
	Min   *int   `json:"min"`
	Max   *int   `json:"max"`
	Color string `json:"color"`
}
