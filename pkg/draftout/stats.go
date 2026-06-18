package draftout

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

func (c *Client) GetPlayerStats(ctx context.Context, username string, page int, filter MatchFilter) (*PlayerStats, error) {
	if page < 1 {
		page = 1
	}
	var result PlayerStats
	vals := url.Values{}
	vals.Add("page", strconv.Itoa(page))
	vals.Add("filter", string(filter))

	err := c.get(ctx, "/api/stats/"+username, vals, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetPlayerMatch(ctx context.Context, username string, matchId int) (*MatchDetail, error) {
	var result MatchDetail

	err := c.get(ctx, fmt.Sprintf("/api/stats/%s/%d", username, matchId), nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetPlayerEloSeries(ctx context.Context, username string) (any, error) {
	var result EloSeries

	err := c.get(ctx, fmt.Sprintf("/api/stats/%s/elo-series", username), nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
