package draftout

import "context"

func (c *Client) GetRanks(ctx context.Context) ([]Rank, error) {
	var result []Rank
	err := c.get(ctx, "/api/ranks", nil, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
