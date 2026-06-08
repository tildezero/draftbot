package draftout

import (
	"context"
	"errors"
	"net/url"
	"strconv"
)

func (c *Client) GetPlayerStats(ctx context.Context, username string, page int, filter MatchFilter) (*PlayerStats, error) {
	if page < 1 {
		page = 1
	}
	var resu PlayerStats
	vals := url.Values{}
	vals.Add("page", strconv.Itoa(page))
	vals.Add("filter", string(filter))

	err := c.get(ctx, "/api/stats/"+username, vals, &resu)
	if err != nil {
		return nil, err
	}
	return &resu, nil
}

func (c *Client) GetPlayerMatch(ctx context.Context, username string, matchId int) (any, error) {
	// TODO
	return nil, errors.New("not implemented yet")
}

func (c *Client) GetPlayerEloSeries(ctx context.Context, username string) (any, error) {
	// TODO
	return nil, errors.New("not implemented yet")
}
